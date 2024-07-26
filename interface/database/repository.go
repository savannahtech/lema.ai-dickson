package database

import (
	"github.com/midedickson/github-service/entity"
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	RemoteID        int    `gorm:"remote_id"`
	OwnerID         uint   `gorm:"owner_id"`
	Owner           *User  `gorm:"foreignKey:OwnerID"`
	Name            string `gorm:"name"`
	Description     string `gorm:"description"`
	URL             string `gorm:"html_url"`
	Language        string `gorm:"language"`
	ForksCount      int    `gorm:"forks_count"`
	StarsCount      int    `gorm:"stargazers_count"`
	OpenIssues      int    `gorm:"open_issues_count"`
	Watchers        int    `gorm:"watchers_count"`
	RemoteCreatedAt string `gorm:"remote_created_at"`
	RemoteUpdatedAt string `gorm:"remote_updated_at"`
}

func (model *Repository) ToEntity() *entity.Repository {
	return &entity.Repository{
		ID:              model.ID,
		RemoteID:        model.RemoteID,
		Owner:           model.Owner.ToEntity(),
		Name:            model.Name,
		Description:     model.Description,
		URL:             model.URL,
		Language:        model.Language,
		ForksCount:      model.ForksCount,
		StarsCount:      model.StarsCount,
		OpenIssues:      model.OpenIssues,
		Watchers:        model.Watchers,
		RemoteCreatedAt: model.RemoteCreatedAt,
		RemoteUpdatedAt: model.RemoteUpdatedAt,
	}
}
