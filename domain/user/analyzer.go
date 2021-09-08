package user

import (
	"container/heap"

	"github.com/ameykpatil/github-data-analyzer/service"
)

// Analyzer encapsulates functionality of analyzing the users
type Analyzer struct {
	eventHandler service.EventHandler
	userMap      map[string]*User
}

// NewAnalyzer creates a new instance of user Analyzer
func NewAnalyzer(eventHandler service.EventHandler) *Analyzer {
	return &Analyzer{
		eventHandler: eventHandler,
		userMap:      indexUsers(eventHandler),
	}
}

// indexRepos creates map of users from the events
func indexUsers(eventHandler service.EventHandler) map[string]*User {
	userMap := make(map[string]*User)
	for _, event := range eventHandler.Events {
		if event.Actor == nil {
			continue
		}
		user, ok := userMap[event.Actor.ID]
		if !ok {
			user = &User{
				ID:             event.Actor.ID,
				Username:       event.Actor.Username,
				EventTypeCount: map[string]int{},
			}
			userMap[user.ID] = user
		}
		user.CommitCount = user.CommitCount + len(event.Commits)
		user.EventTypeCount[event.Type] = user.EventTypeCount[event.Type] + 1
	}

	return userMap
}

// GetTopUsers returns top users based on provided limit & sort function
func (ua *Analyzer) GetTopUsers(limit uint32, fn func(i, j User) bool) []User {
	h := &userHeap{less: fn}
	heap.Init(h)
	for _, user := range ua.userMap {
		heap.Push(h, *user)
	}

	var users []User
	for i := uint32(0); i < limit; i++ {
		users = append(users, heap.Pop(h).(User))
	}

	return users
}
