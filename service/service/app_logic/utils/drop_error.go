package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

func DropError(w http.ResponseWriter, r *http.Request, e error, status int) {
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("templates/not_found.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := struct {
		Status int
	}{
		Status: status,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(e)
}
