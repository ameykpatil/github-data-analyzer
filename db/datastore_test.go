package db

import (
	"testing"

	"github.com/ameykpatil/github-data-analyzer/db/entities"
	"github.com/stretchr/testify/assert"
)

func TestReadActors(t *testing.T) {

	expected := map[string]*entities.Actor{
		"8422699": {
			ID:       "8422699",
			Username: "Apexal",
		},
		"53201765": {
			ID:       "53201765",
			Username: "ArturoCamacho0",
		},
	}

	path := "../data/test-data"
	actorsMap, err := readActors(path)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, actorsMap)
}

func TestReadActorsError(t *testing.T) {

	expected := "open ../data/test-data1/actors.csv: no such file or directory"

	path := "../data/test-data1"
	_, err := readActors(path)

	assert.EqualError(t, err, expected)
}

func TestReadCommits(t *testing.T) {

	expected := map[string]*entities.Commit{
		"5948a6cc5255015e983a9719117c15ff197b4681": {
			Sha:     "5948a6cc5255015e983a9719117c15ff197b4681",
			Message: "Refactor member inde",
			EventID: "11185376329",
		},
		"bf7296401598660b44d8923787a2600f346f9a81": {
			Sha:     "bf7296401598660b44d8923787a2600f346f9a81",
			Message: "Refactor roadmap",
			EventID: "11185376329",
		},
	}

	path := "../data/test-data"
	commitsMap, err := readCommits(path)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, commitsMap)
}

func TestReadCommitsError(t *testing.T) {

	expected := "open ../data/test-data1/commits.csv: no such file or directory"

	path := "../data/test-data1"
	_, err := readCommits(path)

	assert.EqualError(t, err, expected)
}

func TestReadEvents(t *testing.T) {

	expected := map[string]*entities.Event{
		"11185376329": {
			ID:      "11185376329",
			Type:    "PushEvent",
			ActorID: "8422699",
			RepoID:  "224252202",
		},
		"11185376333": {
			ID:      "11185376333",
			Type:    "CreateEvent",
			ActorID: "53201765",
			RepoID:  "231161852",
		},
	}

	path := "../data/test-data"
	eventsMap, err := readEvents(path)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, eventsMap)
}

func TestReadEventsError(t *testing.T) {

	expected := "open ../data/test-data1/events.csv: no such file or directory"

	path := "../data/test-data1"
	_, err := readEvents(path)

	assert.EqualError(t, err, expected)
}

func TestReadRepos(t *testing.T) {

	expected := map[string]*entities.Repo{
		"224252202": {
			ID:   "224252202",
			Name: "DSC-RPI/dsc-portal",
		},
		"231161852": {
			ID:   "231161852",
			Name: "ArturoCamacho0/ProjectResponsive",
		},
	}

	path := "../data/test-data"
	reposMap, err := readRepos(path)

	assert.Nil(t, err)
	assert.EqualValues(t, expected, reposMap)
}

func TestReadReposError(t *testing.T) {

	expected := "open ../data/test-data1/repos.csv: no such file or directory"

	path := "../data/test-data1"
	_, err := readRepos(path)

	assert.EqualError(t, err, expected)
}
