package user

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
		Type:    "PullRequestEvent",
		ActorID: "111",
		RepoID:  "441",
	}
	event2 = entities.Event{
		ID:      "332",
		Type:    "ForkEvent",
		ActorID: "112",
		RepoID:  "442",
	}
	event3 = entities.Event{
		ID:      "333",
		Type:    "ForkEvent",
		ActorID: "112",
		RepoID:  "442",
	}
	event4 = entities.Event{
		ID:      "334",
		Type:    "DeleteEvent",
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
		exp    map[string]*User
	}{
		{
			name: "single commit & event for users",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2}},
			},
			exp: map[string]*User{
				actor1.ID: {actor1.ID, actor1.Username, 1, map[string]int{"PullRequestEvent": 1}},
				actor2.ID: {actor2.ID, actor2.Username, 1, map[string]int{"ForkEvent": 1}},
			},
		},
		{
			name: "multiple commits for users",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2, commit3}},
			},
			exp: map[string]*User{
				actor1.ID: {actor1.ID, actor1.Username, 1, map[string]int{"PullRequestEvent": 1}},
				actor2.ID: {actor2.ID, actor2.Username, 2, map[string]int{"ForkEvent": 1}},
			},
		},
		{
			name: "multiple events for repo",
			events: map[string]*service.Event{
				event1.ID: {event1.ID, event1.Type, &actor1, &repo1, []entities.Commit{commit1}},
				event2.ID: {event2.ID, event2.Type, &actor2, &repo2, []entities.Commit{commit2}},
				event3.ID: {event3.ID, event3.Type, &actor2, &repo2, []entities.Commit{}},
				event4.ID: {event4.ID, event4.Type, &actor2, &repo2, []entities.Commit{}},
			},
			exp: map[string]*User{
				actor1.ID: {actor1.ID, actor1.Username, 1, map[string]int{"PullRequestEvent": 1}},
				actor2.ID: {actor2.ID, actor2.Username, 1, map[string]int{"ForkEvent": 2, "DeleteEvent": 1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventHandler := service.EventHandler{DataStore: nil, Events: tt.events}
			got := indexUsers(eventHandler)
			for k, v := range tt.exp {
				if gotUser, ok := got[k]; ok {
					assert.Equal(t, v.ID, gotUser.ID)
					assert.Equal(t, v.Username, gotUser.Username)
					assert.Equal(t, v.CommitCount, gotUser.CommitCount)
					assert.EqualValues(t, v.EventTypeCount, gotUser.EventTypeCount)
				} else {
					t.Errorf("expected user with id %s not in the result", k)
				}
			}
		})
	}

}
