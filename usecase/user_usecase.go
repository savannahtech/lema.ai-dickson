package usecase

import (
	"errors"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/repository"
	tasks "github.com/midedickson/github-service/interface/task-manager"
)

type UserUseCase interface {
	CreateUser(createUserPayload *dto.CreateUserPayloadDTO) (*entity.User, error)
	GetUser(username string) (*entity.User, error)
}

type UserUseCaseService struct {
	userRepository repository.UserRepository
	task           tasks.Task
}

func NewUserUseCaseService(userRepository repository.UserRepository, task tasks.Task) *UserUseCaseService {
	return &UserUseCaseService{userRepository: userRepository, task: task}
}

func (u *UserUseCaseService) CreateUser(createUserPayload *dto.CreateUserPayloadDTO) (*entity.User, error) {
	dbUser, err := u.userRepository.CreateUser(createUserPayload)
	if err != nil {
		return nil, err
	}
	user := dbUser.ToEntity()
	go u.task.AddUserToGetAllRepoQueue(user)
	return user, nil
}

func (u *UserUseCaseService) GetUser(username string) (*entity.User, error) {
	dbUser, err := u.userRepository.GetUser(username)
	if err != nil {
		return nil, err
	}
	if dbUser == nil {
		return nil, errors.New("user with the username " + username + "not found")
	}
	return dbUser.ToEntity(), nil
}
