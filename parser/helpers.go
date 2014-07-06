package parser

import (
	"fmt"
	"os"
	"os/exec"
)

import (
	"github.com/modcloth/docker-builder/builderfile"
	"github.com/modcloth/docker-builder/parser/tag"
)

/*
IsOpenable examines the Builderfile provided to the Parser and returns a bool
indicating whether or not the file exists and openable.
*/
func (parser *Parser) IsOpenable() bool {

	file, err := os.Open(parser.filename)
	defer file.Close()

	if err != nil {
		return false
	}

	return true
}

// turns InstructionSet structs into CommandSequence structs
func (parser *Parser) commandSequenceFromInstructionSet(is *InstructionSet) *CommandSequence {
	ret := &CommandSequence{
		Commands: []*SubSequence{},
	}

	var containerCommands []exec.Cmd

	for _, v := range is.Containers {
		containerCommands = []exec.Cmd{}

		// ADD BUILD COMMANDS
		uuid, err := parser.NextUUID()
		if err != nil {
			return nil
		}

		name := fmt.Sprintf("%s/%s", v.Registry, v.Project)
		initialTag := fmt.Sprintf("%s:%s", name, uuid)
		buildArgs := []string{"docker", "build", "-t", initialTag}
		buildArgs = append(buildArgs, is.DockerBuildOpts...)
		buildArgs = append(buildArgs, ".")

		containerCommands = append(containerCommands, *&exec.Cmd{
			Path: "docker",
			Args: buildArgs,
		})

		var tagList = []string{}

		// ADD TAG COMMANDS
		for _, t := range v.Tags {
			var tagObj tag.Tag
			tagArg := map[string]string{
				"tag": t,
				"top": parser.top,
			}

			if len(t) > 4 && t[0:4] == "git:" {
				tagObj = tag.NewTag("git", tagArg)
			} else {
				tagObj = tag.NewTag("default", tagArg)
			}

			fullTag := fmt.Sprintf("%s:%s", name, tagObj.Tag())

			tagList = append(tagList, fullTag)

			buildArgs = []string{"docker", "tag"}
			buildArgs = append(buildArgs, is.DockerTagOpts...)
			buildArgs = append(buildArgs, "<IMG>", fullTag)

			containerCommands = append(containerCommands, *&exec.Cmd{
				Path: "docker",
				Args: buildArgs,
			})
		}

		// ADD PUSH COMMANDS
		if !v.SkipPush {
			for _, fullTag := range tagList {
				buildArgs = []string{"docker", "push", fullTag}

				containerCommands = append(containerCommands, *&exec.Cmd{
					Path: "docker",
					Args: buildArgs,
				})
			}
		}

		ret.Commands = append(ret.Commands, &SubSequence{
			Metadata: &SubSequenceMetadata{
				Name:       v.Name,
				Dockerfile: v.Dockerfile,
				Included:   v.Included,
				Excluded:   v.Excluded,
				UUID:       uuid,
			},
			SubCommand: containerCommands,
		})
	}

	return ret
}

func mergeGlobals(container, globals *builderfile.ContainerSection) *builderfile.ContainerSection {
	if container.Excluded == nil {
		container.Excluded = []string{}
	}

	if container.Included == nil {
		container.Included = []string{}
	}

	if container.Tags == nil {
		container.Tags = []string{}
	}

	if container.Dockerfile == "" {
		container.Dockerfile = globals.Dockerfile
	}

	if len(container.Included) == 0 && globals.Included != nil {
		container.Included = globals.Included
	}

	if len(container.Excluded) == 0 && globals.Excluded != nil {
		container.Excluded = globals.Excluded
	}

	if container.Registry == "" {
		container.Registry = globals.Registry
	}

	if container.Project == "" {
		container.Project = globals.Project
	}

	if len(container.Tags) == 0 && globals.Tags != nil {
		container.Tags = globals.Tags
	}

	container.SkipPush = container.SkipPush || globals.SkipPush

	return container
}

// turns Builderfile structs into InstructionSet structs
func (parser *Parser) instructionSetFromBuilderfileStruct(file *builderfile.Builderfile) *InstructionSet {
	ret := &InstructionSet{
		DockerBuildOpts: file.Docker.BuildOpts,
		DockerTagOpts:   file.Docker.TagOpts,
		Containers:      []builderfile.ContainerSection{},
	}

	if file.ContainerArr == nil {
		file.ContainerArr = []*builderfile.ContainerSection{}
	}

	if file.ContainerGlobals == nil {
		file.ContainerGlobals = &builderfile.ContainerSection{}
	}
	globals := file.ContainerGlobals

	for _, container := range file.ContainerArr {
		container = mergeGlobals(container, globals)
		ret.Containers = append(ret.Containers, *container)
	}

	return ret
}
