package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type page struct {
	Title string
	Body  []byte
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := page{Title: "Hello world"}
	t, _ := template.ParseFiles("sensor.html")
	t.Execute(w, p)
}

func webShowSensors(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "<h1>Sensors</h1>")
}

func startWebServer() {
	// Setup simple web server
	log.Info("Starting web server on port 10081")
	http.HandleFunc("/", webShowSensors)  // set router
	go http.ListenAndServe(":10081", nil) // set listen port
}
