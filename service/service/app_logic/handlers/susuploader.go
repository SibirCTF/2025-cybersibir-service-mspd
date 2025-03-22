package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/utils"
)

func Sus_uploader(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	data := struct {
		Username string
	}{
		Username: user.Username,
	}
	tmpl, err := template.ParseFiles("templates/sus_uploader.html")
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		utils.DropError(w, r, err, http.StatusInternalServerError)
		return
	}
}
