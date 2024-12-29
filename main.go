package main

import (
	"fmt"
	"log"
	"sync"

	drivers "github.com/PratikKumar125/go-queue/drivers/redis"
	"github.com/PratikKumar125/go-queue/handler"
)

func greet() (string, error) {
    fmt.Println("Hello, World!")
    return "Hello World", nil
}

func greet2() (any, error) {
    // fmt.Println("Hello, World 2222222! I HAVE BEEN DELAYEDDDDD")
    return "", fmt.Errorf("some error occured")
}

func main() {
    connectionMap := drivers.ConnectionInputStruct{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
        Queue: "prateek",
    }
	rClient := drivers.InitRedisQueue(connectionMap)

	var wg sync.WaitGroup

	rClient.StartWorkers(3, &wg) // concurrency limit for number of workers
	go rClient.StartQueueListener(&wg)
    go rClient.StartDelayedQueueListener(&wg)
	go rClient.StartResultProcessor(&wg)
    wg.Add(3)

    greetJobInputs := handler.CreateInputStruct{
        JobName: "greet",
        Queue: "prateek", 
        Lookup: greet,
    }
    greetJobInputs2 := handler.CreateInputStruct{
        JobName: "greet2",
        Queue: "prateek", 
        Lookup: greet2,
    }


    jobsHandler := handler.NewHandler()
    jobsHandler.CreateJob(greetJobInputs)
    jobsHandler.CreateJob(greetJobInputs2)

    // dispatch1 := handler.DispatchJobStruct{
    //     JobName: "greet",
    //     Data: map[string]any{"name": "prateek", "email": "prateek378@gmail.com"},
    //     Queue: "prateek",
    //     Tries: 0,
    //     Delay: 30,
    // }
    
    // _, err := rClient.DispatchJob(dispatch1)
    // if err != nil {
    //     log.Println(err)
    // }
    
    dispatch2 := handler.DispatchJobStruct{
        JobName: "greet2",
        Data: map[string]any{"name": "prateek", "email": "prateek378@gmail.com"},
        Queue: "prateek",
        Tries: 2,
        Delay: -1,
    }
    
    _, err2 := rClient.DispatchJob(dispatch2)
    if err2 != nil {
        log.Println(err2)
    }

	wg.Wait()


    // originals := storage.DiskStruct{
	// 	Bucket: "testing-media-services",
	// 	Region: "us-east-1",
	// 	Profile: "ankush_prasoon_account",
	// }

    // disks := map[string]storage.DiskStruct{
    //     "originals": originals,
    // }

    // st := storage.InitStorage("originals", disks)
    // items, err := st.GetBucketItems()
	// fmt.Println(items,err, "<<<<<<BUCKET ITEMS")
}
