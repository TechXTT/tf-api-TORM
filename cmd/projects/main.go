package projects

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	tormodels "github.com/hacktues-9/tf-api/models"
	models "github.com/hacktues-9/tf-api/pkg/models"
	"gorm.io/gorm"
)

func GetProjects(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var projects []models.GetProjects
	db.Raw("SELECT * FROM get_projects()").Scan(&projects)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

func GetProjectsTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	srv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}

	projects, err := srv.FindMany(r.Context(), nil, nil, 0, 0)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

func GetProject(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var getProject models.GetProject
	var project models.GetProjectResponse
	db.Raw("SELECT * FROM get_project(?)", mux.Vars(r)["id"]).Scan(&getProject)

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
	project.NextId = getProject.NextId
	project.PrevId = getProject.PrevId

	var creators []models.GetProjectCreators
	db.Raw("SELECT name, concat(grade, class) as class FROM creators WHERE project_id = ?", mux.Vars(r)["id"]).Scan(&creators)
	project.Creators = creators

	var pictures []models.GetProjectPictures
	db.Raw("SELECT url, is_thumbnail FROM pictures WHERE project_id = ?", mux.Vars(r)["id"]).Scan(&pictures)
	project.Pictures = pictures

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func GetProjectTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	srv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}
	projectID := mux.Vars(r)["id"]
	getProject, err := srv.FindFirst(r.Context(), map[string]interface{}{"id": projectID})
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		return
	}

	var project models.GetProjectResponse
	project.ID = uint(getProject.Id)
	project.Name = getProject.Name
	project.Description = getProject.Description
	project.Video = getProject.VideoLink
	project.Type = string(getProject.Type)
	project.Category = string(getProject.Category)
	project.Mentor = getProject.Mentor
	project.HasThumbnail = getProject.HasThumbnail
	project.Links.Github = getProject.GithubLink
	project.Links.Demo = getProject.DemoLink

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func GetProjectsByCategory(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var projects []models.GetProjectByCategory
	db.Raw("SELECT * FROM get_projects_by_category(?)", mux.Vars(r)["category"]).Scan(&projects)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}

func GetProjectsByCategoryTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	srv, err := client.ProjectService()
	if err != nil {
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}
	category := mux.Vars(r)["category"]
	getProjects, err := srv.FindMany(r.Context(), map[string]interface{}{"category": category}, nil, 0, 0)
	if err != nil {
		http.Error(w, "Failed to fetch projects by category", http.StatusInternalServerError)
		return
	}

	var projects []models.GetProjectResponse
	for _, project := range getProjects {
		projects = append(projects, models.GetProjectResponse{
			ID:           uint(project.Id),
			Name:         project.Name,
			Description:  project.Description,
			Video:        project.VideoLink,
			Type:         string(project.Type),
			Category:     string(project.Category),
			Mentor:       project.Mentor,
			HasThumbnail: project.HasThumbnail,
			Links: models.GetProjectLinks{
				Github: project.GithubLink,
				Demo:   project.DemoLink,
			},
		})
	}

	creatorSrv, err := client.CreatorService()
	if err != nil {
		http.Error(w, "Failed to fetch creators", http.StatusInternalServerError)
		return
	}
	// get creators for each project
	for i, project := range projects {
		creators, err := creatorSrv.FindMany(r.Context(), map[string]interface{}{"project_id": project.ID}, nil, 0, 0)
		if err != nil {
			http.Error(w, "Failed to fetch creators", http.StatusInternalServerError)
			return
		}

		var creatorList []models.GetProjectCreators
		for _, creator := range creators {
			creatorList = append(creatorList, models.GetProjectCreators{
				Name:  creator.Name,
				Class: string(creator.Grade) + creator.Class,
			})
		}
		projects[i].Creators = creatorList
	}
	pictureSrv, err := client.PictureService()
	if err != nil {
		http.Error(w, "Failed to fetch pictures", http.StatusInternalServerError)
		return
	}
	// get pictures for each project
	for i, project := range projects {
		pictures, err := pictureSrv.FindMany(r.Context(), map[string]interface{}{"project_id": project.ID}, nil, 0, 0)
		if err != nil {
			http.Error(w, "Failed to fetch pictures", http.StatusInternalServerError)
			return
		}

		var pictureList []models.GetProjectPictures
		for _, picture := range pictures {
			pictureList = append(pictureList, models.GetProjectPictures{
				URL:         picture.Url,
				IsThumbnail: picture.IsThumbnail,
			})
		}
		projects[i].Pictures = pictureList
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}
