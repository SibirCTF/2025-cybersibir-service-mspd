package handlers

import (
	"html/template"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
)

func Claim_manager(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	var resp string
	var claims []models.Claim
	claim_query := db.DB.Preload("Sus").Where("sus_id IN (SELECT id FROM suspects WHERE author_id = ?)", user.ID).Find(&claims)
	if claim_query.Error != nil {
		resp = "В данный момент нет заявок!"
	} else {
		resp = "OK"
	}
	if claim_query.RowsAffected == 0 {
		resp = "В данный момент нет заявок!"
	}
	data := struct {
		Username string
		Resp     string
		Claims   []models.Claim
	}{
		Username: user.Username,
		Resp:     resp,
		Claims:   claims,
	}
	tmpl, err := template.ParseFiles("templates/claim_manager.html")
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
