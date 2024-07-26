package usecase

import (
	"errors"

	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/repository"
	tasks "github.com/midedickson/github-service/interface/task-manager"
)

type CommitUseCase interface {
	GetRepositoryCommits(repoName string) ([]*entity.Commit, error)
	MakeRepoResetRequest(owner, repoName, resetSHA string) error
	GetTopNAuthorsByCommits(topN int) ([]*entity.AuthorCommitCount, error)
}

type CommitUseCaseService struct {
	commitRepository repository.CommitRepository
	repoUseCase      RepoUseCase
	task             tasks.Task
}

func NewCommitUseCaseService(commitRepository repository.CommitRepository, repoUseCase RepoUseCase, task tasks.Task) *CommitUseCaseService {
	return &CommitUseCaseService{commitRepository: commitRepository, repoUseCase: repoUseCase, task: task}
}

func (c *CommitUseCaseService) GetRepositoryCommits(repoName string) ([]*entity.Commit, error) {
	repoCommits, err := c.commitRepository.GetRepositoryCommits(repoName)
	if err != nil {
		return nil, err
	}
	commitEntities := make([]*entity.Commit, len(repoCommits))
	for i, commit := range repoCommits {
		commitEntities[i] = commit.ToEntity()
	}
	return commitEntities, nil
}

func (c *CommitUseCaseService) MakeRepoResetRequest(owner, repoName, resetSHA string) error {
	repo, err := c.repoUseCase.GetRepositoryInfo(owner, repoName)
	if err != nil {
		return err
	}

	if repo == nil {
		return errors.New("this repository does not exist in our databse right now, but we're going to try and get it please check back in a bit")
	}

	go c.task.AddRequestToResetRepositoryQueue(repoName, resetSHA)
	return nil
}

func (c *CommitUseCaseService) GetTopNAuthorsByCommits(topN int) ([]*entity.AuthorCommitCount, error) {
	topAuthors, err := c.commitRepository.FindTopNAuthorsByCommitCounts(topN)
	if err != nil {
		return nil, err
	}
	topAuthorsEntities := make([]*entity.AuthorCommitCount, len(topAuthors))
	for i, commit := range topAuthors {
		topAuthorsEntities[i] = commit.ToEntity()
	}
	return topAuthorsEntities, nil
}
