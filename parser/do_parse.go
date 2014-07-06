package parser

import (
	"fmt"
	"github.com/modcloth/docker-builder/builderfile"
	"io/ioutil"
)

import (
	"github.com/BurntSushi/toml"
)

// Step 1 of Parse
func (parser *Parser) getRaw() (string, error) {
	if !parser.IsOpenable() {
		return "", fmt.Errorf("%s is not openable", parser.filename)
	}

	bytes, err := ioutil.ReadFile(parser.filename)

	if err != nil {
		return "", err
	}

	raw := string(bytes)

	return raw, nil
}

// Step 2 of Parse
func (parser *Parser) rawToStruct() (*builderfile.Builderfile, error) {
	raw, err := parser.getRaw()
	if err != nil {
		return nil, err
	}

	file := &builderfile.Builderfile{}
	if _, err := toml.Decode(raw, &file); err != nil {
		return nil, err
	}

	file.Clean()

	return file, nil
}

// Step 2.5 of Parse - handle Bobfile version

func (parser *Parser) convertBobfileVersion() (*builderfile.Builderfile, error) {
	var err error
	var fileZero *builderfile.Builderfile

	fileZero, err = parser.rawToStruct()
	if err != nil {
		return nil, err
	}

	// check version, do conversion
	if fileZero.Version == 1 {
		return fileZero, nil
	}

	builderfile.Logger(logger)
	return builderfile.Convert0to1(fileZero)
}

// Step 3 of Parse
func (parser *Parser) structToInstructionSet() (*InstructionSet, error) {
	file, err := parser.convertBobfileVersion()
	if err != nil {
		return nil, err
	}

	return parser.instructionSetFromBuilderfileStruct(file), nil
}

// Step 4 of Parse()
func (parser *Parser) instructionSetToCommandSequence() (*CommandSequence, error) {
	is, err := parser.structToInstructionSet()
	if err != nil {
		return nil, err
	}

	return parser.commandSequenceFromInstructionSet(is), nil
}

// wrapper function for the final step
func (parser *Parser) finalStep() (interface{}, error) {
	return parser.instructionSetToCommandSequence()
}

/*
Parse further parses the Builderfile struct into an InstructionSet struct,
merging the global container options into the individual container sections.
*/
func (parser *Parser) Parse() (*CommandSequence, error) {
	ret, err := parser.finalStep()
	if err != nil {
		return nil, err
	}

	return ret.(*CommandSequence), nil
}
