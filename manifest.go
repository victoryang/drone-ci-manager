package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	yamlv2 "gopkg.in/yaml.v2"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-yaml/yaml"
)

type Manifest struct {
	Env 		string
	Timestamp 	string
	Version 	string
	Pipelines 	[]*Pipeline
}

func NewManifest(repoInfo *drone.Repo, buildInfo *drone.Build) (*Manifest,error){
	branch := strings.TrimPrefix(buildInfo.Ref, "refs/heads/")
	var env string
	switch branch {
	case "staging":
		env = "staging"
	case "release":
		env = "release"
	case "prod":
		env = "rc"
	case "sep":
		env = "sep"
	default:
		fmt.Println("env not supported: ", env)
		return nil, errors.New("env not supported")
	}

	timestamp := strconv.FormatInt(buildInfo.Created, 10)
	version := buildInfo.After[:8]

	projects,err := GetProjectsByUrl(repoInfo.SSHURL)
	if err!=nil {
		return nil, err
	}

	pipelines := make([]*Pipeline, 0)
	for _,proj := range projects {
		p := NewPipeline(proj, env)

		from := GetDockerfileFromBytes(proj, env)
		tag := timestamp + "_" + version + "_" + branch + "_" + from
		p.ImageName = BuildImageName(proj, tag)

		pipelines = append(pipelines, p)
	}

	if len(pipelines) == 0 {
		return nil,nil
	}

	return &Manifest {
		Env: env,
		Timestamp: timestamp,
		Version: version,
		Pipelines: pipelines,
	}, nil
}

func (m *Manifest) Compile() (string, error) {
	var content string = ""
	for _, p :=range m.Pipelines {
		steps := p.BuildSteps()
		if len(steps)==0 {
			return "", errors.New("create steps fail")
		}

		pipeline := &yaml.Pipeline {
			Kind: PipelineKind,
			Type: PipelineRunnerExec,
			Name: p.Project,
			Steps: steps,
		}

		d, err := yamlv2.Marshal(pipeline)
		if err!=nil {
			fmt.Println("marshall error:", err)
			return "", err
		}

		content = content + fmt.Sprintf("---\n%v\n", string(d))
	}

	return content, nil
}