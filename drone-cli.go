package main

import (
	"github.com/drone/drone-go/drone"

	"golang.org/x/oauth2"
)

type DroneCli struct {
	Owner 			string
	Endpoint		string
	Token 			string
	Client 			drone.Client
}

func NewDroneClient(owner string, endpoint string, token string) *DroneCli {
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	client := drone.NewClient(endpoint, auther)

	return &DroneCli {
		Owner: owner,
		Endpoint: endpoint,
		Token: token,
		Client: client,
	}
}

func (this *DroneCli) GetRepositoryList() []string {
	repo, err := this.Client.RepoList()
	if err!=nil {
		return nil
	}

	repos := make([]string, 0)
	for _,r :=range repo {
		repos = append(repos, r.Name)
	}

	return repos
}

func (this *DroneCli) EnableRepo(group string, name string) error{
	repo, err := this.Client.RepoEnable(group, name)
	if err!=nil {
		return err
	}

	return nil
}

func (this *DroneCli) DisableRepo(group string, name string) error{
	repo, err := this.Client.RepoDisable(group, name)
	if err!=nil {
		return err
	}

	return nil
}