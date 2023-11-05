package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/erlisb/go-orchestrator/task"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
	Message    string `json:"message,omitempty"`
}

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	err := d.Decode(&te)

	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v", err)
		log.Printf("%s", msg)

		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}

		json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTask(te.Task)
	log.Printf("Added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Printf("No taskID passed in request\n")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	_, ok := a.Worker.Db[tID]
	if !ok {
		log.Printf("No task with ID %v found\n", tID)
		w.WriteHeader(404)
	}
	taskToStop := a.Worker.Db[tID]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed

	a.Worker.AddTask(taskCopy)
	log.Println("Added task %v to stop container %v\n", tID, taskToStop)
	w.WriteHeader(204)

}

func (a *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.Stats)
}
