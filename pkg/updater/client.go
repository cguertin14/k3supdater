package updater

import (
	"context"

	github "github.com/cguertin14/k3supdater/pkg/github"
)

type ClientSet struct {
	client github.Client
}

type Dependencies struct {
	Client      github.Client
	AccessToken string
}

func NewClient(ctx context.Context, deps Dependencies) *ClientSet {
	c := &ClientSet{
		client: deps.Client,
	}

	if deps.Client == nil {
		c.client = github.NewClient(ctx, deps.AccessToken)
	}

	return c
}
