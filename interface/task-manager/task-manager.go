package tasks

import (
	"github.com/midedickson/github-service/discovery"
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
)

type TaskManager struct {
	GetAllRepoForUserQueue       chan *entity.User
	FetchNewlyRequestedRepoQueue chan *dto.RepoRequest
	CheckForUpdateOnAllRepoQueue chan string
	ResetRepositoryQueue         chan *dto.RepoResetRequest
	repoDiscovery                discovery.RepositoryDiscovery
	commitManager                discovery.CommitDiscovery
}

func NewTaskManager(repoDiscovery discovery.RepositoryDiscovery, commitManager discovery.CommitDiscovery) *TaskManager {
	return &TaskManager{
		GetAllRepoForUserQueue:       make(chan *entity.User),
		FetchNewlyRequestedRepoQueue: make(chan *dto.RepoRequest),
		CheckForUpdateOnAllRepoQueue: make(chan string),
		ResetRepositoryQueue:         make(chan *dto.RepoResetRequest),
		repoDiscovery:                repoDiscovery,
		commitManager:                commitManager,
	}
}
