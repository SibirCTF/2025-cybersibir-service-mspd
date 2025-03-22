package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
)

func Sus_browser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	var suspects []models.Suspect
	result := db.DB.Order("ID desc").Limit(100).Find(&suspects)
	if result.Error != nil {
		utils.DropError(w, r, result.Error, http.StatusInternalServerError)
		return
	}
	data := struct {
		Username string
		ID       uint
		Suspects []models.Suspect
	}{
		Username: user.Username,
		ID:       user.ID,
		Suspects: suspects,
	}
	tmpl, err := template.ParseFiles("templates/sus_browser.html")
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
