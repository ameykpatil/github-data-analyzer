package entities

// Event entity denotes a record in events csv
type Event struct {
	ID      string
	Type    string
	ActorID string
	RepoID  string
}
