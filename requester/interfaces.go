package requester

import (
	"github.com/midedickson/github-service/dto"
)

type Requester interface {
	GetRepositoryInfo(owner, repo string) (*dto.RepositoryInfoResponseDTO, error)
	GetRepositoryCommits(owner, repo string, queryParams *dto.CommitQueryParams) (*[]dto.CommitResponseDTO, error)
	GetAllUserRepositories(owner string) (*[]dto.RepositoryInfoResponseDTO, error)
}
