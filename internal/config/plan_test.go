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
		check   func(t *testing.T, plan *ExecutionPlan)
	}{
		{
			name: "valid basic plan",
			yaml: `
name: "Test Setup"
version: "1.0.0"
description: "Test description"

executors:
  shell:
    type: shell
    shell: /bin/bash

steps:
  - id: test_step
    description: "Test step"
    executor: shell
    files:
      - path: "test.sh"
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				assert.Equal(t, "Test Setup", plan.Name)
				assert.Equal(t, "1.0.0", plan.Version)
				assert.Equal(t, "Test description", plan.Description)
				assert.Len(t, plan.Executors, 1)
				assert.Len(t, plan.Steps, 1)
				assert.Equal(t, "test_step", plan.Steps[0].ID)
			},
		},
		{
			name: "plan with dependencies",
			yaml: `
name: "Test Setup"
version: "1.0.0"

executors:
  shell:
    type: shell

steps:
  - id: step1
    executor: shell
    files:
      - path: "step1.sh"
  - id: step2
    executor: shell
    depends_on: [step1]
    files:
      - path: "step2.sh"
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				assert.Len(t, plan.Steps, 2)
				assert.Equal(t, []string{"step1"}, plan.Steps[1].DependsOn)
			},
		},
		{
			name: "plan with platform-specific files",
			yaml: `
name: "Cross-platform Setup"
version: "1.0.0"

platforms:
  darwin:
    PACKAGE_MANAGER: brew
  linux:
    PACKAGE_MANAGER: apt

executors:
  shell:
    type: shell

steps:
  - id: install
    executor: shell
    files:
      - path: "install_mac.sh"
        platform: darwin
      - path: "install_linux.sh"
        platform: linux
`,
			wantErr: false,
			check: func(t *testing.T, plan *ExecutionPlan) {
				assert.Len(t, plan.Steps[0].Files, 2)
				assert.Equal(t, "darwin", plan.Steps[0].Files[0].Platform)
				assert.Equal(t, "linux", plan.Steps[0].Files[1].Platform)
			},
		},
		{
			name: "invalid yaml",
			yaml: `
name: "Test
invalid yaml
`,
			wantErr: true,
		},
		{
			name: "missing required fields",
			yaml: `
description: "Missing name and version"
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.yml")
			err := os.WriteFile(tmpFile, []byte(tt.yaml), 0644)
			require.NoError(t, err)

			// Load the plan
			plan, err := LoadExecutionPlan(tmpFile)

			if tt.wantErr {
				assert.Error(t, err)
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

func TestValidateExecutionPlan(t *testing.T) {
	tests := []struct {
		name    string
		plan    *ExecutionPlan
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid plan",
			plan: &ExecutionPlan{
				Name:    "Test",
				Version: "1.0.0",
				Executors: map[string]ExecutorConfig{
					"shell": {Type: "shell"},
				},
				Steps: []Step{
					{
						ID:       "test",
						Executor: "shell",
						Files:    []FileConfig{{Path: "test.sh"}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			plan: &ExecutionPlan{
				Version: "1.0.0",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing version",
			plan: &ExecutionPlan{
				Name: "Test",
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "undefined executor",
			plan: &ExecutionPlan{
				Name:    "Test",
				Version: "1.0.0",
				Steps: []Step{
					{
						ID:       "test",
						Executor: "undefined",
						Files:    []FileConfig{{Path: "test.sh"}},
					},
				},
			},
			wantErr: true,
			errMsg:  "undefined executor",
		},
		{
			name: "circular dependency",
			plan: &ExecutionPlan{
				Name:    "Test",
				Version: "1.0.0",
				Executors: map[string]ExecutorConfig{
					"shell": {Type: "shell"},
				},
				Steps: []Step{
					{
						ID:        "step1",
						Executor:  "shell",
						DependsOn: []string{"step2"},
						Files:     []FileConfig{{Path: "test1.sh"}},
					},
					{
						ID:        "step2",
						Executor:  "shell",
						DependsOn: []string{"step1"},
						Files:     []FileConfig{{Path: "test2.sh"}},
					},
				},
			},
			wantErr: true,
			errMsg:  "circular dependency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExecutionPlan(tt.plan)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
