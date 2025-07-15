package fluidnc

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// RunGCodeFile sends G-code commands from file line by line
func (c *Client) RunGCodeFile(filePath string, monitor bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Disconnect()

	// Start monitoring if requested
	var ctx context.Context
	var cancel context.CancelFunc
	if monitor {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()

		go c.MonitorStatus(ctx, func(status *FluidNCStatus) {
			c.DisplayStatus(status)
		})
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "(") {
			continue
		}

		if c.config.Verbose {
			fmt.Printf("Line %d: %s\n", lineNum, line)
		}

		response, err := c.SendCommand(line)
		if err != nil {
			return fmt.Errorf("error on line %d: %w", lineNum, err)
		}

		// Check for errors
		if strings.Contains(strings.ToLower(response), "error") {
			return fmt.Errorf("FluidNC error on line %d: %s", lineNum, response)
		}

		if c.config.Verbose {
			fmt.Printf("Response: %s\n", response)
		}

		// Small delay between commands
		if c.config.CommandDelay > 0 {
			time.Sleep(c.config.CommandDelay)
		}
	}

	return scanner.Err()
}

// InteractiveMode starts interactive WebSocket session
func (c *Client) InteractiveMode() error {
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Disconnect()

	fmt.Println("Connected to FluidNC. Commands: exit, status, alarms, hold, start, reset, home, unlock")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())

		if command == "exit" {
			break
		}

		if command == "" {
			fmt.Print("> ")
			continue
		}

		// Handle special commands
		switch command {
		case "status":
			if err := c.SendRealTimeCommand('?'); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				_, message, _ := c.conn.ReadMessage()
				status := c.ParseStatus(strings.TrimSpace(string(message)))
				c.DisplayStatus(status)
				fmt.Println()
			}
		case "hold":
			if err := c.FeedHold(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Feed hold sent")
			}
		case "start":
			if err := c.CycleStart(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Cycle start sent")
			}
		case "reset":
			if err := c.SoftReset(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Soft reset sent")
			}
		case "home":
			if err := c.Home(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Homing started")
			}
		case "unlock":
			if err := c.Unlock(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Machine unlocked")
			}
		case "alarms":
			alarms, err := c.GetAlarms()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				if len(alarms) == 0 {
					fmt.Println("No active alarms")
					continue
				}
				for _, alarm := range alarms {
					fmt.Printf("Alarm %d: %s\n", alarm.Code, alarm.Description)
				}
			}
		default:
			response, err := c.SendCommand(command)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("< %s\n", response)
			}
		}

		fmt.Print("> ")
	}

	return scanner.Err()
}
