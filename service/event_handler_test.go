package service

import (
	"testing"

	"github.com/ameykpatil/github-data-analyzer/db"
	"github.com/ameykpatil/github-data-analyzer/db/entities"
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
	repo1 = entities.Repo{
		ID:   "441",
		Name: "Repo1",
	}
	repo2 = entities.Repo{
		ID:   "442",
		Name: "Repo2",
	}
)

func TestBuildEvents(t *testing.T) {
	tests := []struct {
		name string
		ds   db.DataStore
		exp  map[string]*Event
	}{
		{
			name: "one-to-one commits to event relation",
			ds: db.DataStore{
				ActorStore: map[string]*entities.Actor{
					actor1.ID: &actor1,
					actor2.ID: &actor2,
				},
				CommitStore: map[string]*entities.Commit{
					commit1.Sha: &commit1,
					commit2.Sha: &commit2,
				},
				EventStore: map[string]*entities.Event{
					event1.ID: &event1,
					event2.ID: &event2,
				},
				RepoStore: map[string]*entities.Repo{
					repo1.ID: &repo1,
					repo2.ID: &repo2,
				},
			},
			exp: map[string]*Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2}},
			},
		},
		{
			name: "many-to-one commits to event relation",
			ds: db.DataStore{
				ActorStore: map[string]*entities.Actor{
					actor1.ID: &actor1,
					actor2.ID: &actor2,
				},
				CommitStore: map[string]*entities.Commit{
					commit1.Sha: &commit1,
					commit2.Sha: &commit2,
					commit3.Sha: &commit3,
				},
				EventStore: map[string]*entities.Event{
					event1.ID: &event1,
					event2.ID: &event2,
				},
				RepoStore: map[string]*entities.Repo{
					repo1.ID: &repo1,
					repo2.ID: &repo2,
				},
			},
			exp: map[string]*Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2, commit3}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildEvents(&tt.ds)
			// here the commit array is inside the event struct which is inside the map
			// when we check for equality it can ignore the order of the top-level map but
			// for the slice it needs the order of the elements in the commits slice to be same
			// order of commit slice is not a guaranteed in our case because ultimately it comes from the map
			// hence we will need to check equality field by field
			for k, v := range tt.exp {
				if gotEvent, ok := got[k]; ok {
					assert.Equal(t, v.ID, gotEvent.ID)
					assert.Equal(t, v.Type, gotEvent.Type)
					assert.Equal(t, v.Actor, gotEvent.Actor)
					assert.Equal(t, v.Repo, gotEvent.Repo)
					assert.ElementsMatch(t, v.Commits, gotEvent.Commits)
				} else {
					t.Errorf("expected event with %s not in the result", k)
				}
			}
		})
	}
}
