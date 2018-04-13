package main

import "time"

func main() {
	waitChan := make(chan bool)

	ui := NewUi()
	ui.AddJob("jobA", "This is simple job description")
	ui.AddJob("jobB", "This is simple job description")
	ui.AddJob("jobC", "This is simple job description")
	ui.AddJob("jobD", "This is simple job description")

	ui.SetJobActive("jobA")

	go func() {
		time.Sleep(2 * time.Second)
		ui.setProgress(25)
		ui.SetJobDone("jobA")
		ui.SetJobActive("jobB")

		time.Sleep(2 * time.Second)
		ui.setProgress(50)
		ui.SetJobDone("jobB")
		ui.SetJobActive("jobC")

		time.Sleep(2 * time.Second)
		ui.setProgress(75)
		ui.SetJobDone("jobC")
		ui.SetJobActive("jobD")

		time.Sleep(2 * time.Second)
		ui.setProgress(100)
		ui.SetJobDone("jobD")
		
		ui.setStatus("All jobs done, you may close the app now")
	}()

	<- waitChan
}