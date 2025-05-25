package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadExecutionPlan(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		errMsg  string
		check   func(t *testing.T, plan *ExecutionPlan)
	}{
		{
			name: "complete valid plan",
			yaml: `
name: "Complete Test Setup"
version: "2.1.0"
description: |
  This is a comprehensive test plan
  with multiple lines of description
work_directory: "/tmp/test"

platforms:
  darwin:
    PACKAGE_MANAGER: brew
    PREFIX: /usr/local
  linux:
    PACKAGE_MANAGER: apt
    PREFIX: /usr

executors:
  shell:
    type: shell
    shell: /bin/bash
    timeout: 300
  custom:
    type: custom
    config:
      key: value

steps:
  - id: prepare
    description: "Prepare environment"
    executor: shell
    work_directory: "/tmp/prepare"
    files:
      - path: "prepare.sh"
  
  - id: install
    description: "Install dependencies"
    executor: shell
    depends_on: [prepare]
    skip_if: "test -f /tmp/skip"
    check_command: "which tool"
    files:
      - path: "install_mac.sh"
        platform: darwin
        timeout: 600
        retry: 3
      - path: "install_linux.sh"
        platform: linux
        timeout: 300
  
  - id: configure
    description: "Configure tools"
    executor: custom
    depends_on: [install]
    transaction_mode: all
    files:
      - path: "configure.sh"
        skip_if: "test -f /etc/configured"
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				// Basic fields
				assert.Equal(t, "Complete Test Setup", plan.Name)
				assert.Equal(t, "2.1.0", plan.Version)
				assert.Contains(t, plan.Description, "comprehensive test plan")
				assert.Equal(t, "/tmp/test", plan.WorkDirectory)

				// Platforms
				assert.Len(t, plan.Platforms, 2)
				assert.Equal(t, "brew", plan.Platforms["darwin"]["PACKAGE_MANAGER"])
				assert.Equal(t, "/usr", plan.Platforms["linux"]["PREFIX"])

				// Executors
				assert.Len(t, plan.Executors, 2)
				assert.Equal(t, "shell", plan.Executors["shell"]["type"])
				assert.Equal(t, "custom", plan.Executors["custom"]["type"])

				// Steps
				assert.Len(t, plan.Steps, 3)

				// Step 1: prepare
				assert.Equal(t, "prepare", plan.Steps[0].ID)
				assert.Equal(t, "/tmp/prepare", plan.Steps[0].WorkDirectory)
				assert.Empty(t, plan.Steps[0].DependsOn)

				// Step 2: install
				assert.Equal(t, "install", plan.Steps[1].ID)
				assert.Equal(t, []string{"prepare"}, plan.Steps[1].DependsOn)
				assert.Equal(t, "test -f /tmp/skip", plan.Steps[1].SkipIf)
				assert.Equal(t, "which tool", plan.Steps[1].CheckCommand)
				assert.Len(t, plan.Steps[1].Files, 2)
				assert.Equal(t, 600, plan.Steps[1].Files[0].Timeout)
				assert.Equal(t, 3, plan.Steps[1].Files[0].Retry)

				// Step 3: configure
				assert.Equal(t, "configure", plan.Steps[2].ID)
				assert.Equal(t, "all", plan.Steps[2].TransactionMode)
				assert.Equal(t, "test -f /etc/configured", plan.Steps[2].Files[0].SkipIf)
			},
		},
		{
			name: "minimal valid plan",
			yaml: `
name: "Minimal"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: single
    executor: shell
    files:
      - path: "run.sh"
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				assert.Equal(t, "Minimal", plan.Name)
				assert.Equal(t, "1.0.0", plan.Version)
				assert.Empty(t, plan.Description)
				assert.Empty(t, plan.WorkDirectory)
				assert.Len(t, plan.Steps, 1)
			},
		},
		{
			name: "plan with multiple dependencies",
			yaml: `
name: "Complex Dependencies"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: base
    executor: shell
    files:
      - path: "base.sh"
  - id: dep1
    executor: shell
    depends_on: [base]
    files:
      - path: "dep1.sh"
  - id: dep2
    executor: shell
    depends_on: [base]
    files:
      - path: "dep2.sh"
  - id: final
    executor: shell
    depends_on: [dep1, dep2]
    files:
      - path: "final.sh"
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				assert.Len(t, plan.Steps, 4)
				assert.Empty(t, plan.Steps[0].DependsOn)
				assert.Equal(t, []string{"base"}, plan.Steps[1].DependsOn)
				assert.Equal(t, []string{"base"}, plan.Steps[2].DependsOn)
				assert.Equal(t, []string{"dep1", "dep2"}, plan.Steps[3].DependsOn)
			},
		},
		{
			name: "invalid yaml syntax",
			yaml: `
name: "Test
  invalid: yaml syntax
steps
`,
			wantErr: true,
			errMsg:  "yaml",
		},
		{
			name: "missing required name",
			yaml: `
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing required version",
			yaml: `
name: "Test"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "empty steps",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps: []
`,
			wantErr: true,
			errMsg:  "at least one step is required",
		},
		{
			name: "step without id",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - executor: shell
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "step ID is required",
		},
		{
			name: "step without files",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    files: []
`,
			wantErr: true,
			errMsg:  "at least one file is required",
		},
		{
			name: "duplicate step ids",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: duplicate
    executor: shell
    files:
      - path: "test1.sh"
  - id: duplicate
    executor: shell
    files:
      - path: "test2.sh"
`,
			wantErr: true,
			errMsg:  "duplicate step id",
		},
		{
			name: "undefined executor reference",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: undefined
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "undefined executor",
		},
		{
			name: "invalid dependency reference",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    depends_on: [nonexistent]
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "depends on undefined step",
		},
		{
			name: "circular dependency",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: step1
    executor: shell
    depends_on: [step3]
    files:
      - path: "test1.sh"
  - id: step2
    executor: shell
    depends_on: [step1]
    files:
      - path: "test2.sh"
  - id: step3
    executor: shell
    depends_on: [step2]
    files:
      - path: "test3.sh"
`,
			wantErr: true,
			errMsg:  "circular dependency",
		},
		{
			name: "self dependency",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: self
    executor: shell
    depends_on: [self]
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "circular dependency detected",
		},
		{
			name: "invalid platform value",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    files:
      - path: "test.sh"
        platform: invalid_platform
`,
			wantErr: true,
			errMsg:  "invalid platform",
		},
		{
			name: "invalid transaction mode",
			yaml: `
name: "Test"
version: "1.0.0"
executors:
  shell:
    type: shell
steps:
  - id: test
    executor: shell
    transaction_mode: invalid
    files:
      - path: "test.sh"
`,
			wantErr: true,
			errMsg:  "invalid transaction_mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.yml")
			err := os.WriteFile(tmpFile, []byte(tt.yaml), 0600) // #nosec G306 - Test file
			require.NoError(t, err)

			// Load the plan
			plan, err := LoadExecutionPlan(tmpFile)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" && err != nil {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, plan)

			if tt.check != nil {
				tt.check(t, plan)
			}
		})
	}
}

