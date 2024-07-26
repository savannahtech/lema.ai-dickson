package usecase

import (
	"log"

	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/repository"
	tasks "github.com/midedickson/github-service/interface/task-manager"
	"github.com/midedickson/github-service/utils"
)

type RepoUseCase interface {
	GetRepositoryInfo(owner, repoName string) (*entity.Repository, error)
	GetUserRepositories(username string, repoSearchParams *utils.RepositorySearchParams) ([]*entity.Repository, error)
}

type RepoUseCaseService struct {
	repoRepository repository.RepoRepository
	userUseCase    UserUseCase
	task           tasks.Task
}

func NewRepoUseCaseService(repoRepository repository.RepoRepository, userUseCase UserUseCase, task tasks.Task) *RepoUseCaseService {
	return &RepoUseCaseService{repoRepository: repoRepository, userUseCase: userUseCase, task: task}
}

func (r *RepoUseCaseService) GetRepositoryInfo(username, repoName string) (*entity.Repository, error) {
	user, err := r.userUseCase.GetUser(username)
	if err != nil {
		return nil, err
	}
	repo, err := r.repoRepository.GetRepository(user.ID, repoName)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		go r.task.AddRequestToFetchNewlyRequestedRepoQueue(user.Username, repoName)
		return nil, nil
	}
	return repo.ToEntity(), nil
}

func (r *RepoUseCaseService) GetUserRepositories(username string, repoSearchParams *utils.RepositorySearchParams) ([]*entity.Repository, error) {
	// Logic to fetch repositories for a specific user
	user, err := r.userUseCase.GetUser(username)
	if err != nil {
		return nil, err
	}
	repositories, err := r.repoRepository.SearchRepository(user.ID, repoSearchParams)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	repositoryEntities := make([]*entity.Repository, len(repositories))
	for i, repo := range repositories {
		repositoryEntities[i] = repo.ToEntity()
	}
	return repositoryEntities, nil
}
