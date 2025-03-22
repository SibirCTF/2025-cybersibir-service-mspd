package handlers

import (
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
	"strconv"
)

func Create_claim(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	id, e := strconv.Atoi(r.FormValue("ID"))
	if e != nil {
		http.Redirect(w, r, "/sus_browser?err=wrong_sus_id", http.StatusSeeOther)
		return
	}
	var sus models.Suspect
	if err := db.DB.First(&sus, uint(id)).Error; err != nil {
		http.Redirect(w, r, "/sus_browser?err=wrong_sus_id", http.StatusSeeOther)
		return
	}
	if sus.AuthorID == user.ID {
		http.Redirect(w, r, "/sus_browser?err=wrong_sus_id", http.StatusSeeOther)
		return
	}
	claim := models.Claim{
		UserID: user.ID,
		SusID:  sus.ID,
	}
	var count int64
	resultSearch := db.DB.Model(&models.Claim{}).Where("user_id = ? AND sus_id = ?", user.ID, sus.ID).Count(&count)
	if resultSearch.Error != nil {
		http.Redirect(w, r, "/sus_browser?err=db", http.StatusSeeOther)
		return
	}
	if count > 0 {
		http.Redirect(w, r, "/sus_browser?err=already_claimed", http.StatusSeeOther)
		return
	}
	claim.Status = "Not completed"
	claim.Reward = "Capture the suspect to get bounty!"
	claim_query := db.DB.Create(&claim)
	if claim_query.Error != nil {
		http.Redirect(w, r, "/sus_browser?err=db", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/sus_browser?err=successful_accept&id="+strconv.Itoa(int(claim.ID)), http.StatusSeeOther)
}
