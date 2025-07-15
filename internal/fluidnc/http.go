package fluidnc

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SendHTTPCommand sends a command via HTTP POST to /command endpoint
func (c *Client) SendHTTPCommand(command string, silent bool) (string, error) {
	endpoint := "command"
	if silent {
		endpoint = "command_silent"
	}

	url := fmt.Sprintf("http://%s:%d/%s", c.config.Host, c.config.Port, endpoint)
	resp, err := http.Post(url, "text/plain", strings.NewReader(command))
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("command failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(bodyBytes), nil
}

// HTTPFeedHold sends feed hold via HTTP
func (c *Client) HTTPFeedHold() error {
	url := fmt.Sprintf("http://%s:%d/feedhold_reload", c.config.Host, c.config.Port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send feed hold: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("feed hold failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// HTTPCycleStart sends cycle start via HTTP
func (c *Client) HTTPCycleStart() error {
	url := fmt.Sprintf("http://%s:%d/cyclestart_reload", c.config.Host, c.config.Port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send cycle start: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cycle start failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// HTTPRestart sends restart command via HTTP
func (c *Client) HTTPRestart() error {
	url := fmt.Sprintf("http://%s:%d/restart_reload", c.config.Host, c.config.Port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send restart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("restart failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// CheckDidRestart checks if a restart has occurred
func (c *Client) CheckDidRestart() (bool, error) {
	url := fmt.Sprintf("http://%s:%d/did_restart", c.config.Host, c.config.Port)
	resp, err := c.client.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check restart status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("check restart failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response := strings.TrimSpace(string(bodyBytes))
	return strings.ToLower(response) == "true" || response == "1", nil
}
