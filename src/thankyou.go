package main

import (
	"net/http"
	"html/template"
	"fmt"
)


func main() {
	var app = new(App)
	app.init()
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", app.handll)
	http.ListenAndServe(":8080",nil)
}

type App struct {

	layout *template.Template
}

	func (ap *App)init() {
	}

	func (a *App) handll(writer http.ResponseWriter, request *http.Request) {
		if (request.Method == "POST") {

		}
		a.render(writer)
	}

	func (a *App)render(writer http.ResponseWriter) {
		layout, err := template.ParseFiles("answer.html")
		if (err != nil) {fmt.Println(err)}
		a.layout = layout
		a.layout.ExecuteTemplate(writer, "answer", nil)
	}


