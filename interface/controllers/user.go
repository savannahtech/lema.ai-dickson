package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/utils"
)

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create user logic
	var createUserPayload dto.CreateUserPayloadDTO
	err := json.NewDecoder(r.Body).Decode(&createUserPayload)
	if err != nil {
		log.Printf("Error decoding create user payload: %v", err)
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if createUserPayload.Username == "" || createUserPayload.FullName == "" {
		utils.Dispatch400Error(w, "Invalid Payload", nil)
		return
	}

	user, err := c.userUseCase.CreateUser(&createUserPayload)
	if err != nil {
		log.Printf("Error occured while running %v", err)
		utils.Dispatch500Error(w, err)
		return
	}

	utils.Dispatch200(w, "user created successfully", user)
}
