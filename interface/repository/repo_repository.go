package repository

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/database"
	"github.com/midedickson/github-service/utils"
)

type RepoRepository interface {
	StoreRepositoryInfo(remoteRepoInfo *dto.RepositoryInfoResponseDTO, owner *entity.User) (*database.Repository, error)
	GetRepository(ownerID uint, repoName string) (*database.Repository, error)

	GetAllRepositories() ([]*database.Repository, error)
	SearchRepository(ownerID uint, repoSearchParams *utils.RepositorySearchParams) ([]*database.Repository, error)
}
