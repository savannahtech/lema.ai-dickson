package mocks

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/interface/database"
	"github.com/midedickson/github-service/utils"
	"github.com/stretchr/testify/mock"
)

type MockDBRepository struct {
	mock.Mock
}

func (m *MockDBRepository) CreateUser(createUserPayload *dto.CreateUserPayloadDTO) (*database.User, error) {
	args := m.Called(createUserPayload)
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockDBRepository) GetUser(username string) (*database.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockDBRepository) StoreRepositoryInfo(remoteRepoInfo *dto.RepositoryInfoResponseDTO, owner *database.User) (*database.Repository, error) {
	args := m.Called(remoteRepoInfo, owner)
	return args.Get(0).(*database.Repository), args.Error(1)
}

func (m *MockDBRepository) GetRepository(ownerID uint, repoName string) (*database.Repository, error) {
	args := m.Called(ownerID, repoName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.Repository), args.Error(1)
}

func (m *MockDBRepository) StoreRepositoryCommits(commitRepoInfos *[]dto.CommitResponseDTO, repoName string, owner *database.User) error {
	args := m.Called(commitRepoInfos, repoName, owner)
	return args.Error(0)
}

func (m *MockDBRepository) GetRepositoryCommits(repoName string) ([]*database.Commit, error) {
	args := m.Called(repoName)
	return args.Get(0).([]*database.Commit), args.Error(1)
}

func (m *MockDBRepository) GetAllRepositories() ([]*database.Repository, error) {
	args := m.Called()
	return args.Get(0).([]*database.Repository), args.Error(1)
}

func (m *MockDBRepository) SearchRepository(ownerID uint, repoSearchParams *utils.RepositorySearchParams) ([]*database.Repository, error) {
	args := m.Called(ownerID, repoSearchParams)
	return args.Get(0).([]*database.Repository), args.Error(1)
}
