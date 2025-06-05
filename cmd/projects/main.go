package projects

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/TechXTT/tf-api-TORM/pkg/models"
	tormodels "github.com/TechXTT/tf-api-TORM/torm/models"
	"github.com/gorilla/mux"
)

func GetProjectsTorm(w http.ResponseWriter, r *http.Request, client *tormodels.Client) {
	srv, err := client.ProjectService()
	if err != nil && srv == nil {
		log.Printf("Error initializing project service: %v", err)
		http.Error(w, "Failed to initialize project service", http.StatusInternalServerError)
		return
	}

	projects, err := srv.FindMany(r.Context(), nil, nil, 0, 0)
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	fmt.Println("Projects fetched successfully")
	var projectResponses []models.GetProjectResponse
	for _, project := range projects {

		fmt.Println("Processing project:", project.Name, "Creators Name:", project.Creators[0].Name, "Pictures count:", len(project.Pictures))
		projectResponses = append(projectResponses, models.GetProjectResponse{
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
			Creators: func() []models.GetProjectCreators {
				var creators []models.GetProjectCreators
				for _, creator := range project.Creators {
					creators = append(creators, models.GetProjectCreators{
						Name:  creator.Name,
						Class: string(creator.Grade) + creator.Class,
					})
				}
				return creators
			}(),
			Pictures: func() []models.GetProjectPictures {
				var pictures []models.GetProjectPictures
				for _, picture := range project.Pictures {
					pictures = append(pictures, models.GetProjectPictures{
						URL:         picture.Url,
						IsThumbnail: picture.IsThumbnail,
					})
				}
				return pictures
			}(),
		})

	}

	w.Header().Set("Content-Type", "application/json charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projectResponses)
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
	project.Creators = func() []models.GetProjectCreators {
		var creators []models.GetProjectCreators
		for _, creator := range getProject.Creators {
			creators = append(creators, models.GetProjectCreators{
				Name:  creator.Name,
				Class: string(creator.Grade) + creator.Class,
			})
		}
		return creators
	}()
	// Fetch pictures
	pictureSrv, err := client.PictureService()
	if err != nil {
		http.Error(w, "Failed to fetch pictures", http.StatusInternalServerError)
		return
	}
	pictures, err := pictureSrv.FindMany(r.Context(), map[string]interface{}{"projectid": uint(getProject.Id)}, nil, 0, 0)
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
	project.Pictures = pictureList

	w.Header().Set("Content-Type", "application/json charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
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
		creators := func() []models.GetProjectCreators {
			var creators []models.GetProjectCreators
			for _, creator := range project.Creators {
				creators = append(creators, models.GetProjectCreators{
					Name:  creator.Name,
					Class: string(creator.Grade) + creator.Class,
				})
			}
			return creators
		}()

		picturesList := func() []models.GetProjectPictures {
			var pictures []models.GetProjectPictures
			for _, picture := range project.Pictures {
				pictures = append(pictures, models.GetProjectPictures{
					URL:         picture.Url,
					IsThumbnail: picture.IsThumbnail,
				})
			}
			return pictures
		}()

		// Append project details to the response
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
			Creators: creators,
			Pictures: picturesList,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(projects)
}
