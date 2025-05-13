package api

import "net/http"

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//POST
	case http.MethodPost:
		addTaskHandler(w, r)
	//GET
	case http.MethodGet:
		getTaskHandler(w, r)
	//PUT
	case http.MethodPut:
		updateTaskHandler(w, r)
	//DELETE
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "error 405 Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}
