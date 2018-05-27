package jobsui

import (
	"fmt"
	"log"
	"sync"

	tui "github.com/marcusolsson/tui-go"
)

// UI is a base struct for the interface
type UI struct {
	Jobs            map[string]*Job
	jobWidget       *tui.Box
	theme           *tui.Theme
	uiInternal      tui.UI
	progressWidget  *tui.Progress
	statusBarWidget *tui.StatusBar
	jobsDone        int
	scrollPos       *int
	mutex           *sync.Mutex
	JobsDone        chan bool
}

// NewUI creates new interface with empty jobs list
// The ui event loop is started immediately in a seperate gorotuine
func NewUI() *UI {
	// prepare theme
	theme := tui.NewTheme()
	normalStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite}
	disabledStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite, Underline: tui.DecorationOn}
	activeStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite}
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
	statusBar.SetPermanentText(prepareProgressText(0))

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
	// run ui event processing in separate goroutine
	go func() {
		if err := ui.Run(); err != nil {
			panic(err)
		}
	}()

	return &UI{make(map[string]*Job), jobsList, theme, ui, progressBar, statusBar, 0, scrollPos, &sync.Mutex{}, make(chan bool)}
}

// AddJob adds job to the ui list with given name and description
// The name is used only internally for lookup operations
// The description is what will be visible in the ui
func (ui *UI) AddJob(name, description string) {
	widget := tui.NewLabel("text string")
	newJob := NewJob(name, description, widget)
	ui.jobWidget.Append(widget)
	ui.Jobs[name] = newJob
}

// SetJobState sets job to given state
func (ui *UI) SetJobState(name string, state JobState) error {
	return ui.SetJobStateWithInfo(name, state, "")
}

// SetJobStateWithInfo sets job to given state, with additonal info text appended to its description
func (ui *UI) SetJobStateWithInfo(name string, state JobState, infoText string) error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()

	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	ui.updateJobStateAndProgress(job, state, infoText)
	return nil
}

// SetJobProgress sets job to active state and set its progress to given value
// If value is >= 100 the job state will automatically change to Done
func (ui *UI) SetJobProgress(name string, progress int) error {
	return ui.SetJobProgressWithInfo(name, progress, "")
}

// SetJobProgressWithInfo sets job to active state and set its progress to given value, with additional text appended to its description
// If value is >= 100 the job state will automatically change to Done
func (ui *UI) SetJobProgressWithInfo(name string, progress int, infoText string) error {
	ui.mutex.Lock()
	defer ui.mutex.Unlock()

	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}

	if progress >= 100 {
		if job.State != Done {
			ui.updateJobStateAndProgress(job, Done, infoText)
		}
		return nil
	}

	ui.uiInternal.Update(func() {
		if infoText == "" {
			job.SetProgress(progress)
		} else {
			job.SetProgressWithInfo(progress, infoText)
		}
	})
	return nil
}

// SetStatus sets temporary status bra text
func (ui *UI) SetStatus(statusText string) {
	ui.uiInternal.Update(func() {
		ui.statusBarWidget.SetText(statusText)
	})
}

func (ui *UI) updateJobStateAndProgress(job *Job, state JobState, infoText string) {
	ui.uiInternal.Update(func() {
		if infoText == "" {
			job.SetState(state)
		} else {
			job.SetStateWithInfo(state, infoText)
		}
	})
	if state == Done {
		ui.increaseProgress()
	}
}

func (ui *UI) increaseProgress() {
	jobsCount := len(ui.Jobs)
	ui.jobsDone++
	if ui.jobsDone == jobsCount {
		ui.setProgress(100)
		ui.JobsDone <- true
	} else {
		progressValue := 100 / jobsCount * ui.jobsDone
		ui.setProgress(progressValue)
	}
}

func (ui *UI) setProgress(progress int) {
	ui.uiInternal.Update(func() {
		ui.progressWidget.SetCurrent(progress)
		ui.statusBarWidget.SetPermanentText(prepareProgressText(progress))
	})
}

func prepareProgressText(value int) string {
	return fmt.Sprintf("Progress: %d %%", value)
}
