package database

import (
	"github.com/midedickson/github-service/entity"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName string `gorm:"full_name"`
	Username string `gorm:"username"`
}

func (model *User) ToEntity() *entity.User {
	return &entity.User{
		ID:       model.ID,
		FullName: model.FullName,
		Username: model.Username,
	}
}
