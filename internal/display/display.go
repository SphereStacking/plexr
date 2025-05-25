package display

import (
	"context"
	"fmt"
	"time"

	"github.com/SphereStacking/plexr/internal/config"
	"github.com/SphereStacking/plexr/internal/core"
)

// DisplayMode represents the output mode
type DisplayMode string

const (
	ModeSimple      DisplayMode = "simple"      // Simple text output
	ModeInteractive DisplayMode = "interactive" // Future TUI mode
	ModeJSON        DisplayMode = "json"        // JSON output for CI/CD
)

// StepStatus represents the current status of a step
type StepStatus string

const (
	StatusPending   StepStatus = "pending"
	StatusRunning   StepStatus = "running"
	StatusCompleted StepStatus = "completed"
	StatusFailed    StepStatus = "failed"
	StatusSkipped   StepStatus = "skipped"
)

// ExecutionProgress represents the current execution state
type ExecutionProgress struct {
	Plan            *config.ExecutionPlan
	TotalSteps      int
	CompletedSteps  int
	CurrentStep     string
	StartTime       time.Time
	ElapsedTime     time.Duration
	EstimatedTime   time.Duration
	Steps           []StepProgress
	CurrentActivity string
	Error           error
}

// StepProgress represents progress of a single step
type StepProgress struct {
	ID           string
	Description  string
	Status       StepStatus
	StartTime    *time.Time
	Duration     time.Duration
	Error        error
	Dependencies []string
	Files        []FileProgress
}

// FileProgress represents progress of a single file execution
type FileProgress struct {
	Path     string
	Status   StepStatus
	Error    error
	Output   string
	Duration time.Duration
}

// Display is the interface for rendering execution progress
type Display interface {
	// Initialize the display
	Init(ctx context.Context) error

	// Start displaying execution
	Start(plan *config.ExecutionPlan) error

	// Update progress
	UpdateProgress(progress *ExecutionProgress) error

	// Update step status
	UpdateStep(stepID string, status StepStatus, message string) error

	// Show error
	ShowError(stepID string, err error) error

	// Show output from executors
	ShowOutput(stepID string, output string) error

	// Finish execution
	Finish(success bool, summary string) error

	// Close the display
	Close() error
}

// Factory creates a new Display instance based on mode
func NewDisplay(mode DisplayMode, verbose bool) Display {
	switch mode {
	case ModeInteractive:
		// Future: return NewTUIDisplay(verbose)
		return NewTerminalDisplay(verbose)
	case ModeJSON:
		// Future: return NewJSONDisplay()
		return NewTerminalDisplay(verbose)
	default:
		return NewTerminalDisplay(verbose)
	}
}

// ProgressTracker provides progress information to displays
type ProgressTracker struct {
	plan      *config.ExecutionPlan
	state     *core.ExecutionState
	startTime time.Time
	display   Display
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(plan *config.ExecutionPlan, state *core.ExecutionState, display Display) *ProgressTracker {
	return &ProgressTracker{
		plan:      plan,
		state:     state,
		startTime: time.Now(),
		display:   display,
	}
}

// SetState updates the state reference
func (pt *ProgressTracker) SetState(state *core.ExecutionState) {
	pt.state = state
}

// Start begins tracking
func (pt *ProgressTracker) Start() error {
	return pt.display.Start(pt.plan)
}

// Update updates the current progress
func (pt *ProgressTracker) Update() error {
	progress := pt.buildProgress()
	return pt.display.UpdateProgress(progress)
}

// StepStarted notifies that a step has started
func (pt *ProgressTracker) StepStarted(stepID string) error {
	return pt.display.UpdateStep(stepID, StatusRunning, "Starting...")
}

// StepCompleted notifies that a step has completed
func (pt *ProgressTracker) StepCompleted(stepID string, duration time.Duration) error {
	return pt.display.UpdateStep(stepID, StatusCompleted, "")
}

// StepFailed notifies that a step has failed
func (pt *ProgressTracker) StepFailed(stepID string, err error) error {
	return pt.display.ShowError(stepID, err)
}

// StepSkipped notifies that a step was skipped
func (pt *ProgressTracker) StepSkipped(stepID string, reason string) error {
	return pt.display.UpdateStep(stepID, StatusSkipped, reason)
}

// Output shows output from a step
func (pt *ProgressTracker) Output(stepID string, output string) error {
	return pt.display.ShowOutput(stepID, output)
}

// Finish completes tracking
func (pt *ProgressTracker) Finish(success bool) error {
	summary := pt.buildSummary()
	return pt.display.Finish(success, summary)
}

// buildProgress builds the current progress state
func (pt *ProgressTracker) buildProgress() *ExecutionProgress {
	elapsed := time.Since(pt.startTime)

	completedSteps := 0
	currentStep := ""

	if pt.state != nil {
		completedSteps = len(pt.state.CompletedSteps)
		currentStep = pt.state.CurrentStep
	}

	progress := &ExecutionProgress{
		Plan:           pt.plan,
		TotalSteps:     len(pt.plan.Steps),
		CompletedSteps: completedSteps,
		CurrentStep:    currentStep,
		StartTime:      pt.startTime,
		ElapsedTime:    elapsed,
		Steps:          make([]StepProgress, 0, len(pt.plan.Steps)),
	}

	// Build step progress
	for _, step := range pt.plan.Steps {
		sp := StepProgress{
			ID:           step.ID,
			Description:  step.Description,
			Dependencies: step.DependsOn,
		}

		// Determine status
		completed := false
		if pt.state != nil {
			for _, completedID := range pt.state.CompletedSteps {
				if completedID == step.ID {
					completed = true
					break
				}
			}
		}

		switch {
		case completed:
			sp.Status = StatusCompleted
		case pt.state != nil && step.ID == pt.state.CurrentStep:
			sp.Status = StatusRunning
		default:
			sp.Status = StatusPending
		}

		progress.Steps = append(progress.Steps, sp)
	}

	// Estimate remaining time
	if progress.CompletedSteps > 0 {
		avgStepTime := elapsed / time.Duration(progress.CompletedSteps)
		remainingSteps := progress.TotalSteps - progress.CompletedSteps
		progress.EstimatedTime = avgStepTime * time.Duration(remainingSteps)
	}

	return progress
}

// buildSummary builds execution summary
func (pt *ProgressTracker) buildSummary() string {
	elapsed := time.Since(pt.startTime)
	completed := 0
	if pt.state != nil {
		completed = len(pt.state.CompletedSteps)
	}
	total := len(pt.plan.Steps)

	return fmt.Sprintf("Completed %d/%d steps in %s", completed, total, elapsed.Round(time.Second))
}
