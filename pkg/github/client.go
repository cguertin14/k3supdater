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

type GetRepositoryContentsRequest struct {
	Owner  string
	Repo   string
	Path   string
	Branch string
}

type GetReleaseRequest struct {
	Owner     string
	Repo      string
	ReleaseID int64
}

type GetBranchRequest struct {
	Owner      string
	Repo       string
	BranchName string
}

type CreateBranchRequest struct {
	Owner string
	Repo  string
	*github.Reference
}

type Client interface {
	// GetRepositoryContents
	//
	// Fetches a specific file/folder in a github repo on a given branch.
	GetRepositoryContents(ctx context.Context, req GetRepositoryContentsRequest) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)

	// GetRepositoryReleases
	//
	// Get all releases from a given repository.
	GetRepositoryReleases(ctx context.Context, req CommonRequest) ([]*github.RepositoryRelease, *github.Response, error)

	// GetBranch
	//
	// Returns a branch given a branch name.
	GetBranch(ctx context.Context, req GetBranchRequest) (*github.Reference, *github.Response, error)

	// CreatePullRequest
	//
	// Creates a Pull Requests on a given repository.
	CreatePullRequest(ctx context.Context, req CreatePRRequest) (*github.PullRequest, *github.Response, error)

	// CreateBranch
	//
	// Creates a branch on a given repository.
	CreateBranch(ctx context.Context, req CreateBranchRequest) (*github.Reference, *github.Response, error)

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

func (c *ClientSet) GetRepositoryReleases(ctx context.Context, req CommonRequest) ([]*github.RepositoryRelease, *github.Response, error) {
	return c.github.Repositories.ListReleases(
		ctx,
		req.Owner,
		req.Repo,
		&github.ListOptions{},
	)
}

func (c *ClientSet) GetRepositoryContents(ctx context.Context, req GetRepositoryContentsRequest) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	return c.github.Repositories.GetContents(
		ctx,
		req.Owner,
		req.Repo,
		req.Path,
		&github.RepositoryContentGetOptions{
			Ref: req.Branch,
		},
	)
}

func (c *ClientSet) CreatePullRequest(ctx context.Context, req CreatePRRequest) (*github.PullRequest, *github.Response, error) {
	return c.github.PullRequests.Create(
		ctx,
		req.Owner,
		req.Repo,
		req.NewPullRequest,
	)
}

func (c *ClientSet) GetBranch(ctx context.Context, req GetBranchRequest) (*github.Reference, *github.Response, error) {
	return c.github.Git.GetRef(
		ctx,
		req.Owner,
		req.Repo,
		req.BranchName,
	)
}

func (c *ClientSet) CreateBranch(ctx context.Context, req CreateBranchRequest) (*github.Reference, *github.Response, error) {
	return c.github.Git.CreateRef(
		ctx,
		req.Owner,
		req.Repo,
		req.Reference,
	)
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
