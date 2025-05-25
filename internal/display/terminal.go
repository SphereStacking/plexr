package display

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SphereStacking/plexr/internal/config"
)

// TerminalDisplay implements Display for simple terminal output
type TerminalDisplay struct {
	mu            sync.Mutex
	verbose       bool
	writer        io.Writer
	lastProgress  *ExecutionProgress
	useColor      bool
	progressWidth int
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Status symbols
const (
	symbolPending   = "○"
	symbolRunning   = "◐"
	symbolCompleted = "✓"
	symbolFailed    = "✗"
	symbolSkipped   = "⊘"
)

// NewTerminalDisplay creates a new terminal display
func NewTerminalDisplay(verbose bool) *TerminalDisplay {
	return &TerminalDisplay{
		verbose:       verbose,
		writer:        os.Stdout,
		useColor:      os.Getenv("NO_COLOR") == "",
		progressWidth: 40,
	}
}

// Init initializes the display
func (td *TerminalDisplay) Init(ctx context.Context) error {
	return nil
}

// Start begins display
func (td *TerminalDisplay) Start(plan *config.ExecutionPlan) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	td.clearLine()
	fmt.Fprintf(td.writer, "%s%s Executing: %s%s",
		td.color(colorBold),
		td.color(colorBlue),
		plan.Name,
		td.color(colorReset))

	if plan.Version != "" {
		fmt.Fprintf(td.writer, " %s(v%s)%s",
			td.color(colorGray),
			plan.Version,
			td.color(colorReset))
	}
	fmt.Fprintln(td.writer)

	if plan.Description != "" && td.verbose {
		fmt.Fprintf(td.writer, "%s%s%s\n",
			td.color(colorGray),
			plan.Description,
			td.color(colorReset))
	}
	fmt.Fprintln(td.writer)

	return nil
}

// UpdateProgress updates the progress display
func (td *TerminalDisplay) UpdateProgress(progress *ExecutionProgress) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	td.lastProgress = progress

	// Move cursor up to overwrite previous output
	if !td.verbose {
		td.clearPreviousOutput(progress)
	}

	// Progress bar
	percent := 0
	if progress.TotalSteps > 0 {
		percent = (progress.CompletedSteps * 100) / progress.TotalSteps
	}

	progressBar := td.buildProgressBar(percent)
	fmt.Fprintf(td.writer, "%sProgress:%s %s %d%% (%d/%d steps)\n\n",
		td.color(colorBold),
		td.color(colorReset),
		progressBar,
		percent,
		progress.CompletedSteps,
		progress.TotalSteps)

	// Step list
	for _, step := range progress.Steps {
		td.printStep(step)
	}

	// Elapsed time
	fmt.Fprintf(td.writer, "\n%sElapsed:%s %s",
		td.color(colorGray),
		td.color(colorReset),
		progress.ElapsedTime.Round(time.Second))

	if progress.EstimatedTime > 0 && progress.CurrentStep != "" {
		fmt.Fprintf(td.writer, " %s• Remaining: ~%s%s",
			td.color(colorGray),
			progress.EstimatedTime.Round(time.Second),
			td.color(colorReset))
	}
	fmt.Fprintln(td.writer)

	return nil
}

// UpdateStep updates a step's status
func (td *TerminalDisplay) UpdateStep(stepID string, status StepStatus, message string) error {
	if td.verbose {
		td.mu.Lock()
		defer td.mu.Unlock()

		symbol := td.getStatusSymbol(status)
		color := td.getStatusColor(status)

		fmt.Fprintf(td.writer, "%s%s %s%s",
			color,
			symbol,
			stepID,
			td.color(colorReset))

		if message != "" {
			fmt.Fprintf(td.writer, ": %s", message)
		}
		fmt.Fprintln(td.writer)
	}
	return nil
}

// ShowError displays an error
func (td *TerminalDisplay) ShowError(stepID string, err error) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	fmt.Fprintf(td.writer, "\n%s%s Error in step '%s':%s\n",
		td.color(colorBold),
		td.color(colorRed),
		stepID,
		td.color(colorReset))

	fmt.Fprintf(td.writer, "%s%s%s\n",
		td.color(colorRed),
		err.Error(),
		td.color(colorReset))

	return nil
}

