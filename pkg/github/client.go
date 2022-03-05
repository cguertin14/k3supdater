package legacy

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

type CommonRequest struct {
	Owner string
	Repo  string
}

type CreatePRRequest struct {
	Owner string
	Repo  string
	*github.NewPullRequest
}

type UpdateFileRequest struct {
	Owner    string
	Repo     string
	FilePath string
	*github.RepositoryContentFileOptions
}

type Client interface {
	// GetRepository
	//
	// Fetches a repository given github options.
	GetRepository(ctx context.Context, req CommonRequest) (*github.Repository, *github.Response, error)

	// CreatePullRequest
	//
	// Creates a Pull Requests on a given repository.
	CreatePullRequest(ctx context.Context, req CreatePRRequest) (*github.PullRequest, *github.Response, error)

	// UpdateFile
	//
	// Updates a file in a given repo with new content.
	UpdateFile(ctx context.Context, req UpdateFileRequest) (*github.RepositoryContentResponse, *github.Response, error)
}

type ClientSet struct {
	github *github.Client
}

func NewClient(ctx context.Context, accessToken string) *ClientSet {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &ClientSet{
		github: github.NewClient(tc),
	}
}

// Make sure ClientSet struct
// implements Client interface
var _ Client = &ClientSet{}

func (c *ClientSet) CreatePullRequest(ctx context.Context, req CreatePRRequest) (*github.PullRequest, *github.Response, error) {
	return c.github.PullRequests.Create(
		ctx,
		req.Owner,
		req.Repo,
		req.NewPullRequest,
	)
}

func (c *ClientSet) GetRepository(ctx context.Context, req CommonRequest) (*github.Repository, *github.Response, error) {
	return c.github.Repositories.Get(ctx, req.Owner, req.Repo)
}

func (c *ClientSet) UpdateFile(ctx context.Context, req UpdateFileRequest) (*github.RepositoryContentResponse, *github.Response, error) {
	return c.github.Repositories.UpdateFile(
		ctx,
		req.Owner,
		req.Repo,
		req.FilePath,
		req.RepositoryContentFileOptions,
	)
}
