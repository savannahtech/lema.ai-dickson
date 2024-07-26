package entity

type AuthorCommitCount struct {
	Author      string `json:"author"`
	CommitCount int    `json:"commitCount"`
}
