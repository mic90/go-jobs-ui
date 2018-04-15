package jobsui

import (
	"fmt"
	"log"
	"strings"

	tui "github.com/marcusolsson/tui-go"
)

// UI is a base struct for the interface
type UI struct {
	Jobs            map[string]*Job
	JobsWidgets     map[string]*tui.Label
	JobWidget       *tui.Box
	Theme           *tui.Theme
	UIInternal      tui.UI
	ProgressWidget  *tui.Progress
	StatusBarWidget *tui.StatusBar
	JobsDone        int
	ScrollPos       *int
}

// NewUI creates new interface with empty jobs list
// The ui event loop is started immediately in a seperate gorotuine
func NewUI() *UI {
	// prepare theme
	theme := tui.NewTheme()
	normalStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite}
	disabledStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite, Underline: tui.DecorationOn}
	activeStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorCyan}
	doneStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorGreen}
	failedStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorRed}
	theme.SetStyle("label.normal", normalStyle)
	theme.SetStyle("label.active", activeStyle)
	theme.SetStyle("label.done", doneStyle)
	theme.SetStyle("label.failed", failedStyle)
	theme.SetStyle("label.disabled", disabledStyle)

	// prepare ui with progress bar, list of jobs and status bar
	progressBar := tui.NewProgress(100)
	progressBar.SetCurrent(0)
	progressBar.SetSizePolicy(tui.Minimum, tui.Minimum)

	jobsList := tui.NewVBox()
	jobsList.SetFocused(true)
	scrollArea := tui.NewScrollArea(jobsList)
	jobsWidget := tui.NewVBox(scrollArea)
	jobsWidget.SetBorder(true)
	jobsWidget.SetSizePolicy(tui.Expanding, tui.Expanding)

	statusBar := tui.NewStatusBar("")
	statusBar.SetPermanentText("Progress: 0 %")

	root := tui.NewVBox(progressBar, jobsWidget, statusBar)
	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	scrollPos := new(int)
	*scrollPos = 0

	ui.SetTheme(theme)
	ui.SetKeybinding("Up", func() {
		if *scrollPos > 0 {
			scrollArea.Scroll(0, -1)
			*scrollPos--
		}
	})
	ui.SetKeybinding("Down", func() {
		scrollArea.Scroll(0, 1)
		*scrollPos++
	})
	ui.SetKeybinding("t", func() {
		scrollArea.ScrollToTop()
		*scrollPos = 0
	})

	jobs := make(map[string]*Job)
	jobsWidgets := make(map[string]*tui.Label)

	go func() {
		if err := ui.Run(); err != nil {
			panic(err)
		}
	}()

	return &UI{jobs, jobsWidgets, jobsList, theme, ui, progressBar, statusBar, 0, scrollPos}
}

// AddJob adds job to the ui list with given name and description
// The name is used only internally for lookup operations
// The description is what will be visible in the ui
func (ui *UI) AddJob(name, description string) {
	newJob := NewJob(name, description)
	ui.Jobs[name] = newJob
	jobWidgetText := prepareJobDescription("", description)
	newJobWidget := tui.NewLabel(jobWidgetText)
	newJobWidget.SetStyleName("normal")
	ui.JobsWidgets[name] = newJobWidget
	ui.JobWidget.Append(newJobWidget)
}

// SetJobDisabled sets given job state to disabled.
// Disabled job will have its status changed in the ui to [DISABLED]
// The stylesheet will be changed to underline white on black text
func (ui *UI) SetJobDisabled(name string) error {
	_, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	ui.updateJob(name, "disabled")
	ui.increaseProgress()
	return nil
}

// SetJobFailedText sets given job state to failed and adds message to the current description with format [description]:[message]
// Failed jobs will have its status changed in the ui to [FAILED]
// The stylesheet will be changed to red on black text
func (ui *UI) SetJobFailedText(name, message string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	if message != "" {
		job.Description = fmt.Sprintf("%s : %s", job.Description, message)
	}
	job.Error = true
	ui.updateJob(name, "failed")
	return nil
}

// SetJobFailed sets given job state to failed.
func (ui *UI) SetJobFailed(name string) error {
	return ui.SetJobFailedText(name, "")
}

// SetJobDoneText sets given job state to done and adds message to the current description with format [description]:[message]
// Done jobs will have its status changed in the ui to [DONE]
// The stylesheet will be changed to green on black text
func (ui *UI) SetJobDoneText(name, message string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	if message != "" {
		job.Description = fmt.Sprintf("%s : %s", job.Description, message)
	}
	job.Done = true
	ui.updateJob(name, "done")
	ui.increaseProgress()
	return nil
}

// SetJobDone sets given job state to done.
func (ui *UI) SetJobDone(name string) error {
	return ui.SetJobDoneText(name, "")
}

// SetJobActive sets given job state to active
// Active jobs will have its status changed in the ui to [ACTIVE]
// The stylesheet will be changed to cyan on balck text
func (ui *UI) SetJobActive(name string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}

	job.Active = true
	ui.updateJob(name, "active")
	return nil
}

// SetProgress sets progress of the whole operation to given value
// The value must be betweeon 0 and 100
// This will update the progress bar value and status bar text
func (ui *UI) SetProgress(progress int) {
	ui.UIInternal.Update(func() {
		ui.ProgressWidget.SetCurrent(progress)
		ui.StatusBarWidget.SetPermanentText(prepareProgressText(progress))
	})
}

// SetStatus sets temporary status bra text
func (ui *UI) SetStatus(statusText string) {
	ui.UIInternal.Update(func() {
		ui.StatusBarWidget.SetText(statusText)
	})
}

func (ui *UI) updateJob(jobName, styleName string) error {
	job, found := ui.Jobs[jobName]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", jobName)

	}
	jobWidget, found := ui.JobsWidgets[jobName]
	if found == false {
		return fmt.Errorf("couldn't find job widget named %s", jobName)
	}
	jobWidgetText := prepareJobDescription(styleName, job.Description)
	ui.UIInternal.Update(func() {
		jobWidget.SetText(jobWidgetText)
		jobWidget.SetStyleName(styleName)
	})
	return nil
}

func (ui *UI) increaseProgress() {
	jobsCount := len(ui.Jobs)
	ui.JobsDone++
	if ui.JobsDone == jobsCount {
		ui.SetProgress(100)
	} else {
		progressValue := 100 / jobsCount * ui.JobsDone
		ui.SetProgress(progressValue)
	}
}

func prepareJobDescription(status, text string) string {
	return fmt.Sprintf("[%8s] %s", strings.ToUpper(status), text)
}

func prepareProgressText(value int) string {
	return fmt.Sprintf("Progress: %d %%", value)
}
