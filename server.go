package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
)

type ErrorRequest struct {
	Language string
	Tags     []string `json:"tags"`
}

type Error struct {
	Short    string
	Full     string
	Language string
	Tags     []string
}

// GET /error?language=foo&tags=foo,bar,baz

func main() {

	missingErrErr := Error{
		Short:    "No error found matching your request",
		Full:     "No error found matching your request",
		Language: "json",
		Tags:     []string{"404"},
	}
	missingErrData, _ := json.Marshal(missingErrErr)

	errors := []Error{
		{
			Short:    "TypeError: module.__init__() takes at most 2 arguments (3 given)",
			Full:     "TypeError: module.__init__() takes at most 2 arguments (3 given)",
			Language: "python",
			Tags:     []string{"types", "module", "arity"},
		},
		{
			Short:    "expected '{', found 'type'",
			Language: "go",
			Tags:     []string{"brace", "syntax"},
		},
	}

	// TODO load real errors

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		// Default to a random *any* error
		queryParams := map[string][]string(r.URL.Query())
		var language string
		var tags []string

		if lang, ok := queryParams["lang"]; ok {
			logrus.Infof("lang is %v", lang)
			if len(lang) == 0 {
				logrus.Errorf("Expected language to be at least length 1")
			} else {
				language = lang[0]
			}
		}

		if tags, ok := queryParams["tags"]; ok {
			// TODO index oob
			tags = strings.Split(tags[0], ",")
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Could not read body: %v", err)
		} else {
			var req ErrorRequest
			err := json.Unmarshal(body, &req)
			if err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
			} else {
				if req.Language != "" {
					language = req.Language
				}
				if len(req.Tags) > 0 {
					tags = req.Tags
				}
			}
		}

		// tags all collected, now let's find any errors that match this request

		candidates := make([]scoredErr, len(errors))
		for i, err := range errors {
			candidate := scoredErr{e: err}
			if err.Language == language {
				candidate.score += 100
			}

			candidate.score += numMatch(err.Tags, tags) * 10
			candidates[i] = candidate
		}

		sort.Sort(byScore(candidates))

		w.WriteHeader(500)
		if len(candidates) == 0 {
			w.Write(missingErrData)
			return
		}
		numEqualScores := 0
		bestScore := candidates[0].score
		for ; numEqualScores < len(candidates) && candidates[numEqualScores].score == bestScore; numEqualScores++ {
		}
		choice := candidates[rand.Intn(numEqualScores)]
		logrus.Infof("Error chosen: %v", choice)
		data, err := json.Marshal(choice.e)
		if err != nil {
			logrus.Errorf("Error %v", err)
			w.Write(missingErrData)
			return
		}
		w.Write(data)
	}))
}

type scoredErr struct {
	e     Error
	score int
}

type byScore []scoredErr

func (b byScore) Len() int {
	return len(b)
}

func (b byScore) Less(x, y int) bool {
	return b[x].score > b[y].score
}

func (b byScore) Swap(x, y int) {
	b[x], b[y] = b[y], b[x]
}

func numMatch(x, y []string) int {
	num := 0
	for _, el := range x {
		for _, el2 := range y {
			if el == el2 {
				num++
			}
		}
	}
	return num
}
