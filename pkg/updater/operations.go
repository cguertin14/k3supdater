package updater

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	legacy "github.com/cguertin14/k3supdater/pkg/github"
	"github.com/cguertin14/logger"
	"github.com/google/go-github/v43/github"
	"golang.org/x/mod/semver"
)

type Repository struct {
	Owner  string
	Name   string
	Path   string
	Branch string
}

type UpdateReleaseReq struct {
	Repo        Repository
	ReleaseRepo Repository
}

type getLatestK3sReleaseRequest struct {
	UpdateReleaseReq
	fileContent string
}

type createNewBranchReq struct {
	UpdateReleaseReq
	latestRelease *github.RepositoryRelease
}

type updateFileReq struct {
	fileContent    string
	currentVersion string
	latestRelease  *github.RepositoryRelease
	repoContent    *github.RepositoryContent
	branchName     string
	UpdateReleaseReq
}

type createPRRequest struct {
	currentVersion string
	latestRelease  *github.RepositoryRelease
	branchName     string
	UpdateReleaseReq
}

const (
	k3sVersionKey string = "k3s_release_version"
)

func (c *ClientSet) getGroupVarsFileContent(ctx context.Context, req UpdateReleaseReq) (repoContent *github.RepositoryContent, fileContent string, err error) {
	logger := logger.NewFromContextOrDefault(ctx)
	logger.Infof("Fetching %q from %s/%s...\n", req.Repo.Path, req.Repo.Owner, req.Repo.Name)

	repoContent, _, _, err = c.client.GetRepositoryContents(ctx, legacy.GetRepositoryContentsRequest{
		Owner:  req.Repo.Owner,
		Repo:   req.Repo.Name,
		Path:   req.Repo.Path,
		Branch: req.Repo.Branch,
	})
	if err != nil {
		err = fmt.Errorf("error when fetching repo: %s", err)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(*repoContent.Content)
	if err != nil {
		err = fmt.Errorf("error when decoding %q: %s", req.Repo.Path, err)
		return
	}

	fileContent = string(decoded)
	return
}

func (c *ClientSet) getLatestK3sRelease(ctx context.Context, req getLatestK3sReleaseRequest) (latestRelease *github.RepositoryRelease, currentVersion string, err error) {
	logger := logger.NewFromContextOrDefault(ctx)
	logger.Infof("Fetching the latest k3s release from %s/%s...\n", req.ReleaseRepo.Owner, req.ReleaseRepo.Name)

	regz := regexp.MustCompile(fmt.Sprintf("%s.*", k3sVersionKey))
	extracted := regz.FindStringSubmatch(req.fileContent)
	if len(extracted) == 0 {
		return nil, "", fmt.Errorf("error when extracting k3s version: %q key not found in %q", k3sVersionKey, req.Repo.Path)
	}

	currentVersion = extracted[0][len(k3sVersionKey)+2:]
	releases, _, err := c.client.GetRepositoryReleases(ctx, legacy.CommonRequest{
		Owner: req.ReleaseRepo.Owner,
		Repo:  req.ReleaseRepo.Name,
	})
	if err != nil {
		return nil, "", fmt.Errorf(
			"error when fetching releases from %s/%s: %s",
			req.ReleaseRepo.Owner,
			req.ReleaseRepo.Name,
			err,
		)
	}

	// Find latest stable versions among all releases
	latestStableVersions := make([]*github.RepositoryRelease, 0)
	for _, r := range releases {
		// Exclude release candidates since
		// they're not stable versions
		if strings.Contains(*r.Name, "rc") {
			continue
		}
		if len(latestStableVersions) < 3 {
			latestStableVersions = append(latestStableVersions, r)
		}
	}

	latestRelease = &github.RepositoryRelease{}
	for _, v := range latestStableVersions {
		compared := semver.Compare(currentVersion, *v.Name)
		if compared == -1 {
			// if compared == -1, this means we need to
			// update to that specific version.
			//
			// if compared == 0, this means we're at
			// the latest version available.
			//
			// if compared == 1, this means the current
			// version is more recent than the other ones.
			latestRelease = v
			logger.Warnf("A new k3s version is available: %q\n", *latestRelease.Name)
			break
		}
	}

	return
}

func (c *ClientSet) createNewBranch(ctx context.Context, req createNewBranchReq) (branchName string, err error) {
	branch, _, err := c.client.GetBranch(ctx, legacy.GetBranchRequest{
		Owner:      req.Repo.Owner,
		Repo:       req.Repo.Name,
		BranchName: fmt.Sprintf("refs/heads/%s", req.Repo.Branch),
	})
	if err != nil {
		return "", fmt.Errorf("error when fetching branch %q: %s", req.Repo.Branch, err)
	}

	branchName = fmt.Sprintf("minor/k3s-%s-update", *req.latestRelease.Name)
	_, _, err = c.client.CreateBranch(ctx, legacy.CreateBranchRequest{
		Owner: req.Repo.Owner,
		Repo:  req.Repo.Name,
		Reference: &github.Reference{
			Ref:    github.String(fmt.Sprintf("refs/heads/%s", branchName)),
			Object: branch.Object,
		},
	})
	if err != nil {
		return "", fmt.Errorf("error when creating branch %q: %s", branchName, err)
	}

	return
}

func (c *ClientSet) updateFile(ctx context.Context, req updateFileReq) (err error) {
	now := time.Now()
	newGroupVarsFileContent := strings.ReplaceAll(
		req.fileContent,
		fmt.Sprintf("%s: %s", k3sVersionKey, req.currentVersion),
		fmt.Sprintf("%s: %s", k3sVersionKey, *req.latestRelease.Name),
	)

	_, _, err = c.client.UpdateFile(ctx, legacy.UpdateFileRequest{
		Owner:    req.Repo.Owner,
		Repo:     req.Repo.Name,
		FilePath: req.Repo.Path,
		RepositoryContentFileOptions: &github.RepositoryContentFileOptions{
			Content: []byte(newGroupVarsFileContent),
			Branch:  github.String(req.branchName),
			Committer: &github.CommitAuthor{
				Name:  github.String("k3supdater-bot"),
				Email: github.String("k3supdater-bot@k3s.io"),
				Date:  &now,
			},
			Message: github.String(
				fmt.Sprintf("Updated k3s version %q to %q.", req.currentVersion, *req.latestRelease.Name),
			),
			SHA: req.repoContent.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("error when updating file %q: %s", req.Repo.Path, err)
	}

	return
}

func (c *ClientSet) createPR(ctx context.Context, req createPRRequest) error {
	_, _, err := c.client.CreatePullRequest(ctx, legacy.CreatePRRequest{
		Owner: req.Repo.Owner,
		Repo:  req.Repo.Name,
		NewPullRequest: &github.NewPullRequest{
			Base: github.String(req.Repo.Branch),
			Head: github.String(req.branchName),
			Body: req.latestRelease.Body,
			Title: github.String(
				fmt.Sprintf(
					"Minor: k3s update from %s to %s",
					req.currentVersion,
					*req.latestRelease.Name,
				),
			),
		},
	})
	if err != nil {
		return fmt.Errorf("error when opening pull request on repository: %s", err)
	}

	return nil
}

func (c *ClientSet) UpdateK3sRelease(ctx context.Context, req UpdateReleaseReq) (err error) {
	logger := logger.NewFromContextOrDefault(ctx)
	repoContent, fileContent, err := c.getGroupVarsFileContent(ctx, req)
	if err != nil {
		return
	}

	latestRelease, currentVersion, err := c.getLatestK3sRelease(ctx, getLatestK3sReleaseRequest{
		UpdateReleaseReq: req,
		fileContent:      fileContent,
	})
	if err != nil {
		return
	}

	// No update required in this case
	if latestRelease.Name == nil {
		logger.Infof("Current version %q is the latest version available for k3s, therefore not updating.\n", currentVersion)
		return nil
	}

	// Proceed to make the update
	//
	// Step 1: Create a new branch
	// Step 2: Update file content locally
	// Step 3: Update file content on github repo, on a new branch
	// Step 4: Open pull request with new release details.
	branchName, err := c.createNewBranch(ctx, createNewBranchReq{
		UpdateReleaseReq: req,
		latestRelease:    latestRelease,
	})
	if err != nil {
		return
	}

	// Section 6: update file's content
	if err = c.updateFile(ctx, updateFileReq{
		UpdateReleaseReq: req,
		fileContent:      fileContent,
		currentVersion:   currentVersion,
		latestRelease:    latestRelease,
		repoContent:      repoContent,
		branchName:       branchName,
	}); err != nil {
		return
	}

	// Section 7: create PR
	if err = c.createPR(ctx, createPRRequest{
		UpdateReleaseReq: req,
		currentVersion:   currentVersion,
		latestRelease:    latestRelease,
		branchName:       branchName,
	}); err != nil {
		return
	}

	return nil
}
