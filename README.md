## G-Code Streaming

### Basic Streaming

Stream a G-code file with real-time progress monitoring:

```bash
# Stream with default settings
fluidnc-cli stream job.gcode

# Stream with verbose output
fluidnc-cli --verbose stream job.gcode

# Continue on errors
fluidnc-cli stream --continue-on-error job.gcode

# JSON status output
fluidnc-cli --output json stream job.gcode
