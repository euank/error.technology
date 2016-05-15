package ondisk

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/euank/api.error.technology/errortech"
)

type Provider struct {
	errs []errortech.Error
}

var missingErrErr = errortech.Error{
	Short:    "No error found matching your request",
	Full:     "No error found matching your request",
	Language: "json",
	Tags:     []string{"404"},
}

func New() *Provider {
	errors := []errortech.Error{}

	// Load errors from disk
	files, err := ioutil.ReadDir("base_errors")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join("base_errors", file.Name()))
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		var diskErr errortech.Error
		if err := json.Unmarshal(data, &diskErr); err != nil {
			panic(err)
		}
		errors = append(errors, diskErr)
	}

	return &Provider{
		errs: errors,
	}
}

func (p *Provider) Name() string {
	return "default"
}

func (p *Provider) GetError(lang string, tags []string) errortech.Error {
	candidates := make([]scoredErr, len(p.errs))
	for i, err := range p.errs {
		candidate := scoredErr{e: err}
		if err.Language == lang {
			candidate.score += 100
		}

		candidate.score += numMatch(err.Tags, tags) * 10
		candidates[i] = candidate
	}

	sort.Sort(byScore(candidates))

	if len(candidates) == 0 {
		return missingErrErr
	}

	numEqualScores := 0
	bestScore := candidates[0].score
	for ; numEqualScores < len(candidates) && candidates[numEqualScores].score == bestScore; numEqualScores++ {
	}
	choice := candidates[rand.Intn(numEqualScores)]
	logrus.Infof("Error chosen: %v", choice)
	return choice.e
}

type scoredErr struct {
	e     errortech.Error
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
