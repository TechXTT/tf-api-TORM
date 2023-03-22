package votes

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	models "github.com/hacktues-9/tf-api/pkg/models"
)

func PostVote(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var reqVote models.VoteRequest
	var voteRes models.PostVote
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
	if networksProject.ID == 0 {
		voteRes.Msg = "Networks project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	var softwareProject models.Projects
	db.Where("id = ? AND category = 'software'", reqVote.SoftwareID).First(&softwareProject)
	if softwareProject.ID == 0 {
		voteRes.Msg = "Software project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	var embeddedProject models.Projects
	db.Where("id = ? AND category = 'embedded'", reqVote.EmbeddedID).First(&embeddedProject)
	if embeddedProject.ID == 0 {
		voteRes.Msg = "Embedded project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	var battlebotProject models.Projects
	db.Where("id = ? AND category = 'battlebot'", reqVote.BattleBotID).First(&battlebotProject)
	if battlebotProject.ID == 0 {
		voteRes.Msg = "Battlebot project not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(voteRes)
		if err != nil {
			fmt.Println("[PostVote] Error encoding JSON")
			return
		}
		return
	}

	vote := models.Votes{
		Name:        reqVote.Name,
		Email:       reqVote.Email,
		NetworksID:  reqVote.NetworksID,
		SoftwareID:  reqVote.SoftwareID,
		EmbeddedID:  reqVote.EmbeddedID,
		BattleBotID: reqVote.BattleBotID,
	}

	db.Create(&vote)

	voteRes.Msg = "Successfully voted"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(voteRes)
	if err != nil {
		fmt.Println("[PostVote] Error encoding JSON")
		return
	}
}
