package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/erlisb/go-orchestrator/task"
	"github.com/erlisb/go-orchestrator/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	fmt.Println("Starting the Cube Worker")
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	api := worker.Api{Address: host, Port: port, Worker: &w}
	go runTasks(&w)
	go w.CollectStats()

	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task %s", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently")
		}
		log.Println("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}
