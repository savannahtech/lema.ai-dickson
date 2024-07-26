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

func TestGetRepositoryCommits(t *testing.T) {
	mockCommitUseCase := new(mocks.MockCommitUseCase)
	controller := controllers.NewController(nil, nil, nil, mockCommitUseCase)

	t.Run("successful fetch repository commits", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}/commits", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuserx", "repo": "testrepo"})

		rr := httptest.NewRecorder()
		commits := []*entity.Commit{
			{Message: "Initial commit", Author: "testuserx"},
			{Message: "Added new feature", Author: "testuserx"},
		}
		mockCommitUseCase.On("GetRepositoryCommits", "testrepo").Return(commits, nil)

		http.HandlerFunc(controller.GetRepositoryCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "Repository Commits Fetched Successfully", response.Message)
		mockCommitUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload - missing repo", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos//commits", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuserx"})

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.GetRepositoryCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/{owner}/repos/{repo}/commits", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"owner": "testuserx", "repo": "testrepox"})

		rr := httptest.NewRecorder()

		mockCommitUseCase.On("GetRepositoryCommits", "testrepox").Return(nil, errors.New("some error"))

		http.HandlerFunc(controller.GetRepositoryCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "some error", response.Message)
		mockCommitUseCase.AssertExpectations(t)
	})
}

func TestRequestRepositoryReset(t *testing.T) {
	mockCommitUseCase := new(mocks.MockCommitUseCase)
	controller := controllers.NewController(nil, nil, nil, mockCommitUseCase)

	t.Run("successful repository reset request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/{owner}/repos/{repo}/reset/{reset_sha}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{
			"owner":     "testuser",
			"repo":      "testrepo",
			"reset_sha": "abcdef123456",
		})

		rr := httptest.NewRecorder()
		mockCommitUseCase.On("MakeRepoResetRequest", "testuser", "testrepo", "abcdef123456").Return(nil)

		http.HandlerFunc(controller.RequestRepositoryReset).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "Reset Request sent successfully", response.Message)
		assert.Nil(t, response.Data)
		mockCommitUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload - missing owner", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/{owner}/repos/{repo}/reset/{reset_sha}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{
			"owner":     "",
			"repo":      "testrepo",
			"reset_sha": "abcdef123456",
		})

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.RequestRepositoryReset).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("invalid payload - missing repo", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/{owner}/repos/{repo}/reset/{reset_sha}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{
			"owner":     "testuserx",
			"repo":      "",
			"reset_sha": "abcdef123456",
		})

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.RequestRepositoryReset).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/{owner}/repos/{repo}/reset/{reset_sha}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{
			"owner":     "testuserx",
			"repo":      "testrepo",
			"reset_sha": "abcdef123456",
		})

		rr := httptest.NewRecorder()

		mockCommitUseCase.On("MakeRepoResetRequest", "testuserx", "testrepo", "abcdef123456").Return(errors.New("some error"))

		http.HandlerFunc(controller.RequestRepositoryReset).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "some error", response.Message)
		mockCommitUseCase.AssertExpectations(t)
	})
}

func TestGetTopNAuthorsByCommits(t *testing.T) {
	mockCommitUseCase := new(mocks.MockCommitUseCase)
	controller := controllers.NewController(nil, nil, nil, mockCommitUseCase)

	t.Run("successful fetch top N authors by commits", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/authors/top/{top_n}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"top_n": "3"})

		rr := httptest.NewRecorder()
		authors := []*entity.AuthorCommitCount{
			{Author: "Author1", CommitCount: 50},
			{Author: "Author2", CommitCount: 30},
			{Author: "Author3", CommitCount: 20},
		}
		mockCommitUseCase.On("GetTopNAuthorsByCommits", 3).Return(authors, nil)

		http.HandlerFunc(controller.GetTopNAuthorsByCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, true, response.Success)
		assert.Equal(t, "Top Authors by Commits Fetched Successfully", response.Message)
		mockCommitUseCase.AssertExpectations(t)
	})

	t.Run("invalid payload - non-integer top_n", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/authors/top/{top_n}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"top_n": "abc"})

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.GetTopNAuthorsByCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("invalid payload - top_n less than or equal to zero", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/authors/top/{top_n}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"top_n": "-1"})

		rr := httptest.NewRecorder()

		http.HandlerFunc(controller.GetTopNAuthorsByCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "Invalid Payload", response.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/authors/top/{top_n}", nil)
		assert.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"top_n": "4"})

		rr := httptest.NewRecorder()

		mockCommitUseCase.On("GetTopNAuthorsByCommits", 4).Return(nil, errors.New("some error"))

		http.HandlerFunc(controller.GetTopNAuthorsByCommits).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var response utils.APIResponse
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, false, response.Success)
		assert.Equal(t, "some error", response.Message)
		mockCommitUseCase.AssertExpectations(t)
	})
}
