package votes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TechXTT/TORM/v2/pkg/runtime"
	"github.com/lib/pq"

	email "github.com/TechXTT/tf-api-TORM/pkg/email"
	jwt "github.com/TechXTT/tf-api-TORM/pkg/jwt"
	models "github.com/TechXTT/tf-api-TORM/pkg/models"
	tormodels "github.com/TechXTT/tf-api-TORM/torm/models"
)

func PostVoteTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	privateKey := os.Getenv("PRIVATE_KEY")
	publicKey := os.Getenv("PUBLIC_KEY")

	fmt.Println("[PostVoteTorm] Starting vote process")

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
		fmt.Println("[PostVoteTorm] Error checking cookie:", err)
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}
	if sub != 0 {
		fmt.Println("[PostVoteTorm] User has already voted")
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	var reqVote models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&reqVote); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	existing, err := srv.CountWhere(r.Context(), tormodels.VoteFields.Email.Eq(reqVote.Email))
	if err != nil {
		http.Error(w, "Error checking existing vote", http.StatusInternalServerError)
		return
	}
	if existing > 0 {
		http.Error(w, "Already voted", http.StatusBadRequest)
		return
	}

	choices := []struct {
		id       uint
		category tormodels.Category
	}{
		{reqVote.NetworksID, tormodels.CategoryNetworks},
		{reqVote.SoftwareID, tormodels.CategorySoftware},
		{reqVote.EmbeddedID, tormodels.CategoryEmbedded},
		{reqVote.BattleBotID, tormodels.CategoryBattlebot},
	}
	var filters []runtime.Filter[tormodels.Project]
	for _, choice := range choices {
		filters = append(filters, tormodels.ProjectFields.Id.Eq(int(choice.id)))
	}
	var query tormodels.ProjectQuery
	query.Filter = runtime.Or(filters...)
	projects, err := projSrv.Find(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	categories := make(map[uint]tormodels.Category, len(projects))
	for _, project := range projects {
		categories[uint(project.Id)] = project.Category
	}
	for _, choice := range choices {
		category, found := categories[choice.id]
		if !found {
			http.Error(w, "Project not found", http.StatusBadRequest)
			return
		}
		if category != choice.category {
			http.Error(w, "Invalid "+string(choice.category)+" project", http.StatusBadRequest)
			return
		}
	}

	vote, err := srv.Insert(r.Context(), tormodels.VoteCreate{
		Name:        runtime.Set(reqVote.Name),
		Email:       runtime.Set(reqVote.Email),
		NetworksId:  runtime.Set(int(reqVote.NetworksID)),
		SoftwareId:  runtime.Set(int(reqVote.SoftwareID)),
		EmbeddedId:  runtime.Set(int(reqVote.EmbeddedID)),
		BattleBotId: runtime.Set(int(reqVote.BattleBotID)),
	})
	if err != nil {
		var dbErr *pq.Error
		if errors.As(err, &dbErr) && dbErr.Code == "23505" {
			http.Error(w, "Already voted", http.StatusBadRequest)
			return
		}
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

	emailFilter := tormodels.VoteFields.Email.Eq(claims)
	unverified := runtime.And(emailFilter, tormodels.VoteFields.Verified.Eq(false))
	updated, err := srv.UpdateWhere(r.Context(), unverified, tormodels.VoteUpdate{Verified: runtime.Set(true)})
	if err != nil {
		http.Error(w, "Error saving vote", http.StatusInternalServerError)
		return
	}
	if updated == 0 {
		existing, err := srv.CountWhere(r.Context(), emailFilter)
		if err != nil {
			http.Error(w, "Error finding vote", http.StatusInternalServerError)
			return
		}
		if existing == 0 {
			http.Error(w, "Invalid vote", http.StatusBadRequest)
			return
		}
		http.Error(w, "Already verified", http.StatusBadRequest)
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
