package drivers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/PratikKumar125/go-queue/queue/handler"
)

func RedisWorker(wg *sync.WaitGroup, workerId int, rClient *RedisClient, jobHandler *handler.HandlerStruct) {
	log.Println("Worker", workerId, "Started")
	defer wg.Done()
	for job := range rClient.Jobs { 
		log.Println("Received job:", job, "By worker:", workerId)
		
		var data handler.DispatchJobStruct
		json.Unmarshal([]byte(job), &data)
		
		time.Sleep(5 * time.Second)
		lookup, err := jobHandler.GetJobLookup(rClient.Queue, data.JobName)
		if err != nil {
			log.Println("No lookup method found for job", data.JobName, err)
		} else {
			result, err := lookup()
			if err != nil {
				log.Println("Error in running job", data.JobName, err)
				if data.Data["retryCount"].(float64) < float64(data.Tries) {
					data.Data["retryCount"] = data.Data["retryCount"].(float64) + 1
					rClient.AddJobToQueue(data.Queue, data)
				}
				continue
			} else {
				rClient.Results <- fmt.Sprintf("Done processing job %s with result as %s", data.JobName, result)
			}
		}
	}
}

func (rClient *RedisClient) StartWorkers(concurrency int, wg *sync.WaitGroup) {
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go RedisWorker(wg, i, rClient, rClient.handler)
	}
}