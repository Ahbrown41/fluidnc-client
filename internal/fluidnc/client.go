package fluidnc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client handles communication with FluidNC
type Client struct {
	config      *Config
	client      *http.Client
	conn        *websocket.Conn
	mu          sync.RWMutex
	monitoring  bool
	statusRegex *regexp.Regexp
	alarmRegex  *regexp.Regexp
	errorRegex  *regexp.Regexp
}

// NewClient creates a new FluidNC client
func NewClient(config *Config) *Client {
	// Regex patterns for parsing FluidNC responses
	statusRegex := regexp.MustCompile(`<([^|]+)(?:\|MPos:([^|]+))?(?:\|WPos:([^|]+))?(?:\|FS:([^|]+))?(?:\|Ov:([^|]+))?(?:\|Pn:([^|]+))?(?:\|Bf:([^|]+))?(?:\|Ln:([^>]+))?>`)
	alarmRegex := regexp.MustCompile(`ALARM:(\d+)`)
	errorRegex := regexp.MustCompile(`error:(\d+)`)

	return &Client{
		config:      config,
		client:      &http.Client{Timeout: config.Timeout},
		statusRegex: statusRegex,
		alarmRegex:  alarmRegex,
		errorRegex:  errorRegex,
	}
}

// Connect establishes WebSocket connection
func (c *Client) Connect() error {
	wsURL := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.config.Host, c.config.WebSocketPort),
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	return nil
}

// Disconnect closes WebSocket connection
func (c *Client) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SendCommand sends a command and waits for response
func (c *Client) SendCommand(command string) (string, error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return "", fmt.Errorf("not connected")
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(command+"\n")); err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return strings.TrimSpace(string(message)), nil
}

// SendRealTimeCommand sends real-time command (no newline, immediate)
func (c *Client) SendRealTimeCommand(command byte) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	return conn.WriteMessage(websocket.TextMessage, []byte{command})
}

// ParseStatus parses FluidNC status response
func (c *Client) ParseStatus(response string) *FluidNCStatus {
	status := &FluidNCStatus{
		Timestamp: time.Now(),
		Raw:       response,
	}

	matches := c.statusRegex.FindStringSubmatch(response)
	if len(matches) < 2 {
		return status
	}

	status.State = matches[1]

	// Parse machine position (MPos)
	if len(matches) > 2 && matches[2] != "" {
		coords := strings.Split(matches[2], ",")
		if len(coords) >= 3 {
			status.MachinePos.X, _ = strconv.ParseFloat(coords[0], 64)
			status.MachinePos.Y, _ = strconv.ParseFloat(coords[1], 64)
			status.MachinePos.Z, _ = strconv.ParseFloat(coords[2], 64)
		}
	}

	// Parse work position (WPos)
	if len(matches) > 3 && matches[3] != "" {
		coords := strings.Split(matches[3], ",")
		if len(coords) >= 3 {
			status.WorkPos.X, _ = strconv.ParseFloat(coords[0], 64)
			status.WorkPos.Y, _ = strconv.ParseFloat(coords[1], 64)
			status.WorkPos.Z, _ = strconv.ParseFloat(coords[2], 64)
		}
	}

	// Parse feed rate and spindle speed (FS)
	if len(matches) > 4 && matches[4] != "" {
		fs := strings.Split(matches[4], ",")
		if len(fs) >= 2 {
			status.FeedRate, _ = strconv.Atoi(fs[0])
			status.SpindleSpeed, _ = strconv.Atoi(fs[1])
		}
	}

	// Parse overrides (Ov)
	if len(matches) > 5 && matches[5] != "" {
		ov := strings.Split(matches[5], ",")
		if len(ov) >= 3 {
			status.Overrides.Feed, _ = strconv.Atoi(ov[0])
			status.Overrides.Rapid, _ = strconv.Atoi(ov[1])
			status.Overrides.Spindle, _ = strconv.Atoi(ov[2])
		}
	}

	// Parse pins (Pn)
	if len(matches) > 6 && matches[6] != "" {
		status.Pins = matches[6]
	}

	// Parse buffer (Bf)
	if len(matches) > 7 && matches[7] != "" {
		bf := strings.Split(matches[7], ",")
		if len(bf) >= 2 {
			status.Buffer.Planner, _ = strconv.Atoi(bf[0])
			status.Buffer.Serial, _ = strconv.Atoi(bf[1])
		}
	}

	// Parse line number (Ln)
	if len(matches) > 8 && matches[8] != "" {
		status.LineNumber, _ = strconv.Atoi(matches[8])
	}

	return status
}

// MonitorStatus continuously monitors FluidNC status
func (c *Client) MonitorStatus(ctx context.Context, callback func(*FluidNCStatus)) error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Disconnect()

	c.mu.Lock()
	c.monitoring = true
	c.mu.Unlock()

	ticker := time.NewTicker(c.config.StatusInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.mu.Lock()
			c.monitoring = false
			c.mu.Unlock()
			return nil
		case <-ticker.C:
			if err := c.SendRealTimeCommand('?'); err != nil {
				if c.config.Verbose {
					fmt.Printf("Error requesting status: %v\n", err)
				}
				continue
			}

			// Read response
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if c.config.Verbose {
					fmt.Printf("Error reading status: %v\n", err)
				}
				continue
			}

			response := strings.TrimSpace(string(message))
			status := c.ParseStatus(response)

			if callback != nil {
				callback(status)
			}
		}
	}
}

// GetStatus requests current status
func (c *Client) GetStatus() (*FluidNCStatus, error) {
	response, err := c.SendCommand("?")
	if err != nil {
		return nil, err
	}

	return c.ParseStatus(response), nil
}

// GetAlarms requests alarm information
func (c *Client) GetAlarms() ([]AlarmInfo, error) {
	response, err := c.SendCommand("$alarms")
	if err != nil {
		return nil, err
	}

	var alarms []AlarmInfo
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ALARM:") {
			matches := c.alarmRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				code, _ := strconv.Atoi(matches[1])
				alarms = append(alarms, AlarmInfo{
					Code:        code,
					Description: c.getAlarmDescription(code),
					Timestamp:   time.Now(),
				})
			}
		}
	}

	return alarms, nil
}

// getAlarmDescription returns description for alarm code
func (c *Client) getAlarmDescription(code int) string {
	descriptions := map[int]string{
		1:  "Hard limit triggered",
		2:  "G-code motion target exceeds machine travel",
		3:  "Reset while in motion",
		4:  "Probe fail",
		5:  "Probe fail",
		6:  "Homing fail",
		7:  "Homing fail",
		8:  "Homing fail",
		9:  "Homing fail",
		10: "Homing fail",
	}

	if desc, exists := descriptions[code]; exists {
		return desc
	}
	return fmt.Sprintf("Unknown alarm code: %d", code)
}
