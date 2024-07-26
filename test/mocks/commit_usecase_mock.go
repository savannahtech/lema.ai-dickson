package mocks

import (
	"github.com/midedickson/github-service/entity"
	"github.com/stretchr/testify/mock"
)

type MockCommitUseCase struct {
	mock.Mock
}

func (m *MockCommitUseCase) GetRepositoryCommits(repoName string) ([]*entity.Commit, error) {
	args := m.Called(repoName)
	var commits []*entity.Commit
	if args.Get(0) != nil {
		commits = args.Get(0).([]*entity.Commit)
	}
	return commits, args.Error(1)
}

func (m *MockCommitUseCase) MakeRepoResetRequest(owner, repoName, resetSHA string) error {
	args := m.Called(owner, repoName, resetSHA)
	return args.Error(0)
}

func (m *MockCommitUseCase) GetTopNAuthorsByCommits(topN int) ([]*entity.AuthorCommitCount, error) {
	args := m.Called(topN)
	var authorCounts []*entity.AuthorCommitCount
	if args.Get(0) != nil {
		authorCounts = args.Get(0).([]*entity.AuthorCommitCount)
	}
	return authorCounts, args.Error(1)
}
