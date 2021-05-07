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
	tag := timestamp + "_" + version + "_" + branch + from
	return &BuildInfo {
		Project: proj,
		Env: env,
		Tag: tag,
	}
}

func processBuildEvent(req *webhook.Request, ciAddr string) {
	repo := req.Repo

	switch req.Action {
		case "created":
			projects, repoName, err := GetProjectsByRepo(repo.SSHURL)
			if err!=nil {
				fmt.Println("could not resolve git url:", err)
				return
			}

			fmt.Println("pipeline created for ", repoName)
			for _,proj := range projects {
				info := GetTagFromBuildInfo(proj, req.Build)
				jsonStr := []byte(`{"project":"` + proj + `","tag":"` + info.Tag + `","ciAddr":"` + ciAddr + `"}`)
				Rolling.CreateImage(jsonStr)
			}

		case "updated":
			for _,stage :=range req.Build.Stages {
				project := stage.Name
				info := GetTagFromBuildInfo(project, req.Build)

				switch stage.Status {
				case "running":
				case "success":
					Rolling.UpdateImage(info.Project, info.Tag, info.Env, "Deployable", stage.Error)
				case "failure":
					Rolling.UpdateImage(info.Project, info.Tag, info.Env, "Compile Failed", stage.Error)
				default:
					fmt.Println("status: ", stage.Status)
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