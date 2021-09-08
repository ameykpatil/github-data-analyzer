package repo

import (
	"container/heap"

	"github.com/ameykpatil/github-data-analyzer/service"
)

// Analyzer encapsulates functionality of analyzing the repos
type Analyzer struct {
	eventHandler service.EventHandler
	repoMap      map[string]*Repo
}

// NewAnalyzer creates a new instance of repo Analyzer
func NewAnalyzer(eventHandler service.EventHandler) *Analyzer {
	return &Analyzer{
		eventHandler: eventHandler,
		repoMap:      indexRepos(eventHandler),
	}
}

// indexRepos creates map of repo from the events
func indexRepos(eventHandler service.EventHandler) map[string]*Repo {
	repoMap := make(map[string]*Repo)
	for _, event := range eventHandler.Events {
		if event.Repo == nil {
			continue
		}
		repo, ok := repoMap[event.Repo.ID]
		if !ok {
			repo = &Repo{
				ID:             event.Repo.ID,
				Name:           event.Repo.Name,
				EventTypeCount: map[string]int{},
			}
			repoMap[repo.ID] = repo
		}
		repo.CommitCount = repo.CommitCount + len(event.Commits)
		repo.EventTypeCount[event.Type] = repo.EventTypeCount[event.Type] + 1
	}
	return repoMap
}

// GetTopRepos returns top repos based on provided limit & sort function
func (ra *Analyzer) GetTopRepos(limit uint32, fn func(ri, rj Repo) bool) []Repo {
	h := &repoHeap{less: fn}
	heap.Init(h)
	for _, repo := range ra.repoMap {
		heap.Push(h, *repo)
	}

	var repos []Repo
	for i := uint32(0); i < limit; i++ {
		repos = append(repos, heap.Pop(h).(Repo))
	}

	return repos
}
