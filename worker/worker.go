package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/erlisb/go-orchestrator/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
	Stats *Stats
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}

	taskQueued := t.(task.Task)
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}

	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("we should not transition")
		}
	} else {
		err := fmt.Errorf("invalid transition from %v to %v", taskPersisted, taskQueued)
		result.Error = err
	}

	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	result := d.Run()

	if result.Error != nil {
		log.Printf("Err running task %v: %v\n", t.ID, result.ContainerId)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	// t.ID = uuid.MustParse(result.ContainerId)
	t.State = task.Running
	w.Db[t.ID] = &t

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	result := d.Stop(t.ContainerId)

	if result.Error != nil {
		log.Printf("Error stopping container %v", result.ContainerId)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t

	log.Printf("Stopped and removed container %v for task %v", result.ContainerId, t.ID)

	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) GetTasks() map[uuid.UUID]*task.Task {
	return w.Db
}

func (w *Worker) CollectStats()  {
	for {
		log.Println("Collection stats")
		w.Stats = GetStats()
		// w.TaskCount = w.Stats.TaskCount
		time.Sleep(15 * time.Second)
	}
}