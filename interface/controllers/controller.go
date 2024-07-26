package controllers

import (
	"github.com/midedickson/github-service/requester"
	"github.com/midedickson/github-service/usecase"
)

type Controller struct {
	requester     requester.Requester
	userUseCase   usecase.UserUseCase
	repoUsecase   usecase.RepoUseCase
	commitUsecase usecase.CommitUseCase
}

func NewController(
	requester requester.Requester,
	userUseCase usecase.UserUseCase,
	repoUsecase usecase.RepoUseCase,
	commitUsecase usecase.CommitUseCase,
) *Controller {
	return &Controller{
		requester:     requester,
		userUseCase:   userUseCase,
		repoUsecase:   repoUsecase,
		commitUsecase: commitUsecase,
	}
}
