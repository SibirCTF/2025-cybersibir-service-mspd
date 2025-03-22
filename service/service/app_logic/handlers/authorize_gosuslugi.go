package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
)

func AuthGosuslugi(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Username: r.FormValue("username"),
	}
	guAuthKey := r.FormValue("gosuslugi_auth_key")
	if guAuthKey == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No authorization key!")
		return
	}
	err := utils.CheckAuthKey(guAuthKey, user.Username, r)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Failed to check authorization key!")
		return
	}
	result := db.DB.Take(&user, "username=?", user.Username)
	if result.Error != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Failed to check username!")
		return
	}
	encoded_username, e := utils.SC.Encode("User", r.FormValue("username"))
	if e != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to handle request!")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "User",
		Value: encoded_username,
	})
	log.Println("Sucessful login as", user.Username)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Sucessfully authorized as %s", user.Username)
}
