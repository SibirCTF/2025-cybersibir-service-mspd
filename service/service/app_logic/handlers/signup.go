package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/utils"
)

func Sign_up(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
}
