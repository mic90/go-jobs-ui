package main

import (
	"fmt"
	"github.com/marcusolsson/tui-go"
	"log"
	"strings"
)

type Ui struct {
	Jobs map[string]*Job
	JobsWidgets map[string]*tui.Label
	JobWidget *tui.Box
	Theme *tui.Theme
	UiInternal tui.UI
	Progress uint8
	ProgressWidget *tui.Progress
	StatusBar *tui.StatusBar
}

func NewUi() *Ui {
	// prepare theme
	theme := tui.NewTheme()
	normalStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite}
	disabledStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorWhite, Underline: tui.DecorationOn}
	activeStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorCyan}
	doneStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorGreen}
	errorStyle := tui.Style{Bg: tui.ColorBlack, Fg: tui.ColorRed}
	theme.SetStyle("label.normal", normalStyle)
	theme.SetStyle("label.active", activeStyle)
	theme.SetStyle("label.done", doneStyle)
	theme.SetStyle("label.error", errorStyle)
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
	statusBar.SetPermanentText("Progress: 0%")

	root := tui.NewVBox(progressBar, jobsWidget, statusBar)
	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetTheme(theme)
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Up", func() {
		scrollArea.Scroll(0, -1)
	})
	ui.SetKeybinding("Down", func() {
		scrollArea.Scroll(0, 1)
	})
	ui.SetKeybinding("t", func() {
		scrollArea.ScrollToTop()
	})
	ui.SetKeybinding("b", func() {
		scrollArea.ScrollToBottom()
	})

	jobs := make(map[string]*Job)
	jobsWidgets := make(map[string]*tui.Label)

	go func() {
		if err := ui.Run(); err != nil {
			panic(err)
		}
	}()

	return &Ui{jobs,jobsWidgets,  jobsList, theme, ui, 0, progressBar, statusBar}
}

func (ui *Ui) AddJob(name, description string) {
	newJob := NewJob(name, description)
	ui.Jobs[name] = newJob
	jobWidgetText := fmt.Sprintf("[%8s] %s", "", description)
	newJobWidget := tui.NewLabel(jobWidgetText)
	newJobWidget.SetStyleName("normal")
	ui.JobsWidgets[name] = newJobWidget
	ui.JobWidget.Append(newJobWidget)
}

func (ui *Ui) SetJobDisabled(name string) error {
	_, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	ui.updateJob(name, "disabled")
	return nil
}

func (ui *Ui) SetJobErrorText(name, newText string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	if newText != "" {
		job.Description = newText
	}
	job.Error = true
	ui.updateJob(name, "error")
	return nil
}

func (ui *Ui) SetJobError(name string) error {
	return ui.SetJobErrorText(name, "")
}

func (ui *Ui) SetJobDoneText(name, newText string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}
	if newText != "" {
		job.Description = newText
	}
	job.Done = true
	ui.updateJob(name,  "done")
	return nil
}

func (ui *Ui) SetJobDone(name string) error {
	return ui.SetJobDoneText(name, "")
}

func (ui *Ui) SetJobActive(name string) error {
	job, found := ui.Jobs[name]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", name)
	}

	job.Active = true
	ui.updateJob(name, "active")
	return nil
}

func (ui *Ui) setProgress(progress int) {
	ui.UiInternal.Update(func() {
		ui.ProgressWidget.SetCurrent(progress)
		ui.StatusBar.SetPermanentText(fmt.Sprintf("Progress: %d", progress))
	})
}

func (ui *Ui) setStatus(statusText string) {
	ui.UiInternal.Update(func() {
		ui.StatusBar.SetText(statusText)
	})
}

func (ui *Ui) updateJob(jobName, styleName string) error {
	job, found := ui.Jobs[jobName]
	if found == false {
		return fmt.Errorf("couldn't find job named %s", jobName)

	}
	jobWidget, found := ui.JobsWidgets[jobName]
	if found == false {
		return fmt.Errorf("couldn't find job widget named %s", jobName)
	}
	jobWidgetText := fmt.Sprintf("[%8s] %s", strings.ToUpper(styleName), job.Description)
	ui.UiInternal.Update(func() {
		jobWidget.SetText(jobWidgetText)
		jobWidget.SetStyleName(styleName)
	})
	return nil
}