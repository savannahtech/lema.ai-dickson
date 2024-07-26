package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/midedickson/github-service/entity"
	"github.com/midedickson/github-service/interface/controllers"
	"github.com/midedickson/github-service/test/mocks"
	"github.com/midedickson/github-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetRepositoryInfo(t *testing.T) {
	mockRepoUseCase := new(mocks.MockRepoUseCase)
	controller := controllers.NewController(nil, nil, mockRepoUseCase, nil)

	t.Run("successful fetch repository info", func(t *testing.T) {
		// Create a new HTTP request
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuser", "repo": "testrepo"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		owner := &entity.User{Username: "testuser"}
		repo := &entity.Repository{Name: "testrepo", Owner: owner}
		mockRepoUseCase.On("GetRepositoryInfo", "testuser", "testrepo").Return(repo, nil)

		http.HandlerFunc(controller.GetRepositoryInfo).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "Repository Information Fetched Successfully", response.Message)
		mockRepoUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload - missing owner", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "", "repo": "testrepo"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.GetRepositoryInfo).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("repository not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuser", "repo": "testrepox"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		mockRepoUseCase.On("GetRepositoryInfo", "testuser", "testrepox").Return(nil, nil)

		http.HandlerFunc(controller.GetRepositoryInfo).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Repository not found on Github; kindly check back again.", response.Message)
		mockRepoUseCase.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuser", "repo": "testrepoy"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		mockRepoUseCase.On("GetRepositoryInfo", "testuser", "testrepoy").Return(nil, errors.New("some error"))

		http.HandlerFunc(controller.GetRepositoryInfo).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "some error", response.Message)
		mockRepoUseCase.AssertExpectations(t)
	})
}

func TestGetRepositories(t *testing.T) {
	mockRepoUseCase := new(mocks.MockRepoUseCase)
	controller := controllers.NewController(nil, nil, mockRepoUseCase, nil)

	t.Run("successful fetch repositories", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuser"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		owner := &entity.User{Username: "testuser"}
		repositories := []*entity.Repository{
			{Name: "repo1", Owner: owner},
			{Name: "repo2", Owner: owner},
		}
		repoSearchParams := &utils.RepositorySearchParams{}
		mockRepoUseCase.On("GetUserRepositories", "testuser", repoSearchParams).Return(repositories, nil)

		http.HandlerFunc(controller.GetRepositories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "Repositories Fetched Successfully", response.Message)
		mockRepoUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload - missing owner", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": ""})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.GetRepositories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/", nil)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuserx"})
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		repoSearchParams := &utils.RepositorySearchParams{}
		mockRepoUseCase.On("GetUserRepositories", "testuserx", repoSearchParams).Return(nil, errors.New("some error"))

		http.HandlerFunc(controller.GetRepositories).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "some error", response.Message)
		mockRepoUseCase.AssertExpectations(t)
	})
}
