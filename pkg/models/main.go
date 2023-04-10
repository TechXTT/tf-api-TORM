package models

import (
	"gorm.io/gorm"
)

type Projects struct {
	gorm.Model
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type" gorm:"check:type IN ('diploma', 'class', 'extra')"`
	Category     string `json:"category" gorm:"check:category IN ('software', 'networks', 'embedded', 'battlebot')"`
	Mentor       string `json:"mentor"`
	VideoLink    string `json:"video_link"`
	HasThumbnail bool   `json:"has_thumbnail" gorm:"default:false"`
	DemoLink     string `json:"demo_link"`
	GithubLink   string `json:"github_link"`
}

type Creators struct {
	gorm.Model
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Grade     int      `json:"grade" gorm:"check:grade IN (8, 9, 10, 11, 12)"`
	Class     string   `json:"class" gorm:"check:class IN ('А', 'Б', 'В', 'Г')"`
	ProjectID uint     `json:"project_id" gorm:"unique, not null"`
	Project   Projects `json:"project" gorm:"foreignKey:ProjectID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Pictures struct {
	gorm.Model
	URL         string   `json:"url" gorm:"unique, not null"`
	IsThumbnail bool     `json:"is_thumbnail" gorm:"default:false"`
	ProjectID   uint     `json:"project_id" gorm:"unique, not null"`
	Project     Projects `json:"project" gorm:"foreignKey:ProjectID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Votes struct {
	gorm.Model
	Name        string   `json:"name" gorm:"unique, not null"`
	Email       string   `json:"email" gorm:"unique, not null"`
	Verified    bool     `json:"verified" gorm:"default:false"`
	NetworksID  uint     `json:"networks_id" gorm:"unique, not null"`
	Networks    Projects `json:"networks" gorm:"foreignKey:NetworksID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SoftwareID  uint     `json:"software_id" gorm:"unique, not null"`
	Software    Projects `json:"software" gorm:"foreignKey:SoftwareID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EmbeddedID  uint     `json:"embedded_id" gorm:"unique, not null"`
	Embedded    Projects `json:"embedded" gorm:"foreignKey:EmbeddedID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BattleBotID uint     `json:"battlebot_id" gorm:"unique, not null"`
	BattleBot   Projects `json:"battlebot" gorm:"foreignKey:BattleBotID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
