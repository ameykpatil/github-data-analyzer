package service

import (
	"github.com/ameykpatil/github-data-analyzer/db"
	"github.com/ameykpatil/github-data-analyzer/db/entities"
)

// Event encapsulates all aggregated details related to an event
type Event struct {
	ID      string
	Type    string
	Actor   *entities.Actor
	Repo    *entities.Repo
	Commits []entities.Commit
}

// EventHandler is responsible for building events
type EventHandler struct {
	DataStore *db.DataStore
	Events    map[string]*Event
}

// NewEventHandler creates instance of EventHandler
func NewEventHandler(dataStore *db.DataStore) *EventHandler {
	events := BuildEvents(dataStore)
	return &EventHandler{
		DataStore: dataStore,
		Events:    events,
	}
}

// BuildEvents builds the events from the data-store
func BuildEvents(dataStore *db.DataStore) map[string]*Event {

	eventsMap := make(map[string]*Event)

	for _, eventRec := range dataStore.EventStore {
		event := Event{
			ID:    eventRec.ID,
			Type:  eventRec.Type,
			Actor: dataStore.ActorStore[eventRec.ActorID],
			Repo:  dataStore.RepoStore[eventRec.RepoID],
		}

		eventsMap[event.ID] = &event
	}

	for _, commitRec := range dataStore.CommitStore {
		if event, ok := eventsMap[commitRec.EventID]; ok {
			if event.Commits == nil {
				event.Commits = []entities.Commit{}
			}
			event.Commits = append(event.Commits, *commitRec)
		}
	}

	return eventsMap
}
