package types

import "time"

// Workflow represents a complete workflow definition
type Workflow struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Variables   map[string]string `yaml:"variables"`
	Tasks       []Task            `yaml:"tasks"`
}

// Task represents a single task in the workflow
type Task struct {
	Name            string            `yaml:"name"`
	Type            string            `yaml:"type"` // shell, http, file
	Command         string            `yaml:"command,omitempty"`
	Script          string            `yaml:"script,omitempty"`
	URL             string            `yaml:"url,omitempty"`
	Method          string            `yaml:"method,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Body            string            `yaml:"body,omitempty"`
	FilePath        string            `yaml:"file_path,omitempty"`
	FileContent     string            `yaml:"file_content,omitempty"`
	FileAction      string            `yaml:"file_action,omitempty"` // read, write, delete, copy
	SourcePath      string            `yaml:"source_path,omitempty"`
	DestPath        string            `yaml:"dest_path,omitempty"`
	DependsOn       []string          `yaml:"depends_on,omitempty"`
	Condition       string            `yaml:"condition,omitempty"`
	Retry           RetryConfig       `yaml:"retry,omitempty"`
	Timeout         time.Duration     `yaml:"timeout,omitempty"`
	ContinueOnError bool              `yaml:"continue_on_error,omitempty"`
	Parallel        bool              `yaml:"parallel,omitempty"`
}

// RetryConfig defines retry behavior for a task
type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts"`
	Delay       time.Duration `yaml:"delay"`
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskName  string
	Success   bool
	Error     error
	Output    string
	StartTime time.Time
	EndTime   time.Time
}
