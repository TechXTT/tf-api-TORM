package models

type GetProjects struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
	Video     string `json:"video"`
	Category  string `json:"category"`
}

type GetProject struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Video        string `json:"video"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Mentor       string `json:"mentor"`
	HasThumbnail bool   `json:"has_thumbnail"`
	Demo         string `json:"demo"`
	Github       string `json:"github"`
}

type GetProjectLinks struct {
	Github string `json:"github"`
	Demo   string `json:"demo"`
}

type GetProjectCreators struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

type GetProjectPictures struct {
	URL         string `json:"url"`
	IsThumbnail bool   `json:"is_thumbnail"`
}

type GetProjectResponse struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Video        string               `json:"video"`
	Type         string               `json:"type"`
	Category     string               `json:"category"`
	Mentor       string               `json:"mentor"`
	HasThumbnail bool                 `json:"has_thumbnail"`
	Links        GetProjectLinks      `json:"links"`
	Creators     []GetProjectCreators `json:"creators"`
	Pictures     []GetProjectPictures `json:"pictures"`
}

type GetProjectByCategory struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}
type PostVote struct {
	Msg string `json:"msg"`
}
