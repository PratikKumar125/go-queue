package handler

import (
	"fmt"
	"log"
	"sync"
)

type CreateJobStruct struct {
	JobName string `json:"job"`
	Queue   string `json:"queue"`
	Lookup any	`json:"lookup"`
}

type DispatchJobStruct struct {
	JobName string `json:"job"`
	Data map[string]any `json:"data"`
	Queue   string `json:"queue"`
	Tries int `json:"tries"`
	Delay int `json:"delay"`
}

type HandlerStruct struct {
	jobs map[string]map[string]any
}

var (
	handlerInstance *HandlerStruct
	once            sync.Once
)

func NewHandler() *HandlerStruct{

	once.Do(func() {
		handlerInstance = &HandlerStruct{
			jobs: make(map[string]map[string]any),
		}
	})
	return handlerInstance
}

func (handler *HandlerStruct) InitQueueJobs(queue string) {
	handler.jobs[queue] = make(map[string]any)
}

func (handler *HandlerStruct) CreateJob(inputs CreateJobStruct) (bool, error){
	if inputs.JobName == "" || inputs.Queue == "" {
		return false, fmt.Errorf("job name and queue are required")
	}

	val, ok := handler.jobs[inputs.Queue]
	if ok {
		val[inputs.JobName] = inputs.Lookup
		log.Println("JOB", inputs.JobName, "ADDED SUCCESSFULLY")
		return true, nil
	}
	return false, fmt.Errorf("no queue initialized yet in the jobs")
}

func (handler *HandlerStruct) GetJobLookup(queue string, jobName string) (func() (error), error) {
	log.Println(handler.jobs, queue)
	val, ok := handler.jobs[queue][jobName]
	if !ok {
		return nil, fmt.Errorf("no lookup method found for the job %s", jobName)
	}

	// assert that the value as a function
	lookup, ok := val.(func() (error))
	if !ok {
		return nil, fmt.Errorf("lookup method for job '%s' is not a valid function", jobName)
	}

	return lookup, nil
}

