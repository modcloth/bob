package server

import (
	"net/http"
	"time"

	"github.com/rafecolton/docker-builder/job"
	"github.com/rafecolton/docker-builder/server/webhook"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/rafecolton/vauth"
)

const (
	// BuildRoute is the route used to POST docker builds from JSON
	BuildRoute = "/docker-build"

	// GitHubRoute is the route for GitHub webhooks
	GitHubRoute = "/docker-build/github"

	// HealthRoute is the route for health checks
	HealthRoute = "/health"

	// JobRoute is the route for job operations (various routes, see docs for more info)
	JobRoute = "/jobs"

	// TravisRoute is the route for TravisCI webhooks
	TravisRoute = "/docker-build/travis"
)

var logger *logrus.Logger
var server *martini.ClassicMartini
var skipLogging = map[string]bool{
	"/health": true,
}

//Logger sets the (global) logger for the server package
func Logger(l *logrus.Logger) {
	logger = l
}

// Serve sets everything up and runs the docker-builder server
func Serve(context *cli.Context) {
	// set vars
	setVarsFromContext(context)

	// set up auth functions
	if shouldBasicAuth {
		basicAuthFunc = auth.Basic(un, pwd)
	} else {
		basicAuthFunc = func(http.ResponseWriter, *http.Request) {}
	}
	if shouldTravisAuth {
		travisAuthFunc = vauth.TravisCI(travisToken)
	}
	if shouldGitHubAuth {
		githubAuthFunc = vauth.GitHub(githubSecret)
	}

	// configure webhooks
	webhook.Logger(logger)
	webhook.APIToken(apiToken)

	server = setupServer()

	if shouldTravis {
		server.Post(TravisRoute, travisAuthFunc, webhook.Travis)
	}
	if shouldGitHub {
		server.Post(GitHubRoute, githubAuthFunc, webhook.Github)
	}

	// base routes
	server.Get(HealthRoute, func() (int, string) { return 200, "200 OK" })
	server.Post(BuildRoute, basicAuthFunc, webhook.DockerBuild)

	// job control routes
	server.Group(JobRoute, func(r martini.Router) {
		r.Get("/:id", job.Get)
		r.Get("/:id/tail", job.TailN)
		r.Post("", webhook.DockerBuild)
		r.Get("", job.GetAll)
	}, basicAuthFunc)

	// start server
	http.ListenAndServe(portString, server)
}

func setupServer() *martini.ClassicMartini {
	router := martini.NewRouter()
	server := martini.New()
	server.Use(martini.Recovery())
	server.Use(requestLogger)
	server.MapTo(router, (*martini.Routes)(nil))
	server.Action(router.Handle)
	return &martini.ClassicMartini{server, router}
}

func requestLogger(res http.ResponseWriter, req *http.Request, c martini.Context) {
	if skipLogging[req.URL.Path] {
		return
	}

	start := time.Now()

	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = req.RemoteAddr
		}
	}

	logger.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

	rw := res.(martini.ResponseWriter)
	c.Next()

	logger.Printf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
}
