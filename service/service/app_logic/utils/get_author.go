package utils

import (
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
)

func GetAuthor(sus models.Suspect) models.User {
	author := models.User{
		ID: sus.AuthorID,
	}
	search_query := db.DB.Take(&author, "ID=?", author.ID)
	if search_query.Error != nil {
		author.Username = "Unknown"
		author.ID = 0
	}
	return author
}
