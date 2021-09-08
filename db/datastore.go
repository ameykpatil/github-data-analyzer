package db

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/ameykpatil/github-data-analyzer/db/entities"
)

// DataStore is a collection of entities read from files
type DataStore struct {
	ActorStore  map[string]*entities.Actor
	CommitStore map[string]*entities.Commit
	EventStore  map[string]*entities.Event
	RepoStore   map[string]*entities.Repo
}

// NewDataStore reads files & creates an instance of DataStore
func NewDataStore(path string) (*DataStore, error) {
	actorStore, err := readActors(path)
	if err != nil {
		return nil, err
	}

	commitStore, err := readCommits(path)
	if err != nil {
		return nil, err
	}

	eventStore, err := readEvents(path)
	if err != nil {
		return nil, err
	}

	repoStore, err := readRepos(path)
	if err != nil {
		return nil, err
	}

	return &DataStore{
		ActorStore:  actorStore,
		CommitStore: commitStore,
		EventStore:  eventStore,
		RepoStore:   repoStore,
	}, nil
}

func readActors(path string) (map[string]*entities.Actor, error) {
	actorStore := make(map[string]*entities.Actor)

	in, err := os.Open(path + "/actors.csv")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	reader := csv.NewReader(in)

	// read & skip the first header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		actor := &entities.Actor{
			ID:       record[0],
			Username: record[1],
		}
		actorStore[actor.ID] = actor
	}

	return actorStore, nil
}

func readCommits(path string) (map[string]*entities.Commit, error) {
	commitStore := make(map[string]*entities.Commit)

	in, err := os.Open(path + "/commits.csv")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	reader := csv.NewReader(in)

	// read & skip the first header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		commit := &entities.Commit{
			Sha:     record[0],
			Message: record[1],
			EventID: record[2],
		}
		commitStore[commit.Sha] = commit
	}

	return commitStore, nil
}

func readEvents(path string) (map[string]*entities.Event, error) {
	eventStore := make(map[string]*entities.Event)

	in, err := os.Open(path + "/events.csv")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	reader := csv.NewReader(in)

	// read & skip the first header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		event := &entities.Event{
			ID:      record[0],
			Type:    record[1],
			ActorID: record[2],
			RepoID:  record[3],
		}
		eventStore[event.ID] = event
	}

	return eventStore, nil
}

func readRepos(path string) (map[string]*entities.Repo, error) {
	repoStore := make(map[string]*entities.Repo)

	in, err := os.Open(path + "/repos.csv")
	if err != nil {
		return nil, err
	}
	defer in.Close()

	reader := csv.NewReader(in)

	// read & skip the first header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		repo := &entities.Repo{
			ID:   record[0],
			Name: record[1],
		}
		repoStore[repo.ID] = repo
	}

	return repoStore, nil
}
