package utils

import (
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
)

func AuthCheck(w http.ResponseWriter, r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("User")
	if err != nil {
		return models.User{}, err
	}
	var decoded_username string
	if err := SC.Decode("User", cookie.Value, &decoded_username); err != nil {
		return models.User{}, err
	}
	user := models.User{
		Username: decoded_username,
	}
	search_query := db.DB.Take(&user, "username=?", user.Username)
	if search_query.Error != nil {
		return models.User{}, err
	}
	return user, nil
}
