package main

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/sirupsen/logrus"
)

type DroneServer struct {
	Endpoint 			string
	YamlPluginSecret 		string
	WebhookPluginSecret		string
}

type YamlPlugin struct {

}

func NewYamlPlugin(idx int) http.Handler {

	logrus.SetLevel(logrus.DebugLevel)

	secret := DroneServers[idx].YamlPluginSecret
	handler := config.Handler(
		&YamlPlugin{},
		secret,
		logrus.StandardLogger(),
	)

	return handler
}

func (p *YamlPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {

	logrus.Info("New coming request")
	logrus.Info("Repo Info", req.Repo)
	logrus.Info("Build Info", req.Build)

	m,err := NewManifest(&req.Repo, &req.Build)
	if err!=nil {
		return nil,err
	}

	data, err := m.Compile()
	if err!=nil {
		return nil, err
	}

	return &drone.Config {
		Data: data,
	}, nil
}

type RegistryPlugin struct {

}

func NewRegistryPlugin() http.Handler {
	logrus.SetLevel(logrus.DebugLevel)

	handler := registry.Handler(
		YamlPluginSecret,
		&RegistryPlugin{},
		logrus.StandardLogger(),
	)

	return handler
}

func (p *RegistryPlugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {
	credentials := []*drone.Registry{
		{
			Address:  HarborBaseUrl,
			Username: HarborUser,
			Password: HarborSecret,
		},
	}

	return credentials, nil
}

type WebhookPlugin struct {
	Id 			int
}

func NewWebhookPlugin(idx int) http.Handler {

	logrus.SetLevel(logrus.DebugLevel)

	secret := DroneServers[idx].WebhookPluginSecret
	handler := webhook.Handler(
		&WebhookPlugin{},
		secret,
		logrus.StandardLogger(),
	)

	return handler
}

func (p *WebhookPlugin) Deliver(ctx context.Context, req *webhook.Request) error {
	switch req.Event {
		case "build":
			go processBuildEvent(req)
		case "user":
			go processUserEvent(req)
		case "repo":
			go processRepoEvent(req)
		default:
	}

	return nil
}