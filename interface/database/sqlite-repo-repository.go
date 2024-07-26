package database

import (
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/utils"
	"gorm.io/gorm"
)

type SqliteRepoRepository struct {
	DB *gorm.DB
}

func NewSqliteRepoRepository(db *gorm.DB) *SqliteRepoRepository {
	return &SqliteRepoRepository{DB: db}
}

func (s *SqliteRepoRepository) StoreRepositoryInfo(remoteRepoInfo *dto.RepositoryInfoResponseDTO, owner *entity.User) (*Repository, error) {
	//  logic to store repository info in the database

	// check if this remote repository already exists in our database
	existingRepo, err := s.GetRepositoryInfoByRemoteId(remoteRepoInfo.ID)
	if err != nil {
		return nil, err
	}
	if existingRepo != nil {
		// repository already exists, update existing record;
		if existingRepo.RemoteUpdatedAt != remoteRepoInfo.UpdatedAt {
			// but if only there has been an update
			return existingRepo, nil
		}
		existingRepo.Name = remoteRepoInfo.Name
		existingRepo.Description = remoteRepoInfo.Description
		existingRepo.URL = remoteRepoInfo.URL
		existingRepo.Language = remoteRepoInfo.Language
		existingRepo.ForksCount = remoteRepoInfo.ForksCount
		existingRepo.StarsCount = remoteRepoInfo.StarsCount
		existingRepo.OpenIssues = remoteRepoInfo.OpenIssues
		existingRepo.Watchers = remoteRepoInfo.Watchers
		return existingRepo, s.DB.Save(existingRepo).Error
	}
	newRepo := &Repository{
		RemoteID:        remoteRepoInfo.ID,
		OwnerID:         owner.ID,
		Name:            remoteRepoInfo.Name,
		Description:     remoteRepoInfo.Description,
		URL:             remoteRepoInfo.HtmlUrl,
		Language:        remoteRepoInfo.Language,
		ForksCount:      remoteRepoInfo.ForksCount,
		StarsCount:      remoteRepoInfo.StarsCount,
		OpenIssues:      remoteRepoInfo.OpenIssues,
		Watchers:        remoteRepoInfo.Watchers,
		RemoteCreatedAt: remoteRepoInfo.CreatedAt,
		RemoteUpdatedAt: remoteRepoInfo.UpdatedAt,
	}
	err = s.DB.Create(newRepo).Error
	if err != nil {
		return nil, err
	}
	var repository Repository
	s.DB.Preload("Owner").First(&repository, newRepo.ID)
	return &repository, nil
}

func (s *SqliteRepoRepository) GetRepositoryInfoByRemoteId(remoteID int) (*Repository, error) {
	//  logic to retrieve repository info from the database by remote ID
	repo := &Repository{}
	err := s.DB.Where("remote_id =?", remoteID).Preload("Owner").First(repo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return repo, nil
}

func (s *SqliteRepoRepository) GetRepository(ownerID uint, repoName string) (*Repository, error) {
	//  logic to retrieve repository info from the database by ID
	repo := &Repository{}
	err := s.DB.Where("owner_id =?", ownerID).Where("name =?", repoName).Preload("Owner").First(repo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return repo, nil
}

func (s *SqliteRepoRepository) SearchRepository(ownerID uint, repoSearchParams *utils.RepositorySearchParams) ([]*Repository, error) {
	//  logic to retrieve all repositories from the database
	repos := &[]*Repository{}
	dbQueryBuilder := s.DB.Preload("Owner").Where("owner_id =?", ownerID)
	if repoSearchParams.TopStarsCount > 0 {
		dbQueryBuilder = dbQueryBuilder.Order("stars_count DESC").Limit(repoSearchParams.TopStarsCount)
	}
	if repoSearchParams.Name != "" {
		dbQueryBuilder = dbQueryBuilder.Where("name LIKE?", "%"+repoSearchParams.Name+"%")
	}
	if repoSearchParams.Language != "" {
		dbQueryBuilder = dbQueryBuilder.Where("language =?", repoSearchParams.Language)
	}

	err := dbQueryBuilder.Preload("Owner").Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return *repos, nil
}

func (s *SqliteRepoRepository) GetAllRepositories() ([]*Repository, error) {
	//  logic to retrieve all repositories from the database
	repos := &[]*Repository{}
	err := s.DB.Preload("Owner").Find(&repos).Error
	if err != nil {
		return nil, err
	}
	return *repos, nil
}
