package database

import (
	"github.com/midedickson/github-service/entity"
	"gorm.io/gorm"
)

type AuthorCommitCount struct {
	gorm.Model
	Author      string `gorm:"author"`
	CommitCount int    `gorm:"commit_count"`
}

func (model *AuthorCommitCount) ToEntity() *entity.AuthorCommitCount {
	return &entity.AuthorCommitCount{
		Author:      model.Author,
		CommitCount: model.CommitCount,
	}
}
