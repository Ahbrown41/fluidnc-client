package fluidnc

import "time"

// Config represents the application configuration
type Config struct {
	Host           string        `yaml:"host" mapstructure:"host"`
	Port           int           `yaml:"port" mapstructure:"port"`
	WebSocketPort  int           `yaml:"websocket_port" mapstructure:"websocket_port"`
	Timeout        time.Duration `yaml:"timeout" mapstructure:"timeout"`
	RetryAttempts  int           `yaml:"retry_attempts" mapstructure:"retry_attempts"`
	RetryDelay     time.Duration `yaml:"retry_delay" mapstructure:"retry_delay"`
	OutputFormat   string        `yaml:"output_format" mapstructure:"output_format"`
	Verbose        bool          `yaml:"verbose" mapstructure:"verbose"`
	StatusInterval time.Duration `yaml:"status_interval" mapstructure:"status_interval"`
	CommandDelay   time.Duration `yaml:"command_delay" mapstructure:"command_delay"`
}

// FluidNCStatus represents parsed status from FluidNC
type FluidNCStatus struct {
	State        string    `json:"state"`
	MachinePos   Position  `json:"machine_position"`
	WorkPos      Position  `json:"work_position"`
	FeedRate     int       `json:"feed_rate"`
	SpindleSpeed int       `json:"spindle_speed"`
	Overrides    Overrides `json:"overrides"`
	Pins         string    `json:"pins"`
	Buffer       Buffer    `json:"buffer"`
	LineNumber   int       `json:"line_number"`
	Timestamp    time.Time `json:"timestamp"`
	Raw          string    `json:"raw_response"`
}

// Position represents X, Y, Z coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Overrides represents feed/rapid/spindle overrides
type Overrides struct {
	Feed    int `json:"feed"`
	Rapid   int `json:"rapid"`
	Spindle int `json:"spindle"`
}

// Buffer represents planner and serial buffer info
type Buffer struct {
	Planner int `json:"planner"`
	Serial  int `json:"serial"`
}

// AlarmInfo represents alarm information
type AlarmInfo struct {
	Code        int       `json:"code"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

// FileInfo represents file information from FluidNC
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"` // "file" or "directory"
}

// FileListResponse represents the response from listing files
type FileListResponse struct {
	Files []FileInfo `json:"files"`
	Path  string     `json:"path"`
}
