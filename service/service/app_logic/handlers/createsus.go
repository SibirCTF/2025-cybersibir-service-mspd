package handlers

import (
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
	"strconv"
)

func Create_sus(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	sus := models.Suspect{
		AuthorID:  user.ID,
		SusName:   r.FormValue("sus_name"),
		SusDesc:   r.FormValue("sus_desc"),
		CrimeDesc: r.FormValue("crime_desc"),
		Sbertoken: r.FormValue("sbertoken"),
	}
	sus_query := db.DB.Create(&sus)
	if sus_query.Error != nil {
		http.Redirect(w, r, "/sus_uploader", http.StatusSeeOther)
		return
	}
	utils.Upload_sus_img(w, r, sus)
	http.Redirect(w, r, "/sus_browser?err=successful_upload&id="+strconv.Itoa(int(sus.ID)), http.StatusSeeOther)
}
