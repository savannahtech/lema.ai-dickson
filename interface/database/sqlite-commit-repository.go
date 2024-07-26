package database

import (
	"fmt"
	"log"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"gorm.io/gorm"
)

type SqliteCommitRepository struct {
	DB *gorm.DB
}

func NewSqliteCommitRepository(db *gorm.DB) *SqliteCommitRepository {
	return &SqliteCommitRepository{DB: db}
}

func (s *SqliteCommitRepository) StoreRepositoryCommits(commitRepoInfos *[]dto.CommitResponseDTO, repoName string, owner *entity.User) error {
	//  logic to store commit info in the database
	repo := &Repository{}
	err := s.DB.Where("owner_id =?", owner.ID).Where("name =?", repoName).First(repo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("repository not found for owner %v and repo %v", owner.Username, repoName)

		}
		return err
	}

	for _, commit := range *commitRepoInfos {
		// check if this commit already exists in our database
		existingCommit, err := s.GetCommitBySHA(commit.SHA)
		if err != nil {
			log.Println("Error in checking existing commits by sha")
			continue
		}
		if existingCommit != nil {
			// commit already exists, skip;
			log.Printf("Commit with SHA: %s already exists; skipping", existingCommit.SHA)
			continue
		}
		newCommit := &Commit{
			RepositoryName: repoName,
			SHA:            commit.SHA,
			Message:        commit.Message,
			Author:         commit.Author,
			Date:           commit.Date,
		}
		log.Printf("New commit to be created: %v", newCommit)
		err = s.DB.Create(newCommit).Error
		if err != nil {
			log.Printf("Error in saving commits with SHA: %s", newCommit.SHA)
			return err
		}
	}
	return nil
}

func (s *SqliteCommitRepository) GetCommitBySHA(sha string) (*Commit, error) {
	commit := &Commit{}
	err := s.DB.Where("sha =?", sha).First(commit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return commit, nil
}

func (s *SqliteCommitRepository) GetRepositoryCommits(repoName string) ([]*Commit, error) {
	//  logic to retrieve commit info from the database by repository name
	commits := &[]*Commit{}
	err := s.DB.Where("repository_name =?", repoName).Find(commits).Error
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return *commits, nil
}

func (s *SqliteCommitRepository) GetMostRecentCommitInRepository(repoName string) (*Commit, error) {
	commit := &Commit{}
	err := s.DB.Where("repository_name =?", repoName).Order("created_at DESC").First(commit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return commit, nil
}

func (s *SqliteCommitRepository) DeleteUntilSHA(repoName, sha string) error {
	allCommits := &[]*Commit{}

	// get all the commits in descending order of when they were created
	err := s.DB.Where("repository_name =?", repoName).Order("created_at DESC").Find(allCommits).Error
	if err != nil {
		log.Printf("Error fetching all commits in created at order: %v", err)
		return err
	}
	// we remove all the items from the most recent commits to the preferred sha we want to resr into
	for _, commit := range *allCommits {
		if commit.SHA == sha {
			break
		}
		err = s.DB.Delete(commit).Error
		if err != nil {
			log.Printf("Error deleting commits after sha %s: %v", sha, err)
			return err
		}
	}

	return nil
}

func (s *SqliteCommitRepository) AddAuthorCommitCount(author string, count int) error {
	authorCommitCount := &AuthorCommitCount{}
	err := s.DB.Where("author =?", author).First(authorCommitCount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newAuthorCommitCount := &AuthorCommitCount{
				Author:      author,
				CommitCount: count,
			}
			err = s.DB.Create(newAuthorCommitCount).Error
			if err != nil {
				log.Printf("Error creating author commit count for author %s: %v", author, err)
				return err
			}
		} else {
			log.Printf("Error fetching author commit count for author %s: %v", author, err)
			return err
		}
	}
	authorCommitCount.CommitCount += count
	return s.DB.Save(authorCommitCount).Error
}

func (s *SqliteCommitRepository) FindTopNAuthorsByCommitCounts(topN int) ([]*AuthorCommitCount, error) {
	authorCounts := &[]*AuthorCommitCount{}
	err := s.DB.Order("commit_count DESC").Limit(topN).Find(authorCounts).Error
	if err != nil {
		log.Printf("Error fetching top %d authors by commit counts: %v", topN, err)
		return nil, err
	}
	return *authorCounts, nil
}
