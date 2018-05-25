package main

import (
	"html/template"
	"net/http"
)

type page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(cfg.templatesPath + "/hoststatistics.html"))

	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.Execute(w, data)
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func startWebServer() {
	// Setup simple web server
	log.Info("Starting web server on port 10081")
	http.HandleFunc("/", viewHandler)

	go http.ListenAndServe(":10081", nil) // set listen port
}
