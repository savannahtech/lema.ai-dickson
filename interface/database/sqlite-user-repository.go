package database

import (
	"github.com/midedickson/github-service/dto"
	"gorm.io/gorm"
)

type SqliteUserRepository struct {
	DB *gorm.DB
}

func NewSqliteUserRepository(db *gorm.DB) *SqliteUserRepository {
	return &SqliteUserRepository{DB: db}
}

func (s *SqliteUserRepository) CreateUser(createUserPaylod *dto.CreateUserPayloadDTO) (*User, error) {
	// Create a user from payload
	existingUser, err := s.GetUser(createUserPaylod.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		// user already exists, update existing record;
		existingUser.FullName = createUserPaylod.FullName
		return existingUser, s.DB.Save(existingUser).Error
	}
	newUser := &User{
		Username: createUserPaylod.Username,
		FullName: createUserPaylod.FullName,
	}
	// add users into the pool to get more
	return newUser, s.DB.Create(newUser).Error
}

func (s *SqliteUserRepository) GetUser(username string) (*User, error) {
	// Get user by username
	var user User
	err := s.DB.Where("username =?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, nil
}
