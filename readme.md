# GitHub Service

## Overview

GitHub Service is a Go-based application that interacts with the GitHub API to manage repositories and their commits. The service allows users to register their GitHub usernames, fetch repository information, and retrieve commit details for repositories.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Dependencies](#dependencies)
3. [API Endpoints](#api-endpoints)
4. [Usage](#usage)
5. [Video Explanation](#video-example)
6. [Running Tests](#running-tests)

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed on your machine:

- [Go](https://golang.org/doc/install) (version 1.15+)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/midedickson/github-service.git
   ```

2. Navigate to the project directory:

   ```sh
   cd github-service
   ```

3. Download the project dependencies:

   ```sh
   go mod download
   ```

4. Create new env file:
   ```sh
   cp .env.sample .env
   ```

### Running the Application

To run the application locally:

```sh
go run cmd/main.go
```

The application will start on `http://localhost:8080`.

## Dependencies

The project uses the following dependencies:

- [gorilla/mux](https://github.com/gorilla/mux) - URL router and dispatcher for Go
- [testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks
- [mockery](https://github.com/vektra/mockery) - A mock code autogenerator for Golang
- [gorm.io/gorm](https://gorm.io/) - The fantastic ORM library for Golang
- [github.com/golang/mock](https://github.com/golang/mock) - GoMock is a mocking framework for the Go programming language.

## API Endpoints

Find the documentation to the API endpoints here: https://documenter.getpostman.com/view/26825676/2sA3kPpjD1

## Usage

#### Get All Commits for a Repository:

- Commits are already stored in the database with a foreignKey to the repository with the repository name being the reference attached to the commit table.

```go
type Commit struct {
	gorm.Model
	RepositoryName string      `gorm:"repository_name"`
	Repository     *Repository `gorm:"foreignKey:RepositoryName"`
	Message        string      `gorm:"message" json:"message"`
	Author         string      `gorm:"author" json:"author"`
	Date           string      `gorm:"string" json:"date"`
	URL            string      `gorm:"html_url" json:"html_url"`
	SHA            string      `gorm:"sha" json:"sha"`
}
```

- Hence, The query for getting all the commits for a repository by Name, is a single query that returns all the commits that meet the criteria.

```go
func (s *SqliteCommitRepository) GetRepositoryCommits(repoName string) ([]*Commit, error) {
	//  logic to retrieve commit info from the database by repository name
	commits := &[]*Commit{}
	err := s.DB.Where("repository_name =?", repoName).Find(commits).Error
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return *commits, nil
}
```

- Fetching commits for a repository by repository name is achieved via this endpoint: `/{owner}/repos/{repo}/commits`

#### Get Top N Authors by Commits:

- Getting the Top N Authors by Commits Count has an interesting approach.
- Normally, it would be achieved by making a single a wide query into the entire commit and then use a group by aggregate to get the top authors.
- However, in this system commits can grow very quickly, it's an exponential growth, hence, running aggregate queries can impale performance and computational load on the Database at scale.
- Moreso, this data doesn't need to be completely real-time, reemphasizing that this system prioritizes eventual consistency.
- That said, the implementation uses a pre-aggregated data store in a different AuthorCommitCount table.

```go
type AuthorCommitCount struct {
	gorm.Model
	Author      string `gorm:"author"`
	CommitCount int    `gorm:"commit_count"`
}
```

- With this table, all authors in our system will have their row in this table holding their current total commit count.
- The row will be updated anytime there are new commits being pulled into our database.

```go
func (s *SqliteCommitRepository) AddAuthorCommitCount(author string, count int) error {
	authorCommitCount := &AuthorCommitCount{}
	err := s.DB.Where("author =?", author).First(authorCommitCount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newAuthorCommitCount := &AuthorCommitCount{
				Author:      author,
				CommitCount: count,
			}
			err = s.DB.Create(newAuthorCommitCount).Error
			if err != nil {
				log.Printf("Error creating author commit count for author %s: %v", author, err)
				return err
			}
		} else {
			log.Printf("Error fetching author commit count for author %s: %v", author, err)
			return err
		}
	}
	authorCommitCount.CommitCount += count
	return s.DB.Save(authorCommitCount).Error
}
```

- With this implementation, we only need to run a simplpe query to OrderBy the total commit count!

```go
func (s *SqliteCommitRepository) FindTopNAuthorsByCommitCounts(topN int) ([]*AuthorCommitCount, error) {
	authorCounts := &[]*AuthorCommitCount{}
	err := s.DB.Order("commit_count DESC").Limit(topN).Find(authorCounts).Error
	if err != nil {
		log.Printf("Error fetching top %d authors by commit counts: %v", topN, err)
		return nil, err
	}
	return *authorCounts, nil
}

```

- Use the `/authors/top/{top_n}` endpoint to fetch the top N authors by commit count.
  Example: Get the top 3 authors by commit.

## Video Explanation

### Folder Structure Walkthrough:

https://www.loom.com/share/67491e4aeec64cb6be3e86ef0e376176?sid=7cb07326-428f-4fc2-a9ac-f29f3b13601e

### Endpoint and Unit Test Demo:

https://www.loom.com/share/12f4fbce610a40048eab4d26d43a4bc8?sid=858c5c8d-68f4-4e08-98b0-b7ea42c370e6

## Running Tests

The project includes unit tests for the controller methods. To run the tests, use the following command:

```sh
go test ./...
```
