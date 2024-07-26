package mocks

import (
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/utils"
	"github.com/stretchr/testify/mock"
)

type MockRepoUseCase struct {
	mock.Mock
}

func (m *MockRepoUseCase) GetRepositoryInfo(owner, repoName string) (*entity.Repository, error) {
	args := m.Called(owner, repoName)
	var repo *entity.Repository
	if args.Get(0) != nil {
		repo = args.Get(0).(*entity.Repository)
	}
	return repo, args.Error(1)
}

func (m *MockRepoUseCase) GetUserRepositories(username string, repoSearchParams *utils.RepositorySearchParams) ([]*entity.Repository, error) {
	args := m.Called(username, repoSearchParams)
	var repositories []*entity.Repository
	if args.Get(0) != nil {
		repositories = args.Get(0).([]*entity.Repository)
	}
	return repositories, args.Error(1)
}
