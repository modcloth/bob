package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/modcloth/go-fileutils"
	"github.com/onsi/gocleanup"

	"github.com/rafecolton/docker-builder/job"
)

const (
	// AsyncSuccessCode is the http status code returned for a successful async job
	AsyncSuccessCode = 202

	// AsyncSuccessMessage is the generic message returned for a successful async job
	AsyncSuccessMessage = "202 accepted"

	// SyncSuccessCode is the http status code returned for a successful synchronous job
	SyncSuccessCode = 201

	// SyncSuccessMessage is the generic message returned for a successful synchronous job
	SyncSuccessMessage = "201 created"
)

var logger *logrus.Logger
var apiToken string
var testMode bool

//Logger sets the (global) logger for the webhook package
func Logger(l *logrus.Logger) {
	logger = l
}

//APIToken sets the (global) apiToken for the webhook package
func APIToken(t string) {
	apiToken = t
}

//TestMode sets the (global) testMode variable for the webhook package
func TestMode(b bool) {
	testMode = b
}

func processJobHelper(spec *job.Spec, w http.ResponseWriter, req *http.Request) (int, string) {
	// If tests are running, don't actually attempt to build containers, just return success.
	// This is meant to allow testing ot the HTTP interactions for the webhooks
	if testMode {
		if spec.Sync {
			return SyncSuccessCode, SyncSuccessMessage
		}
		return AsyncSuccessCode, AsyncSuccessMessage
	}

	if err := spec.Validate(); err != nil {
		return 412, "412 precondition failed"
	}

	workdir, err := ioutil.TempDir("", "docker-build-worker")
	if err != nil {
		return 500, "500 internal server error"
	}

	gocleanup.Register(func() {
		fileutils.RmRF(workdir)
	})

	jobConfig := &job.Config{
		Logger:         logger,
		Workdir:        workdir,
		GitHubAPIToken: apiToken,
	}

	j := job.NewJob(jobConfig, spec, req)

	// if sync
	if spec.Sync {
		if err = j.Process(); err != nil {
			return 417, `{"error": "` + err.Error() + `"}`
		}
		retBytes, err := json.Marshal(j)
		if err != nil {
			return 417, `{"error": "` + err.Error() + `"}`
		}

		return SyncSuccessCode, string(retBytes)
	}

	// if async
	go j.Process()

	retBytes, err := json.Marshal(j)
	if err != nil {
		return 409, `{"error": "` + err.Error() + `"}`
	}

	return AsyncSuccessCode, string(retBytes)
}