// ShowOutput displays output from executors
func (td *TerminalDisplay) ShowOutput(stepID string, output string) error {
	if td.verbose {
		td.mu.Lock()
		defer td.mu.Unlock()

		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			if line != "" {
				fmt.Fprintf(td.writer, "%s│%s %s\n",
					td.color(colorGray),
					td.color(colorReset),
					line)
			}
		}
	}
	return nil
}

// Finish completes the display
func (td *TerminalDisplay) Finish(success bool, summary string) error {
	td.mu.Lock()
	defer td.mu.Unlock()

	fmt.Fprintln(td.writer)

	if success {
		fmt.Fprintf(td.writer, "%s%s✓ Success!%s %s\n",
			td.color(colorBold),
			td.color(colorGreen),
			td.color(colorReset),
			summary)
	} else {
		fmt.Fprintf(td.writer, "%s%s✗ Failed!%s %s\n",
			td.color(colorBold),
			td.color(colorRed),
			td.color(colorReset),
			summary)
	}

	return nil
}

// Close closes the display
func (td *TerminalDisplay) Close() error {
	return nil
}

// Helper methods

func (td *TerminalDisplay) color(code string) string {
	if td.useColor {
		return code
	}
	return ""
}

func (td *TerminalDisplay) getStatusSymbol(status StepStatus) string {
	switch status {
	case StatusRunning:
		return symbolRunning
	case StatusCompleted:
		return symbolCompleted
	case StatusFailed:
		return symbolFailed
	case StatusSkipped:
		return symbolSkipped
	default:
		return symbolPending
	}
}

func (td *TerminalDisplay) getStatusColor(status StepStatus) string {
	switch status {
	case StatusRunning:
		return td.color(colorCyan)
	case StatusCompleted:
		return td.color(colorGreen)
	case StatusFailed:
		return td.color(colorRed)
	case StatusSkipped:
		return td.color(colorYellow)
	default:
		return td.color(colorGray)
	}
}

func (td *TerminalDisplay) buildProgressBar(percent int) string {
	filled := (td.progressWidth * percent) / 100
	empty := td.progressWidth - filled

	bar := td.color(colorGreen) + strings.Repeat("█", filled)
	bar += td.color(colorGray) + strings.Repeat("░", empty)
	bar += td.color(colorReset)

	return fmt.Sprintf("[%s]", bar)
}

func (td *TerminalDisplay) printStep(step StepProgress) {
	symbol := td.getStatusSymbol(step.Status)
	color := td.getStatusColor(step.Status)

	fmt.Fprintf(td.writer, "  %s%s %s%s",
		color,
		symbol,
		step.ID,
		td.color(colorReset))

	// Duration for completed steps
	switch {
	case step.Status == StatusCompleted && step.Duration > 0:
		fmt.Fprintf(td.writer, "  %s%6s%s",
			td.color(colorGray),
			step.Duration.Round(time.Second),
			td.color(colorReset))
	case step.Status == StatusRunning:
		fmt.Fprintf(td.writer, "  %s%6s%s",
			td.color(colorCyan),
			"...",
			td.color(colorReset))
	default:
		fmt.Fprintf(td.writer, "  %s%6s%s",
			td.color(colorGray),
			"-",
			td.color(colorReset))
	}

	// Description
	fmt.Fprintf(td.writer, "  %s", step.Description)

	// Current activity for running step
	if step.Status == StatusRunning && td.lastProgress != nil && td.lastProgress.CurrentActivity != "" {
		fmt.Fprintf(td.writer, " %s(%s)%s",
			td.color(colorCyan),
			td.lastProgress.CurrentActivity,
			td.color(colorReset))
	}

	fmt.Fprintln(td.writer)
}

func (td *TerminalDisplay) clearLine() {
	fmt.Fprintf(td.writer, "\r\033[K")
}

func (td *TerminalDisplay) clearPreviousOutput(progress *ExecutionProgress) {
	// Calculate lines to clear
	lines := 4 + len(progress.Steps) + 2 // header + steps + footer

	// Move cursor up and clear lines
	for i := 0; i < lines; i++ {
		fmt.Fprintf(td.writer, "\033[1A\033[K")
	}
}
