package projects

import (
	"encoding/json"
	models "github.com/hacktues-9/tf-api/pkg/models"
	"gorm.io/gorm"
	"net/http"
)

func GetProjects(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var projects []models.GetProjects
	db.Raw("SELECT * FROM get_projects()").Scan(&projects)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

func GetProject(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var getProject models.GetProject
	var project models.GetProjectResponse
	db.Raw("SELECT * FROM get_project(?)", r.URL.Query().Get("id")).Scan(&getProject)

	project.ID = getProject.ID
	project.Name = getProject.Name
	project.Description = getProject.Description
	project.Video = getProject.Video
	project.Type = getProject.Type
	project.Category = getProject.Category
	project.Mentor = getProject.Mentor
	project.HasThumbnail = getProject.HasThumbnail
	project.Links.Github = getProject.Github
	project.Links.Demo = getProject.Demo

	var creators []models.GetProjectCreators
	db.Raw("SELECT name, concat(grade, class) as class FROM creators WHERE project_id = ?", r.URL.Query().Get("id")).Scan(&creators)
	project.Creators = creators

	var pictures []models.GetProjectPictures
	db.Raw("SELECT url, is_thumbnail FROM pictures WHERE project_id = ?", r.URL.Query().Get("id")).Scan(&pictures)
	project.Pictures = pictures

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func GetProjectsByCategory(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var projects []models.GetProjectByCategory
	db.Raw("SELECT * FROM get_projects_by_category(?)", r.URL.Query().Get("category")).Scan(&projects)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}
