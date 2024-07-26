package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/midedickson/github-service/config"
	"github.com/midedickson/github-service/discovery"
	"github.com/midedickson/github-service/interface/controllers"
	"github.com/midedickson/github-service/interface/database"
	tasks "github.com/midedickson/github-service/interface/task-manager"
	"github.com/midedickson/github-service/requester"
	"github.com/midedickson/github-service/routes"
	"github.com/midedickson/github-service/usecase"
)

func main() {
	log.Println("Starting server...")

	// loading configuration file and get Database URL
	dbUrl := config.GetDBUrl()

	database.ConnectToDB(dbUrl)
	database.AutoMigrate()
	// Use a WaitGroup to manage goroutines
	var wg sync.WaitGroup

	repoRequester := requester.NewRepositoryRequester()

	// databasae repositories for each domain/service
	userRepository := database.NewSqliteUserRepository(database.DB)
	repoRepository := database.NewSqliteRepoRepository(database.DB)
	commitRepository := database.NewSqliteCommitRepository(database.DB)

	// commit manager for handling commit discovery and monitoring task execution
	commitManager := discovery.NewCommitDiscoveryService(repoRepository, repoRequester, commitRepository, config.GetCommitStartDate(), config.GetCommitEndDate())

	// repo discovery for executing tasks relating to finding repositories
	repoDiscovery := discovery.NewRepositoryDiscoveryService(repoRequester, userRepository, repoRepository, commitRepository, commitManager)

	// task manager for managing the queueing and execution of tasks
	taskManager := tasks.NewTaskManager(repoDiscovery, commitManager)

	// Usecase services for each domain/service
	userUseCase := usecase.NewUserUseCaseService(userRepository, taskManager)
	repoUseCase := usecase.NewRepoUseCaseService(repoRepository, userUseCase, taskManager)
	commitUseCase := usecase.NewCommitUseCaseService(commitRepository, repoUseCase, taskManager)

	// creation of application handler
	controller := controllers.NewController(repoRequester, userUseCase, repoUseCase, commitUseCase)

	// Starting goroutines to fetch repositories and check for updates
	wg.Add(1)
	go taskManager.GetAllRepoForUser(&wg)
	wg.Add(1)
	go taskManager.FetchNewlyRequestedRepo(&wg)
	wg.Add(1)
	go taskManager.CheckForUpdateOnAllRepo(&wg)
	wg.Add(1)
	go taskManager.HandleRequestedRepoReset(&wg)

	// start the task for updating all repositories
	go taskManager.AddSignalToCheckForUpdateOnAllRepoQueue()

	// create mux router and connect handlers to router
	r := mux.NewRouter()
	routes.ConnectRoutes(r, controller)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	server := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	log.Println("Server started on :8080")
	<-stop
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close channels to signal workers to stop
	close(taskManager.GetAllRepoForUserQueue)
	close(taskManager.FetchNewlyRequestedRepoQueue)
	close(taskManager.CheckForUpdateOnAllRepoQueue)
	close(taskManager.ResetRepositoryQueue)

	// Wait for all goroutines to complete
	wg.Wait()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
