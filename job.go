package jobsui

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

// JobState represents the job state enum type
type JobState int

const (
	// Idle default state on created jobs
	Idle JobState = iota
	// Active states that the job is currently running
	Active
	// Skipped job was never run and will not be
	Skipped
	// Done the job was run and finished without error
	Done
	// Error the job was run and finished with error
	Error
)

// Job represents a simple job with name, description and status
type Job struct {
	Name        string
	Description string
	Progress    int
	State       JobState
	widget      *tui.Label
}

// NewJob creates new Job object with given name
// The name parameter is not visible and is used to search for given job
// The description parameter is what will be shown in the ui
func NewJob(name, description string, widget *tui.Label) *Job {
	widgetText := preapreJobDescriptionNoStatus(description)
	widget.SetStyleName("normal")
	widget.SetText(widgetText)
	return &Job{name, description, 0, Idle, widget}
}

// SetState changes this job state to given one. The job will change it's text color and state text
func (job *Job) SetState(state JobState) {
	job.SetStateWithInfo(state, "")
}

// SetStateWithInfo changes this job state to the given one. The job will change it's text color and state text.
// Also additional text will be appended to the job description after ':' sign
func (job *Job) SetStateWithInfo(state JobState, additionalText string) {
	job.State = state
	switch state {
	case Idle:
		job.setStateText("", additionalText)
		job.widget.SetStyleName("normal")
	case Active:
		job.setStateText(" -> ", additionalText)
		job.widget.SetStyleName("active")
	case Skipped:
		job.setStateText("SKIP", additionalText)
		job.widget.SetStyleName("disabled")
	case Done:
		job.setStateText("DONE", additionalText)
		job.widget.SetStyleName("done")
	case Error:
		job.setStateText("FAIL", additionalText)
		job.widget.SetStyleName("failed")
	}
}

// SetProgress changes this job progress to given value.
// The progress value will be shown in the status text with % sign. for example [  40%  ]
func (job *Job) SetProgress(progress int) {
	if job.State == Done {
		return
	}

	if progress > 100 {
		progress = 100
	} else if progress < 0 {
		progress = 0
	}
	job.SetProgressWithInfo(progress, "")
}

// SetProgressWithInfo changes this job progress to given value.
// The progress value will be shown in the status text with % sign. for example [  40%  ]
func (job *Job) SetProgressWithInfo(progress int, infoText string) {
	job.Progress = progress
	job.widget.SetText(prepareJobProgressDescriptionWithInfo(progress, job.Description, infoText))
}

func (job *Job) setStateText(stateText, infoText string) {
	if infoText == "" {
		job.widget.SetText(prepareJobDescription(stateText, job.Description))
	} else {
		job.widget.SetText(preapreJobDescriptionWithInfo(stateText, job.Description, infoText))
	}
}

func preapreJobDescriptionNoStatus(description string) string {
	return prepareJobDescription("", description)
}

func prepareJobDescription(status, description string) string {
	return fmt.Sprintf("[%4s] %s", status, description)
}

func preapreJobDescriptionWithInfo(status, description, infoText string) string {
	return fmt.Sprintf("[%4s] %s : %s", status, description, infoText)
}

func prepareJobProgressDescription(progress int, text string) string {
	return fmt.Sprintf("[%4d%%] %s", progress, text)
}

func prepareJobProgressDescriptionWithInfo(progress int, description, infoText string) string {
	return fmt.Sprintf("[%3d%%] %s : %s", progress, description, infoText)
}
