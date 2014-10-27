package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rafecolton/docker-builder/job"
)

var (
	githubSupportedEvents = map[string]bool{
		"push": true,
	}
)

type githubOwner struct {
	Name string `json:"name"`
}

type githubRepository struct {
	Name  string      `json:"name"`
	Owner githubOwner `json:"owner"`
}

type githubPushPayload struct {
	Repository githubRepository `json:"repository"`
	CommitSHA  string           `json:"after"`
	Ref        string           `json:"ref"`
}

/*
Github parses a Github webhook HTTP request and returns a job.Spec.
*/
func Github(w http.ResponseWriter, req *http.Request) (int, string) {
	event := req.Header.Get("X-Github-Event")
	if !githubSupportedEvents[event] {
		logger.Errorf("Github event type %s is not supported.", event)
		return 400, "400 bad request"
	}
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		logger.Error(err)
		return 400, "400 bad request"
	}
	var payload = &githubPushPayload{}
	if err = json.Unmarshal([]byte(body), payload); err != nil {
		logger.Error(err)
		return 400, "400 bad request"
	}

	spec := &job.Spec{
		RepoOwner: payload.Repository.Owner.Name,
		RepoName:  payload.Repository.Name,
		GitRef:    strings.TrimPrefix(payload.Ref, "refs/heads/"),
	}

	return processJobHelper(spec, w, req)
}
