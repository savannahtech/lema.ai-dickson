package repository

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/database"
)

type CommitRepository interface {
	StoreRepositoryCommits(commitRepoInfos *[]dto.CommitResponseDTO, repoName string, owner *entity.User) error
	GetRepositoryCommits(repoName string) ([]*database.Commit, error)
	GetMostRecentCommitInRepository(repoName string) (*database.Commit, error)
	DeleteUntilSHA(repoName, sha string) error
	FindTopNAuthorsByCommitCounts(topN int) ([]*database.AuthorCommitCount, error)
	AddAuthorCommitCount(author string, count int) error
}
