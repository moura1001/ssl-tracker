package main

import (
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	Title string
	Done  bool
}

type TaskPageData struct {
	PageTitle string
	Tasks     []Task
}

func tasks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/static/views/task/show.html"))

	data := TaskPageData{
		PageTitle: "Todays Task",
		Tasks: []Task{
			{Title: "title 1", Done: false},
			{Title: "title 2", Done: true},
			{Title: "title 3", Done: true},
		},
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./src/static"))
	http.Handle("/src/static/", http.StripPrefix("/src/static/", fs))

	http.HandleFunc("/tasks", tasks)

	log.Println("Server is listening on port 3000...")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
