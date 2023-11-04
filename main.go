package main

import (
	"fmt"
	"time"

	"github.com/erlisb/go-orchestrator/task"
	"github.com/erlisb/go-orchestrator/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("Starting task")
	w.AddTask(t)

	result := w.RunTask()

	if result.Error != nil {
		panic(result.Error)
	}

	

	fmt.Printf("task %s is running in container %s\n", t.ID, result.ContainerId)

	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 30)

	fmt.Printf("stopping task %s\n", t.ID)

	t.ContainerId = result.ContainerId

	t.State = task.Completed
	fmt.Println("Task ID: ", t.ID)

	w.AddTask(t)

	result = w.RunTask()

	if result.Error != nil {
		panic(result.Error)
	}
}
