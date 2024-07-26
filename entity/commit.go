package entity

type Commit struct {
	ID         uint
	Repository *Repository
	Message    string
	Author     string
	Date       string
	URL        string
	SHA        string
}
