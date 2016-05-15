package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/euank/api.error.technology/errortech"
	"github.com/euank/api.error.technology/providers"
)

type ErrorRequest struct {
	Language string   `json:"language"`
	Tags     []string `json:"tags"`
	Source   string   `json:"source"`
}

// GET /error?language=foo&tags=foo,bar,baz

var missingErrErr = errortech.Error{
	Short:    "No error found matching your request",
	Full:     "No error found matching your request",
	Language: "json",
	Tags:     []string{"404"},
}
var missingErrData, _ = json.Marshal(missingErrErr)

func main() {

	errProv := providers.NewDefaultProviders()

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Request: %v", r)

		queryParams := map[string][]string(r.URL.Query())
		var language string
		var tags []string
		var full bool
		var source string

		if lang, ok := queryParams["lang"]; ok {
			logrus.Infof("lang is %v", lang)
			if len(lang) == 0 {
				logrus.Errorf("Expected language to be at least length 1")
			} else {
				language = lang[0]
			}
		}

		if pquery, ok := queryParams["source"]; ok {
			logrus.Infof("provider is %v", pquery)
			if len(pquery) == 0 {
				logrus.Errorf("Expected at least a provider")
			} else {
				source = pquery[0]
			}
		}

		if tags, ok := queryParams["tags"]; ok {
			// TODO index oob
			tags = strings.Split(tags[0], ",")
		}

		_, full = queryParams["full"]
		if full {
			logrus.Infof("full is true")
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
				if req.Source != "" {
					source = req.Source
				}
			}
		}

		// tags all collected, now let's find any errors that match this request

		var provider providers.ErrorProvider
		for _, prov := range errProv.All() {
			if prov.Name() == source {
				provider = prov
			}
		}

		if provider == nil {
			provider = errProv.Random()
		}

		chosenErr := provider.GetError(language, tags)

		w.WriteHeader(500)
		writeErr(full, chosenErr, w)
	}))
}

func writeErr(full bool, e errortech.Error, w io.Writer) {
	if e.Full == "" {
		e.Full = e.Short
	}
	if full {
		w.Write([]byte(e.Full))
		return
	}
	data, err := json.Marshal(e)
	if err != nil {
		logrus.Errorf("Error %v", err)
		w.Write(missingErrData)
		return
	}
	w.Write(data)
}
