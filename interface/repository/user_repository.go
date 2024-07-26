package repository

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/interface/database"
)

type UserRepository interface {
	CreateUser(createUserPaylod *dto.CreateUserPayloadDTO) (*database.User, error)
	GetUser(username string) (*database.User, error)
}
