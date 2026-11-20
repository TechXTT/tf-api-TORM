package projects

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/TechXTT/tf-api-TORM/pkg/models"
	tormodels "github.com/TechXTT/tf-api-TORM/torm/models"
	"github.com/gorilla/mux"
)

var include = tormodels.ProjectInclude{Creators: true, Pictures: true}

func text(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func projectResponse(project *tormodels.Project) models.GetProjectResponse {
	var response = models.GetProjectResponse{
		ID:           uint(project.Id),
		Name:         project.Name,
		Description:  text(project.Description),
		Video:        text(project.VideoLink),
		Type:         string(project.Type),
		Category:     string(project.Category),
		Mentor:       text(project.Mentor),
		HasThumbnail: project.HasThumbnail,
		Links:        models.GetProjectLinks{Github: text(project.GithubLink), Demo: text(project.DemoLink)},
	}

	for _, creator := range project.Creators {
		response.Creators = append(response.Creators, models.GetProjectCreators{
			Name: creator.Name, Class: fmt.Sprint(creator.Grade) + creator.Class,
		})
	}
	for _, picture := range project.Pictures {
		response.Pictures = append(response.Pictures, models.GetProjectPictures{
			URL: picture.Url, IsThumbnail: picture.IsThumbnail,
		})
	}
	return response
}

func GetProjectsTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	var query = tormodels.ProjectQuery{Include: include}

	if category, ok := mux.Vars(r)["category"]; ok {
		switch value := tormodels.Category(category); value {
		case tormodels.CategoryNetworks, tormodels.CategorySoftware, tormodels.CategoryEmbedded, tormodels.CategoryBattlebot:
			query.Filter = tormodels.ProjectFields.Category.Eq(value)
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]models.GetProjectResponse(nil))
			return
		}
	}
	srv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}
	projects, err := srv.Find(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	var responses []models.GetProjectResponse
	for _, project := range projects {
		responses = append(responses, projectResponse(project))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

func GetProjectTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}
	srv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}
	project, err := srv.Fetch(r.Context(), tormodels.ProjectById(id), include)
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectResponse(project))
}
