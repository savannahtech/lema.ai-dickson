package entity

type Repository struct {
	ID              uint
	RemoteID        int
	Owner           *User
	Name            string
	Description     string
	URL             string
	Language        string
	ForksCount      int
	StarsCount      int
	OpenIssues      int
	Watchers        int
	RemoteCreatedAt string
	RemoteUpdatedAt string
}
