package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/controllers"
	"github.com/midedickson/github-service/test/mocks"
	"github.com/midedickson/github-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	mockUserUseCase := new(mocks.MockUserUseCase)
	controller := controllers.NewController(nil, mockUserUseCase, nil, nil)

	t.Run("successful create user", func(t *testing.T) {
		payload := &dto.CreateUserPayloadDTO{
			Username: "testuser",
			FullName: "Test User",
		}
		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		user := &entity.User{Username: "testuser", FullName: "Test User", ID: 1}
		mockUserUseCase.On("CreateUser", payload).Return(user, nil)

		http.HandlerFunc(controller.CreateUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "user created successfully", response.Message)
		mockUserUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"invalid": "payload"}`)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.CreateUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("create user error", func(t *testing.T) {
		payload := &dto.CreateUserPayloadDTO{
			Username: "testuserx",
			FullName: "Test User",
		}
		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadBytes))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		mockUserUseCase.On("CreateUser", payload).Return(nil, errors.New("some error"))

		http.HandlerFunc(controller.CreateUser).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		mockUserUseCase.AssertExpectations(t)
	})
}
