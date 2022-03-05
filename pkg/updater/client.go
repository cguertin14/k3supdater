package updater

import (
	"context"

	github "github.com/cguertin14/k3s-ansible-updater/pkg/github"
)

type ClientSet struct {
	client         github.Client
	repoURI        string
	releaseRepoURI string
}

type Dependencies struct {
	Client         github.Client
	AccessToken    string
	RepoURI        string
	ReleaseRepoURI string
}

func NewClient(ctx context.Context, deps Dependencies) *ClientSet {
	c := &ClientSet{
		repoURI:        deps.RepoURI,
		releaseRepoURI: deps.ReleaseRepoURI,
	}

	if deps.Client == nil {
		c.client = github.NewClient(ctx, deps.AccessToken)
	}

	return c
}
