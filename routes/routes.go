package routes

import (
	"github.com/gorilla/mux"
	"github.com/midedickson/github-service/interface/controllers"
)

func ConnectRoutes(r *mux.Router, controller *controllers.Controller) {
	r.HandleFunc("/register", controller.CreateUser).Methods("POST")
	r.HandleFunc("/{owner}/repos", controller.GetRepositories).Methods("GET")
	r.HandleFunc("/{owner}/repos/{repo}", controller.GetRepositoryInfo).Methods("GET")
	r.HandleFunc("/{owner}/repos/{repo}/commits", controller.GetRepositoryCommits).Methods("GET")
	r.HandleFunc("/{owner}/repos/{repo}/commits/reset/{reset_sha}", controller.RequestRepositoryReset).Methods("GET")
	r.HandleFunc("/authors/top/{top_n}", controller.GetTopNAuthorsByCommits).Methods("GET")
}
