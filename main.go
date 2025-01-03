package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/PratikKumar125/go-queue/queue"
)

func greet(data map[string]any) (error) { 
    fmt.Println(data, "<<<<<DATA RECEIVEDD")
    return fmt.Errorf("dummy error")
}

func main() {
    connectionMap := queue.RedisConnectionInputStruct{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
        Queue: "prateek",
    }
	rClient := queue.NewRedisQueue(connectionMap)
	handler := queue.NewHandler()

	var queueWg sync.WaitGroup
    rClient.StartConnection(&queueWg)
	rClient.StartWorkers(3, &queueWg)

	// Add a job
	greetJobInputs := queue.CreateJobStruct{
		JobName: "greet",
		Queue:   "prateek",
		Lookup:  greet,
	}
	handler.CreateJob(greetJobInputs)

    dispatchJob1 := queue.DispatchJobStruct{
        JobName: "greet",
        Data: map[string]any{
            "name": "Prateek",
            "email": "prateek@378.com",
        },
        Queue: "prateek",
        Tries: 2,
        Delay: -1,
    }
    _, dispatch1Err := rClient.DispatchJob(dispatchJob1)
    if dispatch1Err != nil {
        log.Println("Error occured", dispatch1Err)
    }

	queueWg.Wait()
	fmt.Println("All tasks completed")
}