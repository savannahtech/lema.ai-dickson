package tasks

import (
	"log"
	"sync"
	"time"
)

func (t *TaskManager) GetAllRepoForUser(wg *sync.WaitGroup) {
	//  logic to fetch all repositories for the given user
	// Use the GetAllRepoForUserQueue channel to send and recieve the user to and from the worker pool
	defer wg.Done()
	for user := range t.GetAllRepoForUserQueue {

		// Handover task to repository discovery
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.repoDiscovery.GetAllUserRepositories(user)
			t.AddUserToGetAllRepoQueue(user)
		}()
	}
}

func (t *TaskManager) FetchNewlyRequestedRepo(wg *sync.WaitGroup) {
	//  logic to fetch a newly requested repo and commits for the given repository
	defer wg.Done()
	log.Println("waiting for newly requested repos...")

	for repoRequest := range t.FetchNewlyRequestedRepoQueue {
		log.Println("checking for newly requested repos...")
		wg.Add(1)
		go t.repoDiscovery.FetchNewlyRequestedRepo(repoRequest, wg)
	}

	log.Println("exiting checking for newly requested repos...")
}

func (t *TaskManager) HandleRequestedRepoReset(wg *sync.WaitGroup) {
	//  logic to fetch a newly requested repo and commits for the given repository
	defer wg.Done()
	log.Println("waiting for newly requested repos...")

	for repoResetRequest := range t.ResetRepositoryQueue {
		log.Println("checking for newly requested repos...")
		wg.Add(1)
		go t.commitManager.ResetCommitToSHA(repoResetRequest.RepoName, repoResetRequest.ResetSHA)
	}

	log.Println("exiting checking for newly requested repos...")
}

func (t *TaskManager) CheckForUpdateOnAllRepo(wg *sync.WaitGroup) {
	//  logic to check for updates on all repositories in the database
	defer wg.Done()
	for {
		_, ok := <-t.CheckForUpdateOnAllRepoQueue
		if !ok {
			log.Println("No more signal to check for updates on all repositories")
			return
		}
		err := t.repoDiscovery.CheckForUpdateOnAllRepo()
		if err != nil {
			log.Printf("Error in checking for updates on all repositories: %v", err)
		}

		// trigger the update again after 3days (currently passed as seconds)
		time.Sleep(3 * time.Second)
		go t.AddSignalToCheckForUpdateOnAllRepoQueue()
	}
}
