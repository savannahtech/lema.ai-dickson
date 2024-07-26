package mocks

import (
	"github.com/midedickson/github-service/interface/database"
	"github.com/stretchr/testify/mock"
)

type MockTask struct {
	mock.Mock
}

func (m *MockTask) AddUserToGetAllRepoQueue(user *database.User) {
	m.Called(user)
}

func (m *MockTask) AddRequestToFetchNewlyRequestedRepoQueue(username, repoName string) {
	m.Called(username, repoName)
}