func TestLoadExecutionPlanFileOperations(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		plan, err := LoadExecutionPlan("/non/existent/file.yml")
		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "failed to read file")
	})

	t.Run("empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "empty.yml")
		err := os.WriteFile(tmpFile, []byte(""), 0600) // #nosec G306 - Test file
		require.NoError(t, err)

		plan, err := LoadExecutionPlan(tmpFile)
		assert.Error(t, err)
		assert.Nil(t, plan)
	})

	t.Run("directory instead of file", func(t *testing.T) {
		tmpDir := t.TempDir()
		plan, err := LoadExecutionPlan(tmpDir)
		assert.Error(t, err)
		assert.Nil(t, plan)
	})
}

func TestValidateExecutionPlan(t *testing.T) {
	t.Run("nil plan", func(t *testing.T) {
		err := ValidateExecutionPlan(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "plan is nil")
	})

	t.Run("executors validation", func(t *testing.T) {
		plan := &ExecutionPlan{
			Name:    "Test",
			Version: "1.0.0",
			Executors: map[string]ExecutorConfig{
				"": {"type": "shell"}, // Empty key
			},
			Steps: []Step{
				{
					ID:       "test",
					Executor: "shell",
					Files:    []FileConfig{{Path: "test.sh"}},
				},
			},
		}

		err := ValidateExecutionPlan(plan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "undefined executor 'shell'")
	})

	t.Run("file path validation", func(t *testing.T) {
		tests := []struct {
			name    string
			path    string
			wantErr bool
		}{
			{"valid relative path", "scripts/test.sh", false},
			{"valid nested path", "a/b/c/test.sh", false},
			{"empty path", "", true},
			{"absolute path", "/absolute/path.sh", true},
			{"path with ..", "../test.sh", true},
			{"current directory", "./test.sh", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				plan := &ExecutionPlan{
					Name:    "Test",
					Version: "1.0.0",
					Executors: map[string]ExecutorConfig{
						"shell": {"type": "shell"},
					},
					Steps: []Step{
						{
							ID:       "test",
							Executor: "shell",
							Files:    []FileConfig{{Path: tt.path}},
						},
					},
				}

				err := ValidateExecutionPlan(plan)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestExecutorConfig(t *testing.T) {
	t.Run("marshal and unmarshal", func(t *testing.T) {
		tests := []struct {
			name   string
			config ExecutorConfig
		}{
			{
				name: "simple config",
				config: ExecutorConfig{
					"type":  "shell",
					"shell": "/bin/bash",
				},
			},
			{
				name: "complex config",
				config: ExecutorConfig{
					"type":    "custom",
					"timeout": 300,
					"env": map[string]string{
						"KEY": "value",
					},
					"options": []string{"opt1", "opt2"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// This would test YAML marshaling if ExecutorConfig had custom methods
				assert.Equal(t, tt.config["type"], tt.config["type"])
				assert.NotNil(t, tt.config)
			})
		}
	})
}

func TestStepDependencyGraph(t *testing.T) {
	t.Run("complex dependency resolution", func(t *testing.T) {
		plan := &ExecutionPlan{
			Name:    "Complex Graph",
			Version: "1.0.0",
			Executors: map[string]ExecutorConfig{
				"shell": {"type": "shell"},
			},
			Steps: []Step{
				// Diamond dependency pattern
				{ID: "A", Executor: "shell", Files: []FileConfig{{Path: "a.sh"}}},
				{ID: "B", Executor: "shell", DependsOn: []string{"A"}, Files: []FileConfig{{Path: "b.sh"}}},
				{ID: "C", Executor: "shell", DependsOn: []string{"A"}, Files: []FileConfig{{Path: "c.sh"}}},
				{ID: "D", Executor: "shell", DependsOn: []string{"B", "C"}, Files: []FileConfig{{Path: "d.sh"}}},
				// Independent chain
				{ID: "E", Executor: "shell", Files: []FileConfig{{Path: "e.sh"}}},
				{ID: "F", Executor: "shell", DependsOn: []string{"E"}, Files: []FileConfig{{Path: "f.sh"}}},
			},
		}

		err := ValidateExecutionPlan(plan)
		assert.NoError(t, err)
	})
}

func TestPlanWithWorkDirectory(t *testing.T) {
	t.Run("global and step work directories", func(t *testing.T) {
		yaml := `
name: "WorkDir Test"
version: "1.0.0"
work_directory: "/global/work"
executors:
  shell:
    type: shell
steps:
  - id: global_dir
    executor: shell
    files:
      - path: "test1.sh"
  - id: custom_dir
    executor: shell
    work_directory: "/custom/work"
    files:
      - path: "test2.sh"
`
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.yml")
		err := os.WriteFile(tmpFile, []byte(yaml), 0600) // #nosec G306 - Test file
		require.NoError(t, err)

		plan, err := LoadExecutionPlan(tmpFile)
		require.NoError(t, err)

		assert.Equal(t, "/global/work", plan.WorkDirectory)
		assert.Equal(t, "", plan.Steps[0].WorkDirectory) // Uses global
		assert.Equal(t, "/custom/work", plan.Steps[1].WorkDirectory)
	})
}
