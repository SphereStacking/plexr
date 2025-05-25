package config

// ExecutionPlan represents the top-level structure of a YAML execution plan
type ExecutionPlan struct {
	Name          string                       `yaml:"name"`
	Version       string                       `yaml:"version"`
	Description   string                       `yaml:"description"`
	WorkDirectory string                       `yaml:"work_directory,omitempty"`
	Platforms     map[string]map[string]string `yaml:"platforms,omitempty"`
	Executors     map[string]ExecutorConfig    `yaml:"executors"`
	Steps         []Step                       `yaml:"steps"`
}

// ExecutorConfig represents the configuration for an executor
type ExecutorConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:",inline"`
}

// Step represents a single execution step
type Step struct {
	ID             string         `yaml:"id"`
	Description    string         `yaml:"description"`
	Executor       string         `yaml:"executor"`
	DependsOn      []string       `yaml:"depends_on,omitempty"`
	SkipIf         string         `yaml:"skip_if,omitempty"`
	CheckCommand   string         `yaml:"check_command,omitempty"`
	WorkDirectory  string         `yaml:"work_directory,omitempty"`
	Files          []FileConfig   `yaml:"files"`
	TransactionMode string        `yaml:"transaction_mode,omitempty"`
}

// FileConfig represents the configuration for a file to be executed
type FileConfig struct {
	Path     string `yaml:"path"`
	Timeout  int    `yaml:"timeout,omitempty"`
	Retry    int    `yaml:"retry,omitempty"`
	Platform string `yaml:"platform,omitempty"`
	SkipIf   string `yaml:"skip_if,omitempty"`
}
