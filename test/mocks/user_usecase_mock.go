package mocks

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) CreateUser(createUserPayload *dto.CreateUserPayloadDTO) (*entity.User, error) {
	args := m.Called(createUserPayload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetUser(username string) (*entity.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
