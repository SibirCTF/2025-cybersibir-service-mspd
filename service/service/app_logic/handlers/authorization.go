package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Username: r.FormValue("username"),
		Password: utils.HashPassword(r.FormValue("password")),
	}
	//result := db.DB.Take(&user, "username=? AND password=?", user.Username, user.Password)
	/*if result.Error != nil {
		http.Redirect(w, r, "/sign_in?err=wrong_user", http.StatusSeeOther)
		return
	}
	*/
	queryResult := db.DB.Select("id, username, password").First(&user, fmt.Sprintf("username = '%v' AND password = '%v'", user.Username, user.Password))
	if queryResult.Error != nil {
		http.Redirect(w, r, "/sign_in?err=wrong_user", http.StatusSeeOther)
		return
	}
	encoded_username, e := utils.SC.Encode("User", user.Username)
	if e != nil {
		http.Redirect(w, r, "/sign_in?err=wrong_user", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "User",
		Value: encoded_username,
	})
	log.Println("Sucessful login as", user.Username)
	http.Redirect(w, r, "/sus_browser", http.StatusSeeOther)
}
