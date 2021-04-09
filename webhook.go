package main

import (
	"github.com/drone/drone-go/plugin/webhook"
)

func processBuildEvent(req *webhook.Request) {

	droneInfo := ProcessRepoAndEventInfo(req.Repo, req.Build)

	switch req.Action {
	case "created":
		Rolling.CreateImage(droneInfo.Project, droneInfo.Tag)
	case "updated":
		if req.Build.Status == "success" {
            Rolling.UpdateImage(droneInfo.Project, droneInfo.Tag, droneInfo.Env, "Deployable", req.Build.Error)
        } else if req.Build.Status == "failure" {
            Rolling.UpdateImage(droneInfo.Project, droneInfo.Tag, droneInfo.Env, "Compile Failed", req.Build.Error)
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
	switch req.Action {
	case "activated":
		NewRepository(req.Repo, req.SSHURL)
	case "deactivated":
		DisableRepository(req.Repo)
	}
	return
}