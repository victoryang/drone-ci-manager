package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
)

type BuildInfo struct {
	Project 	string
	Env 		string
	Tag 		string
}

func GetTagFromBuildInfo(proj string, buildInfo *drone.Build) *BuildInfo {
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
		return nil
	}

	timestamp := strconv.FormatInt(buildInfo.Created, 10)
	version := buildInfo.After[:8]

	from := GetDockerfileFromBytes(proj, env)
	tag := timestamp + "_" + version + "_" + branch + "_" + "base-" + from
	return &BuildInfo {
		Project: proj,
		Env: env,
		Tag: tag,
	}
}

func processBuildEvent(req *webhook.Request) {

	for _, stage := range req.Build.Stages {
		project := stage.Name
		info := GetTagFromBuildInfo(project, req.Build)

		switch req.Action {
		case "created":
			Rolling.CreateImage(info.Project, info.Tag)
		case "updated":
			if req.Build.Status == "success" {
				Rolling.UpdateImage(info.Project, info.Tag, info.Env, "Deployable", req.Build.Error)
			} else if req.Build.Status == "failure" {
				Rolling.UpdateImage(info.Project, info.Tag, info.Env, "Compile Failed", req.Build.Error)
			}
		}
	}

	return
}

func processUserEvent(req *webhook.Request) {
	// TODO
	return
}

func processRepoEvent(req *webhook.Request) {
	// TODO
	return
}