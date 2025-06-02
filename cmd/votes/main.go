package votes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	tormodels "github.com/hacktues-9/tf-api/models"
	email "github.com/hacktues-9/tf-api/pkg/email"
	jwt "github.com/hacktues-9/tf-api/pkg/jwt"
	models "github.com/hacktues-9/tf-api/pkg/models"
)

func PostVote(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	privateKey := os.Getenv("PRIVATE_KEY")
	publicKey := os.Getenv("PUBLIC_KEY")

	var reqVote models.VoteRequest
	var voteRes models.PostVote

	sub, err := jwt.CheckCookie(r)
	if err != nil {
		voteRes.Msg = "Already voted"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	if sub != 0 {
		voteRes.Msg = "Already voted"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	//create a query for gorm
	fieldsToOmit := []string{}
	if err := json.NewDecoder(r.Body).Decode(&reqVote); err != nil {
		voteRes.Msg = "Error decoding JSON"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	var dupVote models.Votes
	db.Where("email = ? ", reqVote.Email).First(&dupVote)
	if dupVote.ID != 0 {
		voteRes.Msg = "Already voted"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	var networksProject models.Projects
	db.Where("id = ? AND category = 'networks'", reqVote.NetworksID).First(&networksProject)
	if networksProject.ID == 0 && reqVote.NetworksID != 0 {
		voteRes.Msg = "Networks project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	} else {
		fieldsToOmit = append(fieldsToOmit, "NetworksID")
	}

	var softwareProject models.Projects
	db.Where("id = ? AND category = 'software'", reqVote.SoftwareID).First(&softwareProject)
	if softwareProject.ID == 0 && reqVote.SoftwareID != 0 {
		voteRes.Msg = "Software project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	} else {
		fieldsToOmit = append(fieldsToOmit, "SoftwareID")
	}

	var embeddedProject models.Projects
	db.Where("id = ? AND category = 'embedded'", reqVote.EmbeddedID).First(&embeddedProject)
	if embeddedProject.ID == 0 && reqVote.EmbeddedID != 0 {
		voteRes.Msg = "Embedded project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	} else {
		fieldsToOmit = append(fieldsToOmit, "EmbeddedID")
	}

	var battlebotProject models.Projects
	db.Where("id = ? AND category = 'battlebot'", reqVote.BattleBotID).First(&battlebotProject)
	if battlebotProject.ID == 0 && reqVote.BattleBotID != 0 {
		voteRes.Msg = "Battlebot project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	} else {
		fieldsToOmit = append(fieldsToOmit, "BattleBotID")
	}

	vote := models.Votes{
		Name:        reqVote.Name,
		Email:       reqVote.Email,
		NetworksID:  reqVote.NetworksID,
		SoftwareID:  reqVote.SoftwareID,
		EmbeddedID:  reqVote.EmbeddedID,
		BattleBotID: reqVote.BattleBotID,
	}

	if err := db.Create(&vote).Error; err != nil {
		//print the query

		voteRes.Msg = "Error creating vote"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	data := struct {
		RecieverName     string
		SenderName       string
		VerificationLink string
	}{
		RecieverName:     reqVote.Name,
		SenderName:       "TuesFest 2023",
		VerificationLink: email.GenerateVerificationLink(reqVote.Email, privateKey, publicKey, 24*time.Hour),
	}

	if data.VerificationLink == "" {
		voteRes.Msg = "Error generating verification link"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	email.OAuthGmailService()
	_, err = email.SendEmailOAUTH2(reqVote.Email, data, "template.txt")
	if err != nil {
		voteRes.Msg = "Error sending email"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	token, err := jwt.CreateToken(24*time.Hour, vote.ID, os.Getenv("PRIVATE_KEY"), os.Getenv("PUBLIC_KEY"))
	if err != nil {
		voteRes.Msg = "Error creating token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	tokenCookie := http.Cookie{
		Name:     "vote",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Domain:   os.Getenv("DOMAIN"),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &tokenCookie)
	voteRes.Msg = "Successfully voted"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(voteRes)
	if err != nil {
		fmt.Println("[PostVote] Error encoding JSON")
		return
	}
}

func PostVoteTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	privateKey := os.Getenv("PRIVATE_KEY")
	publicKey := os.Getenv("PUBLIC_KEY")

	srv, err := client.VoteService()
	if err != nil {
		http.Error(w, "Failed to initialize vote service", http.StatusInternalServerError)
		return
	}

	var voteRes models.PostVote

	projSrv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}

	// Check if the user has already voted
	sub, err := jwt.CheckCookie(r)
	if err != nil {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}
	if sub != 0 {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	// Check if the user has voted
	if _, err := srv.FindUnique(r.Context(), map[string]interface{}{"email": sub}); err != nil {
		if err.Error() == "record not found" {
			// User has not voted yet, proceed with creating a new vote
			fmt.Println("[PostVoteTorm] User has not voted yet")
		} else {
			http.Error(w, "Error checking existing vote", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	var reqVote models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqVote); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// check if each project exists
	if netProj, err := projSrv.FindUnique(r.Context(), map[string]interface{}{"id": reqVote.NetworksID}); err != nil {
		http.Error(w, "Project not found", http.StatusBadRequest)
		return
	} else if netProj.Category != "networks" {
		http.Error(w, "Invalid networks project", http.StatusBadRequest)
		return
	}
	if softProj, err := projSrv.FindUnique(r.Context(), map[string]interface{}{"id": reqVote.SoftwareID}); err != nil {
		http.Error(w, "Project not found", http.StatusBadRequest)
		return
	} else if softProj.Category != "software" {
		http.Error(w, "Invalid software project", http.StatusBadRequest)
		return
	}
	if embProj, err := projSrv.FindUnique(r.Context(), map[string]interface{}{"id": reqVote.EmbeddedID}); err != nil {
		http.Error(w, "Project not found", http.StatusBadRequest)
		return
	} else if embProj.Category != "embedded" {
		http.Error(w, "Invalid embedded project", http.StatusBadRequest)
		return
	}
	if battProj, err := projSrv.FindUnique(r.Context(), map[string]interface{}{"id": reqVote.BattleBotID}); err != nil {
		http.Error(w, "Project not found", http.StatusBadRequest)
		return
	} else if battProj.Category != "battlebot" {
		http.Error(w, "Invalid battlebot project", http.StatusBadRequest)
		return
	}
	// Create the vote

	vote := tormodels.Vote{
		Name:        reqVote.Name,
		Email:       reqVote.Email,
		NetworksId:  int(reqVote.NetworksID),
		SoftwareId:  int(reqVote.SoftwareID),
		EmbeddedId:  int(reqVote.EmbeddedID),
		BattleBotId: int(reqVote.BattleBotID),
	}

	if err := srv.Create(r.Context(), &vote); err != nil {
		http.Error(w, "Error creating vote", http.StatusInternalServerError)
		return
	}
	data := struct {
		RecieverName     string
		SenderName       string
		VerificationLink string
	}{
		RecieverName:     reqVote.Name,
		SenderName:       "TuesFest 2023",
		VerificationLink: email.GenerateVerificationLink(reqVote.Email, privateKey, publicKey, 24*time.Hour),
	}

	if data.VerificationLink == "" {
		voteRes.Msg = "Error generating verification link"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	email.OAuthGmailService()
	_, err = email.SendEmailOAUTH2(reqVote.Email, data, "template.txt")
	if err != nil {
		voteRes.Msg = "Error sending email"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	token, err := jwt.CreateToken(24*time.Hour, vote.Id, privateKey, publicKey)
	if err != nil {
		voteRes.Msg = "Error creating token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	tokenCookie := http.Cookie{
		Name:     "vote",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Domain:   os.Getenv("DOMAIN"),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}

	http.SetCookie(w, &tokenCookie)
	voteRes.Msg = "Successfully voted"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(voteRes)
	if err != nil {
		fmt.Println("[PostVote] Error encoding JSON")
		return
	}
}

func VerifyVote(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var reqVote models.VerifyVoteRequest
	var voteRes models.PostVote

	if err := json.NewDecoder(r.Body).Decode(&reqVote); err != nil {
		voteRes.Msg = "Error"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Error decoding JSON")
		return
	}

	if reqVote.Token == "" {
		voteRes.Msg = "Invalid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Token not found")
		return
	}

	claims, err := email.ValidateEmailToken(reqVote.Token)
	if err != nil {
		voteRes.Msg = "Invalid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Error validating token")
		return
	}

	var vote models.Votes
	db.Where("email = ?", claims).First(&vote)
	if vote.ID == 0 {
		voteRes.Msg = "Invalid"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Error finding vote")
		return
	}

	if vote.Verified {
		voteRes.Msg = "Already"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Vote already verified")
		return
	}

	vote.Verified = true

	if err := db.Save(&vote).Error; err != nil {
		voteRes.Msg = "Error"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[VerifyVote] Error encoding JSON")
			return
		}
		fmt.Println("[VerifyVote] Error saving vote")
		return
	}

	voteRes.Msg = "Successfully verified vote"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(voteRes)
	if err != nil {
		fmt.Println("[VerifyVote] Error encoding JSON")
		return
	}
}

func VerifyVoteTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	srv, err := client.VoteService()
	if err != nil {
		http.Error(w, "Failed to initialize vote service", http.StatusInternalServerError)
		return
	}

	var reqVote models.VerifyVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqVote); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	if reqVote.Token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	claims, err := email.ValidateEmailToken(reqVote.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	vote, err := srv.FindUnique(r.Context(), map[string]interface{}{"email": claims})
	if err != nil {
		if err.Error() == "record not found" {
			http.Error(w, "Invalid vote", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error finding vote", http.StatusInternalServerError)
		return
	}

	if vote.Verified {
		http.Error(w, "Already verified", http.StatusBadRequest)
		return
	}

	vote.Verified = true

	if err := srv.Update(r.Context(), map[string]interface{}{"id": vote.Id}, vote); err != nil {
		http.Error(w, "Error saving vote", http.StatusInternalServerError)
		return
	}

	voteRes := models.PostVote{Msg: "Successfully verified vote"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(voteRes)
	if err != nil {
		fmt.Println("[VerifyVoteTorm] Error encoding JSON")
		return
	}
}
