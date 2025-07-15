package fluidnc

// FeedHold sends feed hold command
func (c *Client) FeedHold() error {
	return c.SendRealTimeCommand('!')
}

// CycleStart sends cycle start command
func (c *Client) CycleStart() error {
	return c.SendRealTimeCommand('~')
}

// SoftReset sends soft reset command
func (c *Client) SoftReset() error {
	return c.SendRealTimeCommand(0x18) // Ctrl-X
}

// Home sends homing command
func (c *Client) Home() error {
	_, err := c.SendCommand("$H")
	return err
}

// Unlock sends unlock command
func (c *Client) Unlock() error {
	_, err := c.SendCommand("$X")
	return err
}

// GetSettings gets FluidNC settings
func (c *Client) GetSettings() (string, error) {
	return c.SendCommand("$$")
}

// GetCommands gets available commands
func (c *Client) GetCommands() (string, error) {
	return c.SendCommand("$")
}

// GetVersion gets FluidNC version
func (c *Client) GetVersion() (string, error) {
	return c.SendCommand("$I")
}
