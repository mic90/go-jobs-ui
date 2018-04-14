package main

import "time"

func main() {
	waitChan := make(chan bool)

	ui := NewUI()
	ui.AddJob("jobA", "This is simple job description")
	ui.AddJob("jobB", "This is simple job description")
	ui.AddJob("jobC", "This is simple job description")
	ui.AddJob("jobD", "This is simple job description")
	ui.AddJob("jobE", "This is simple job description")
	ui.AddJob("jobF", "This is simple job description")

	ui.SetJobActive("jobA")

	go func() {
		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobA")
		ui.SetJobActive("jobB")

		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobB")
		ui.SetJobActive("jobC")

		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobC")
		ui.SetJobActive("jobD")

		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobD")
		ui.SetJobActive("jobE")

		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobE")
		ui.SetJobActive("jobF")

		time.Sleep(2 * time.Second)
		ui.SetJobDone("jobF")

		ui.SetStatus("All jobs done, you may close the app now")
	}()

	<-waitChan
}
