package repo

import (
	"testing"

	"github.com/ameykpatil/github-data-analyzer/db/entities"
	"github.com/ameykpatil/github-data-analyzer/service"
	"github.com/stretchr/testify/assert"
)

var (
	actor1 = entities.Actor{
		ID:       "111",
		Username: "Actor1",
	}
	actor2 = entities.Actor{
		ID:       "112",
		Username: "Actor2",
	}
	commit1 = entities.Commit{
		Sha:     "221",
		Message: "Message 1",
		EventID: "331",
	}
	commit2 = entities.Commit{
		Sha:     "222",
		Message: "Message 2",
		EventID: "332",
	}
	commit3 = entities.Commit{
		Sha:     "223",
		Message: "Message 3",
		EventID: "332",
	}
	event1 = entities.Event{
		ID:      "331",
		Type:    "PushEvent",
		ActorID: "111",
		RepoID:  "441",
	}
	event2 = entities.Event{
		ID:      "332",
		Type:    "CreateEvent",
		ActorID: "112",
		RepoID:  "442",
	}
	event3 = entities.Event{
		ID:      "333",
		Type:    "CreateEvent",
		ActorID: "112",
		RepoID:  "442",
	}
	event4 = entities.Event{
		ID:      "334",
		Type:    "WatchEvent",
		ActorID: "112",
		RepoID:  "442",
	}
	repo1 = entities.Repo{
		ID:   "441",
		Name: "Repo1",
	}
	repo2 = entities.Repo{
		ID:   "442",
		Name: "Repo2",
	}
)

func TestIndexRepos(t *testing.T) {

	tests := []struct {
		name   string
		events map[string]*service.Event
		exp    map[string]*Repo
	}{
		{
			name: "single commit & event for repo",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2}},
			},
			exp: map[string]*Repo{
				repo1.ID: {repo1.ID, repo1.Name, 1, map[string]int{"PushEvent": 1}},
				repo2.ID: {repo2.ID, repo2.Name, 1, map[string]int{"CreateEvent": 1}},
			},
		},
		{
			name: "multiple commits for repo",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2, commit3}},
			},
			exp: map[string]*Repo{
				repo1.ID: {repo1.ID, repo1.Name, 1, map[string]int{"PushEvent": 1}},
				repo2.ID: {repo2.ID, repo2.Name, 2, map[string]int{"CreateEvent": 1}},
			},
		},
		{
			name: "multiple events for repo",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2, commit3}},
				event3.ID: {event3.ID, event3.Type, &actor2, &repo2, []entities.Commit{}},
				event4.ID: {event4.ID, event4.Type, &actor2, &repo2, []entities.Commit{}},
			},
			exp: map[string]*Repo{
				repo1.ID: {repo1.ID, repo1.Name, 1, map[string]int{"PushEvent": 1}},
				repo2.ID: {repo2.ID, repo2.Name, 2, map[string]int{"CreateEvent": 2, "WatchEvent": 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventHandler := service.EventHandler{DataStore: nil, Events: tt.events}
			got := indexRepos(eventHandler)
			for k, v := range tt.exp {
				if gotRepo, ok := got[k]; ok {
					assert.Equal(t, v.ID, gotRepo.ID)
					assert.Equal(t, v.Name, gotRepo.Name)
					assert.Equal(t, v.CommitCount, gotRepo.CommitCount)
					assert.EqualValues(t, v.EventTypeCount, gotRepo.EventTypeCount)
				} else {
					t.Errorf("expected repo with id %s not in the result", k)
				}
			}
		})
	}

}
