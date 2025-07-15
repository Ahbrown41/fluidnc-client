package fluidnc

import (
	"encoding/json"
	"fmt"
)

// DisplayStatus shows current status
func (c *Client) DisplayStatus(status *FluidNCStatus) {
	if c.config.OutputFormat == "json" {
		jsonOutput, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Printf("\rState: %s | MPos: X%.3f Y%.3f Z%.3f | WPos: X%.3f Y%.3f Z%.3f | F:%d S:%d | Line:%d",
			status.State,
			status.MachinePos.X, status.MachinePos.Y, status.MachinePos.Z,
			status.WorkPos.X, status.WorkPos.Y, status.WorkPos.Z,
			status.FeedRate, status.SpindleSpeed,
			status.LineNumber,
		)
	}
}
