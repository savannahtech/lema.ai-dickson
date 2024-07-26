package discovery

import (
	"log"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/repository"
	"github.com/midedickson/github-service/requester"
)

type CommitDiscoveryService struct {
	repoRepository   repository.RepoRepository
	commitRepository repository.CommitRepository
	requester        requester.Requester
	startDateLimit   string
	endDateLimit     string
}

func NewCommitDiscoveryService(repoRepository repository.RepoRepository,
	requester requester.Requester,
	commitRepository repository.CommitRepository, startDateLimit, endDateLimit string) *CommitDiscoveryService {
	return &CommitDiscoveryService{
		repoRepository:   repoRepository,
		commitRepository: commitRepository,
		requester:        requester,
		startDateLimit:   startDateLimit,
		endDateLimit:     endDateLimit,
	}
}

func (cd *CommitDiscoveryService) GetLatestCommitSHAInRepository(repoName string) (string, error) {
	mostRecentCommit, err := cd.commitRepository.GetMostRecentCommitInRepository(repoName)
	if err != nil {
		return "", err
	}
	if mostRecentCommit != nil {
		return mostRecentCommit.SHA, nil
	}
	return "", nil
}

func (cd *CommitDiscoveryService) CheckForNewCommits(repo *entity.Repository) error {
	log.Printf("fetching new repository commits for repo: %s...", repo.Name)
	mostRecentSHA, err := cd.GetLatestCommitSHAInRepository(repo.Name)
	if err != nil {
		log.Printf("Error in fetching most recent commit SHA: %v", err)
		return err
	}
	newRemoteCommits, err := cd.requester.GetRepositoryCommits(repo.Owner.Username, repo.Name, &dto.CommitQueryParams{SHA: mostRecentSHA, Since: cd.startDateLimit, Until: cd.endDateLimit})
	if err != nil {
		log.Printf("Error in fetching new commits: %v", err)
		return err
	}
	err = cd.commitRepository.StoreRepositoryCommits(newRemoteCommits, repo.Name, repo.Owner)
	if err != nil {
		log.Printf("Error in saving new commits: %v", err)
		return err
	}
	return nil
}

func (cd *CommitDiscoveryService) GetCommitsForNewRepo(repo *entity.Repository) error {
	log.Printf("fetching repository commits for repo: %s...", repo.Name)
	remoteCommits, err := cd.requester.GetRepositoryCommits(repo.Owner.Username, repo.Name, &dto.CommitQueryParams{Since: cd.startDateLimit, Until: cd.endDateLimit})
	if err != nil {
		log.Printf("Error in fetching commits: %v", err)
		return err
	}
	cd.UpdateAuthorCountInNewCommits(*remoteCommits)
	err = cd.commitRepository.StoreRepositoryCommits(remoteCommits, repo.Name, repo.Owner)
	if err != nil {
		log.Printf("Error in saving commits: %v", err)
		return err
	}
	return nil
}

func (cd *CommitDiscoveryService) UpdateAuthorCountInNewCommits(remoteCommits []dto.CommitResponseDTO) {
	authorCommitCounts := make(map[string]int)
	for _, c := range remoteCommits {
		_, ok := authorCommitCounts[c.Author]
		if !ok {
			authorCommitCounts[c.Author] = 1
		} else {
			authorCommitCounts[c.Author]++
		}
	}
	for author := range authorCommitCounts {
		cd.commitRepository.AddAuthorCommitCount(author, authorCommitCounts[author])
	}
}

func (cd *CommitDiscoveryService) ResetCommitToSHA(repoName, resetSha string) error {
	log.Printf("resetting commits for repo: %s to SHA: %s...", repoName, resetSha)
	err := cd.commitRepository.DeleteUntilSHA(repoName, resetSha)
	if err != nil {
		log.Printf("Error in resetting commits: %v", err)
		return err
	}
	return nil
}
