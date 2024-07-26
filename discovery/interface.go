package discovery

import (
	"sync"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
)

type RepositoryDiscovery interface {
	GetAllUserRepositories(user *entity.User)
	FetchNewlyRequestedRepo(repoRequest *dto.RepoRequest, wg *sync.WaitGroup)
	CheckForUpdateOnAllRepo() error
}

type CommitDiscovery interface {
	CheckForNewCommits(repo *entity.Repository) error
	GetCommitsForNewRepo(repo *entity.Repository) error
	ResetCommitToSHA(repoName, resetSha string) error
}
