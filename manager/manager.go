package manager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/erlisb/go-orchestrator/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]task.Task
	EventDb       map[string][]task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	LastWorker int
}

func (m *Manager) SelectWorker() string {
	var newWorker int
	if m.LastWorker + 1 < len(m.Workers){
		newWorker = m.LastWorker + 1
	}else {
		newWorker = 0
		m.LastWorker = 0
	}

	return m.Workers[newWorker]
}

func (m *Manager) UpdateTasks() {
	fmt.Println("I will update tasks")
}

func (m *Manager) SendWork() {
	if m.Pending.Len() > 0 {
		w := m.SelectWorker()
		e := m.Pending.Dequeue()

		te := e.(task.TaskEvent)
		t := te.Task
		log.Printf("Pulled %v off pending queue", t)

		m.EventDb[te.ID] = &te
		m.WorkerTaskMap[w] = append(m.WorkerTaskMap[w], t)
		m.TaskWorkerMap[t.ID] = w

		t.State = task.Scheduled
		m.TaskDb[string(t.ID)] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Unable to marshal task object")
		}

	}
}


