package jobsui

type Job struct {
	Name string
	Description string
	Active bool
	Error bool
	Done bool
}

// NewJob creates new Job object with given name
// The name parameter is not visible and is used to search for given job
// The description parameter is what will be shown in the ui
func NewJob(name, description string) *Job {
	return &Job{name, description,false,false, false}
}