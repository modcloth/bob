package config

import (
	"fmt"
	"os"
)

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/onsi/gocleanup"
	builderlogger "github.com/rafecolton/bob/log"
)

var (
	parser *flags.Parser
	opts   Options
)

/*
Usage is like running the builder with -h/--help - it simply prints the usage
message to stderr.
*/
func Usage() {
	parser.WriteHelp(os.Stderr)

}

/*
Runtime is a struct of convenience, used for keeping track of our conf options
(i.e. passed on the command line or specified otherwise) as well as other
useful, global-ish things.
*/
type Runtime struct {
	builderlogger.Logger
	Options
}

/*
NewRuntime returns a new Runtime struct instance that contains all of the
global-ish things specific to this invokation of builder.
*/
func NewRuntime() *Runtime {
	parser = flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		arg1 := os.Args[1]
		if arg1 == "-h" || arg1 == "--help" {
			gocleanup.Exit(0)
		} else {
			fmt.Println("Unable to parse args")
			gocleanup.Exit(3)
		}
	}

	logger := builderlogger.Initialize(opts.Quiet)

	runtime := &Runtime{
		Options: opts,
		Logger:  logger,
	}

	return runtime
}

/*
Options are our command line options, set using the
https://github.com/jessevdk/go-flags library.
*/
type Options struct {
	// Inform and Exit
	Version     bool `short:"v" description:"Print version and exit"`
	VersionFull bool `long:"version" description:"Print long version and exit"`
	Branch      bool `long:"branch" description:"Print branch and exit"`
	Rev         bool `long:"rev" description:"Print revision and exit"`

	// Runtime Options
	Quiet bool `short:"q" long:"quiet" description:"Produce no output, only exit codes" default:"false"`

	// Features
	Lintfile    string `short:"l" long:"lint" description:"Lint the provided file. Compatible with -q/--quiet"`
	Builderfile string `short:"b" long:"build" description:"The configuration file for Builder"`
}
