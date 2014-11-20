package server

import (
	"fmt"
	"net/http"

	"github.com/rafecolton/docker-builder/conf"
	"github.com/rafecolton/docker-builder/job"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
)

var apiToken, githubSecret, portString, pwd, travisToken, un string
var port int
var skipPush bool
var shouldTravis, shouldGitHub bool
var shouldBasicAuth, shouldTravisAuth, shouldGitHubAuth bool

var basicAuthFunc martini.Handler
var travisAuthFunc = func(http.ResponseWriter, *http.Request) {}
var githubAuthFunc = func(http.ResponseWriter, *http.Request) {}

func setVarsFromContext(c *cli.Context) {
	config := conf.Config
	/// lowest priority

	// ENV
	un = config.Username
	pwd = config.Password
	apiToken = config.APIToken
	travisToken = config.TravisToken
	githubSecret = config.GitHubSecret
	port = config.Port

	// command line
	cliUn := c.String("username")
	cliPwd := c.String("password")
	cliAPIToken := c.String("api-token")
	cliTravisToken := c.String("travis-token")
	cliGitHubSecret := c.String("github-secret")
	cliPort := c.Int("port")

	if cliTravisToken != "" {
		travisToken = cliTravisToken
	}

	if cliGitHubSecret != "" {
		githubSecret = cliGitHubSecret
	}

	// if username passed on command line, use cl one instead
	if cliUn != "" {
		un = cliUn
	}

	// if password passed on command line, used cl one instead
	if cliPwd != "" {
		pwd = cliPwd
	}

	// get api token
	if cliAPIToken != "" {
		apiToken = cliAPIToken
	}

	//set port
	if cliPort != 0 {
		port = cliPort
	}

	// get port
	portString = fmt.Sprintf(":%d", port)

	// get skip-push
	skipPush = c.Bool("skip-push") || config.SkipPush

	// check if should travis
	shouldTravis = !c.Bool("no-travis") && !config.NoTravis

	// check if should github
	shouldGitHub = !c.Bool("no-github") && !config.NoGitHub

	shouldBasicAuth = (un != "" && pwd != "")
	shouldTravisAuth = (travisToken != "")
	shouldGitHubAuth = (githubSecret != "")

	/// highest priority

	job.SkipPush = skipPush
}
