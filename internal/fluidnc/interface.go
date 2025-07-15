package fluidnc

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ClientInterface defines the interface for FluidNC client operations
type ClientInterface interface {
	// Connection management
	Connect() error
	Disconnect() error

	// WebSocket operations
	SendCommand(command string) (string, error)
	SendRealTimeCommand(command byte) error

	// HTTP operations
	SendHTTPCommand(command string, silent bool) (string, error)

	// Status and monitoring
	GetStatus() (*FluidNCStatus, error)
	ParseStatus(response string) *FluidNCStatus
	DisplayStatus(status *FluidNCStatus)
	MonitorStatus(ctx context.Context, callback func(*FluidNCStatus)) error

	// Control operations
	FeedHold() error
	CycleStart() error
	SoftReset() error
	Home() error
	Unlock() error

	// HTTP control operations
	HTTPFeedHold() error
	HTTPCycleStart() error
	HTTPRestart() error
	CheckDidRestart() (bool, error)

	// File operations
	ListFiles() (*FileListResponse, error)
	UploadFile(filePath, destination, endpoint string) error
	UploadToLocalFS(filePath string, destinationName string) error
	UploadToSD(filePath string, destinationName string) error
	UpdateFirmware(firmwarePath string) error

	// Information
	GetAlarms() ([]AlarmInfo, error)
	GetSettings() (string, error)
	GetCommands() (string, error)
	GetVersion() (string, error)

	// G-code execution
	RunGCodeFile(filePath string, monitor bool) error
	InteractiveMode() error
}

// Validate that Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)

// ClientOptions provides options for creating a new client
type ClientOptions struct {
	Config         *Config
	HTTPClient     *http.Client
	ConnectTimeout time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
}

// NewClientWithOptions creates a new FluidNC client with custom options
func NewClientWithOptions(opts *ClientOptions) *Client {
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{Timeout: opts.Config.Timeout}
	}

	client := NewClient(opts.Config)
	client.client = opts.HTTPClient

	return client
}

// IsConnected returns true if WebSocket connection is active
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// Ping sends a ping to test connection
func (c *Client) Ping() error {
	url := fmt.Sprintf("http://%s:%d/", c.config.Host, c.config.Port)
	resp, err := c.client.Get(url)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping failed with status: %d", resp.StatusCode)
	}

	return nil
}
