package parser

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

/*
Lint parses a builderfile and returns either nil if the file was parsed
successfully or an error indicating that parsing failed and the file is
invalid.
*/
func (parser *Parser) Lint() error {
	//TODO: this should probably call Parse instead of the unexported method rawToStruct
	_, err := parser.rawToStruct()
	return err
}

/*
AssertLint is like Lint except that instead of returning an nil/error to
indicate success/failure, it exits nonzero if linting fails.
*/
func (parser *Parser) AssertLint() int {
	if !parser.IsOpenable() {
		if parser.filename == "" {
			parser.printLintFailMessage(
				errors.New("no file provided for linting"),
			)
		} else {
			parser.printLintFailMessage(
				errors.New("unable to open file"),
			)
		}
		return 17
	}

	if err := parser.Lint(); err != nil {
		parser.printLintFailMessage(err)
		return 5
	}

	parser.printLintSuccessMessage()
	return 0
}

// helper functions
func (parser *Parser) printLintSuccessMessage() {
	parser.WithFields(logrus.Fields{"file": parser.filename}).Info("file is a valid Bobfile")
}

func (parser *Parser) printLintFailMessage(err error) {
	var errMsg string
	fields := logrus.Fields{"error": err}

	if parser.filename == "" {
		errMsg = "no file provided for linting"
	} else {
		fields["filename"] = parser.filename
		errMsg = "file provided is not a valid Bobfile"
	}
	parser.WithFields(fields).Error(errMsg)
}
