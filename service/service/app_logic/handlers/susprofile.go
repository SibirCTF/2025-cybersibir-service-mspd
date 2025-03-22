package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
	"strconv"
	"strings"
)

func Sus_profile(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	urlnum := strings.TrimPrefix(r.URL.Path, "/sus")
	urlnum = strings.Trim(urlnum, "/")
	sus_ID, err := strconv.Atoi(urlnum)
	if err != nil {
		http.Redirect(w, r, "/sus_browser?err=wrong_sus_id", http.StatusSeeOther)
		return
	}
	sus := models.Suspect{
		ID: uint(sus_ID),
	}
	result := db.DB.Take(&sus, "ID=?", sus.ID)
	if result.Error != nil {
		http.Redirect(w, r, "/sus_browser?err=wrong_sus_id", http.StatusSeeOther)
		return
	}
	author := utils.GetAuthor(sus)
	data := struct {
		Username string
		ID       uint
		Sus      models.Suspect
		Author   models.User
	}{
		Username: user.Username,
		ID:       user.ID,
		Sus:      sus,
		Author:   author,
	}
	tmpl, err := template.ParseFiles("templates/sus_profile.html")
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
