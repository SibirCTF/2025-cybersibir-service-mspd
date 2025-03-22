package handlers

import (
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
	"strconv"
)

func Download_profile(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	susID, e := strconv.Atoi(r.FormValue("ID"))
	ext := r.FormValue("ext")
	if e != nil {
		susID = 0
	}
	sus := models.Suspect{
		ID: uint(susID),
	}
	if err := db.DB.First(&sus, uint(susID)).Error; err != nil {
		sus.AuthorID = 0
		sus.CrimeDesc = "Unknown"
		sus.SusDesc = "Unknown"
		sus.SusName = "Unknown"
	}
	utils.CreateSusProfile(w, r, susID, r.FormValue("ID"), ext, sus, user)
}
