package entities

// Commit entity denotes a record in commits csv
type Commit struct {
	Sha     string
	Message string
	EventID string
}
