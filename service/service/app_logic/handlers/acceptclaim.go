package handlers

import (
	"net/http"
	"sibir2025/service/app_logic/db"
	"sibir2025/service/app_logic/models"
	"sibir2025/service/app_logic/utils"
	"strconv"
)

func Accept_claim(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthCheck(w, r)
	if err != nil {
		utils.DropError(w, r, err, http.StatusForbidden)
		return
	}
	id, e := strconv.Atoi(r.FormValue("ID"))
	if e != nil {
		http.Redirect(w, r, "/claim_manager?err=wrong_claim_id", http.StatusSeeOther)
		return
	}
	claim := models.Claim{
		ID: uint(id),
	}
	claim_query := db.DB.Take(&claim, "ID=?", claim.ID)
	if claim_query.Error != nil {
		http.Redirect(w, r, "/claim_manager?err=wrong_claim_id", http.StatusSeeOther)
		return
	}
	sus := models.Suspect{
		ID: claim.SusID,
	}
	sus_query := db.DB.Take(&sus, "ID=?", sus.ID)
	if sus_query.Error != nil {
		http.Redirect(w, r, "/claim_manager?err=db", http.StatusSeeOther)
		return
	}
	if claim.UserID == user.ID {
		http.Redirect(w, r, "/claim_manager?err=wrong_claim_id", http.StatusSeeOther)
		return
	}
	if claim.Status == "Completed" {
		http.Redirect(w, r, "/claim_manager?err=already_completed", http.StatusSeeOther)
		return
	}
	claim.Reward = sus.Sbertoken
	claim.Status = "Completed"
	db.DB.Save(&claim)
	http.Redirect(w, r, "/claim_manager?err=successful_accept", http.StatusSeeOther)
}
