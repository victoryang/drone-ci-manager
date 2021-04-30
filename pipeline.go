package main

import (
	"fmt"
	"path"

	"github.com/drone/drone-yaml/yaml"
)

const (
	PipelineKind = "pipeline"
	PipelineRunnerExec = "exec"
	PackagingWorkspace string = "/data/rolling-build/projects/"
)

type Pipeline struct {
	Project 	string
	Env 		string
	BuildCmd 	string
	Target		string
	ImageName 	string
}

func NewPipeline(project string, env string) *Pipeline {
	info := Rolling.GetBuildInfo(project)
	if info == nil {
		fmt.Println("rolling build info empty")
		return nil
	}

	return &Pipeline{
		Project: project,
		Env: env,
		BuildCmd: info.BuildCmd,
		Target: info.Target,
	}
}

func (p *Pipeline) BuildSteps() []*yaml.Container {
	steps := make([]*yaml.Container, 0)

	// Build step
	buildStep := p.CreateBuildStep()
	steps = append(steps, buildStep)

	// Packaging step
	packagingStep := p.CreatePackagingStep()
	steps = append(steps, packagingStep)

	// Publish step
	publishStep := p.CreatePublishStep()
	steps = append(steps, publishStep)

	// Clean up
	cleanUpStep := p.CreateCleanupStep()
	steps = append(steps, cleanUpStep)

	return steps
}

func (p *Pipeline) CreateBuildStep() *yaml.Container {

	buildCommands := []string {
		"bash -c \"" + "set -ex\n\n" + p.BuildCmd + "\"",
	}
	postBuildCommands := p.CreatePostBuildCommands()
	buildCommands = append(buildCommands, postBuildCommands...)

	return &yaml.Container {
		Name: "build",
		Commands: buildCommands,
	}
}

func (p *Pipeline) CreatePostBuildCommands() []string {
	from := p.Target
	to := path.Join(PackagingWorkspace, p.Project, "release-"+p.Env, p.Project+".zip")

	return []string {
		"cp -f " + from + " " + to,
	}
}

func (p *Pipeline) CreatePackagingStep() *yaml.Container {

	packagingCommand := []string {
		"cd " + path.Join(PackagingWorkspace, p.Project, "release-"+p.Env),
		"docker build -t " + p.ImageName + " .",
	}

	return &yaml.Container {
		Name: "packaging",
		Commands: packagingCommand,
	}
}

func (p *Pipeline) CreatePublishStep() *yaml.Container {

	publishCommand := []string {
		"echo $CI_JOB_TOKEN | docker login --username $CI_USER --password-stdin $CI_REGISTRY",
		"docker push " + p.ImageName,
	}

	return &yaml.Container {
		Name: "publish",
		Commands: publishCommand,
	}
}

func (p *Pipeline) CreateCleanupStep() *yaml.Container {

	cleanUpCommand := []string {
		"docker rmi " + p.ImageName,
	}

	return &yaml.Container {
		Name: "cleanup",
		Commands: cleanUpCommand,
	}
}