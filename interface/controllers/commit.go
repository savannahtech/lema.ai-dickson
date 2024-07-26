package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/midedickson/github-service/utils"
)

func (c *Controller) GetRepositoryCommits(w http.ResponseWriter, r *http.Request) {
	repoName, err := utils.GetPathParam(r, "repo")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if repoName == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	commits, err := c.commitUsecase.GetRepositoryCommits(repoName)
	if err != nil {
		log.Printf("%v", err)
		utils.Dispatch500Error(w, err)
		return
	}
	utils.Dispatch200(w, "Repository Commits Fetched Successfully", commits)
}

func (c *Controller) RequestRepositoryReset(w http.ResponseWriter, r *http.Request) {
	owner, err := utils.GetPathParam(r, "owner")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if owner == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	repoName, err := utils.GetPathParam(r, "repo")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if repoName == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	resetSHA, err := utils.GetPathParam(r, "reset_sha")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	err = c.commitUsecase.MakeRepoResetRequest(owner, repoName, resetSHA)
	if err != nil {
		log.Printf("Error occured while trying to make a reset repo request: %v", err)
		utils.Dispatch500Error(w, err)
		return
	}
	utils.Dispatch200(w, "Reset Request sent successfully", nil)
}

func (c *Controller) GetTopNAuthorsByCommits(w http.ResponseWriter, r *http.Request) {
	topNString, err := utils.GetPathParam(r, "top_n")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	topN, err := strconv.Atoi(topNString)
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if topN <= 0 {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	authors, err := c.commitUsecase.GetTopNAuthorsByCommits(topN)
	if err != nil {
		log.Printf("Error in getting top n authors: %v", err)
		utils.Dispatch500Error(w, err)
		return
	}
	utils.Dispatch200(w, "Top Authors by Commits Fetched Successfully", authors)
}
