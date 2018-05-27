package main

import (
	"fmt"
	"time"

	"github.com/mic90/go-jobs-ui"
)

func main() {
	ui := jobsui.NewUI()
	ui.AddJob("jobA", "This is simple job description")
	ui.AddJob("jobB", "This is simple job description")
	ui.AddJob("jobC", "This is simple job description")
	ui.AddJob("jobD", "This is simple job description")
	ui.AddJob("jobE", "This is simple job description")
	ui.AddJob("jobF", "This is simple job description")
	ui.AddJob("jobProgress", "This is job with progress indicator")
	ui.AddJob("anotherJobProgress", "This is another job with progress indicator")

	ui.SetJobState("jobA", jobsui.Active)

	go func() {
		time.Sleep(2 * time.Second)
		ui.SetJobState("jobA", jobsui.Done)
		ui.SetJobState("jobB", jobsui.Active)

		time.Sleep(2 * time.Second)
		ui.SetJobState("jobB", jobsui.Done)
		ui.SetJobState("jobC", jobsui.Active)

		time.Sleep(2 * time.Second)
		ui.SetJobState("jobC", jobsui.Done)
		ui.SetJobState("jobD", jobsui.Active)

		time.Sleep(2 * time.Second)
		ui.SetJobState("jobD", jobsui.Done)
		ui.SetJobState("jobE", jobsui.Active)

		time.Sleep(2 * time.Second)
		ui.SetJobState("jobE", jobsui.Done)
		ui.SetJobState("jobF", jobsui.Active)

		time.Sleep(2 * time.Second)
		ui.SetJobState("jobF", jobsui.Done)
	}()

	go func() {
		progress := 0
		for {
			ui.SetJobProgress("jobProgress", progress)
			time.Sleep(100 * time.Millisecond)
			progress++
		}
	}()

	go func() {
		progress := 0
		for {
			ui.SetJobProgressWithInfo("anotherJobProgress", progress, fmt.Sprintf("current progress: %d %%", progress))
			time.Sleep(100 * time.Millisecond)
			progress++
		}
	}()

	<-ui.JobsDone
	ui.SetStatus("All jobs done, you may now close the app")

	time.Sleep(5 * time.Second)
}
