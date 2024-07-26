package tasks

import (
	"github.com/midedickson/github-service/entity"
)

type Task interface {
	AddUserToGetAllRepoQueue(user *entity.User)
	AddRequestToFetchNewlyRequestedRepoQueue(username, repoName string)
	AddRequestToResetRepositoryQueue(repoName, resetSHA string)
}
