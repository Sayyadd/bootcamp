package main

import (
	"container/list"
	"fmt"
)

type Set map[string]struct{}

func (s Set) Add(element string) {
	s[element] = struct{}{}
}

func (s Set) Remove(element string) {
	delete(s, element)
}

func (s Set) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

func (s Set) Len() int {
	return len(s)
}

// Job represents a job with an ID and a list of dependencies
type Job struct {
	ID           string
	Dependencies Set
}

// JobScheduler handles scheduling and processing of jobs
type JobScheduler struct {
	// your data structures
	jobs       map[string]*Job
	jobQueue   *list.List
	dependents map[string]Set
}

// NewJobScheduler initializes a new JobScheduler
func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		jobs:       make(map[string]*Job),
		jobQueue:   list.New(),
		dependents: make(map[string]Set),
	}
}

// AddJob adds a new job to the scheduler
// Parameters:
// - id: the unique identifier for the job
// - dependencies: a list of job IDs that this job depends on
func (js *JobScheduler) AddJob(id string, dependencies []string) {
	depSet := make(Set)
	for _, dep := range dependencies {
		depSet.Add(dep)
	}

	job := &Job{ID: id, Dependencies: depSet}
	js.jobs[id] = job
	js.jobQueue.PushBack(job)

	for _, dep := range dependencies {
		if js.dependents[dep] == nil {
			js.dependents[dep] = make(Set)
		}
		js.dependents[dep].Add(id)
	}
}

// RemoveJob removes a job from the scheduler
// Parameters:
// - id: the unique identifier of the job to remove
func (js *JobScheduler) RemoveJob(id string) {
	job, exists := js.jobs[id]
	if !exists {
		return
	}

	for e := js.jobQueue.Front(); e != nil; e = e.Next() {
		if e.Value.(*Job).ID == id {
			js.jobQueue.Remove(e)
			break
		}
	}

	delete(js.jobs, id)

	for dep := range job.Dependencies {
		if dependentSet, ok := js.dependents[dep]; ok {
			dependentSet.Remove(id)
			// Clean up empty sets
			if dependentSet.Len() == 0 {
				delete(js.dependents, dep)
			}
		}
	}

	delete(js.dependents, id)
}

// AddDependency adds a dependency to a job
// Parameters:
// - jobID: the ID of the job to add a dependency to
// - dependencyID: the ID of the job that is a dependency
func (js *JobScheduler) AddDependency(jobID, dependencyID string) {
	job, exists := js.jobs[jobID]
	if !exists {
		return
	}
	job.Dependencies.Add(dependencyID)
	if js.dependents[dependencyID] == nil {
		js.dependents[dependencyID] = make(Set)
	}
	js.dependents[dependencyID].Add(jobID)
}

// RemoveDependency removes a dependency from a job
// Parameters:
// - jobID: the ID of the job to remove a dependency from
// - dependencyID: the ID of the dependency to remove
func (js *JobScheduler) RemoveDependency(jobID, dependencyID string) {
	job, exists := js.jobs[jobID]
	if !exists {
		return
	}
	job.Dependencies.Remove(dependencyID)
	if js.dependents[dependencyID] != nil {
		js.dependents[dependencyID].Remove(jobID)
	}
}

// GetNextJob retrieves the next job to be processed based on dependencies
// Returns:
// - *Job: the next job to be processed or nil if no jobs are available
func (js *JobScheduler) GetNextJob() *Job {
	for e := js.jobQueue.Front(); e != nil; e = e.Next() {
		job := e.Value.(*Job)
		if job.Dependencies.Len() == 0 {
			js.jobQueue.Remove(e)
			return job
		}
	}
	return nil
}

// ProcessJob processes the next job in the queue
func (js *JobScheduler) ProcessJob() {
	job := js.GetNextJob()
	if job == nil {
		fmt.Println("No job to process")
		return
	}
	fmt.Println("Processing job:", job.ID)
	delete(js.jobs, job.ID)

	// Update dependents and add to readyQueue if needed
	for dependentID := range js.dependents[job.ID] {
		dependentJob := js.jobs[dependentID]
		dependentJob.Dependencies.Remove(job.ID)
	}
	delete(js.dependents, job.ID)
}

// DisplayJobQueue displays the current job queue
func (js *JobScheduler) DisplayJobQueue() {
	fmt.Println("Current job queue:")
	for element := js.jobQueue.Front(); element != nil; element = element.Next() {
		job := element.Value.(*Job)
		fmt.Println("Job ID:", job.ID, "Dependencies:", job.Dependencies.Len())
	}
}

// DisplayJobs displays all jobs with their details
func (js *JobScheduler) DisplayJobs() {
	fmt.Println("All jobs:")
	for id, job := range js.jobs {
		dependencies := make([]string, 0, job.Dependencies.Len())
		for dep := range job.Dependencies {
			dependencies = append(dependencies, dep)
		}
		fmt.Printf("Job ID: %s, Dependencies: %v\n", id, dependencies)
	}
}
