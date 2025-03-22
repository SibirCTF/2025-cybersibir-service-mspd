package handlers

import (
	"log"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Username: r.FormValue("username"),
		Password: utils.HashPassword(r.FormValue("password")),
	}
	var count int64
	resultSearch := db.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if resultSearch.Error != nil {
		http.Redirect(w, r, "/sign_up?err=wrong_user", http.StatusSeeOther)
		return
	}
	if count > 0 {
		http.Redirect(w, r, "/sign_up?err=wrong_user", http.StatusSeeOther)
		return
	}
	result := db.DB.Create(&user)
	if result.Error != nil {
		http.Redirect(w, r, "/sign_up?err=registration_failed", http.StatusSeeOther)
		return
	}
	encoded_username, e := utils.SC.Encode("User", user.Username)
	if e != nil {
		http.Redirect(w, r, "/sign_up?err=wrong_user", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "User",
		Value: encoded_username,
	})
	log.Println("Sucessful register as", user.Username)
	http.Redirect(w, r, "/sus_browser", http.StatusSeeOther)
}
