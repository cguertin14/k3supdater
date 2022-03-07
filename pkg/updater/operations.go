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

const (
	k3sVersionKey string = "k3s_release_version"
)

func (c *ClientSet) UpdateK3sRelease(ctx context.Context, req UpdateReleaseReq) error {
	logger := logger.NewFromContextOrDefault(ctx)

	contents, _, _, err := c.client.GetRepositoryContents(ctx, legacy.GetRepositoryContentsRequest{
		Owner:  req.Repo.Owner,
		Repo:   req.Repo.Name,
		Path:   req.Repo.Path,
		Branch: req.Repo.Branch,
	})
	if err != nil {
		return fmt.Errorf("error when fetching repo: %s", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(*contents.Content)
	if err != nil {
		return fmt.Errorf("error when decoding %q: %s", req.Repo.Path, err)
	}

	groupVarsFileContent := string(decoded)
	regz := regexp.MustCompile(fmt.Sprintf("%s.*", k3sVersionKey))
	extracted := regz.FindStringSubmatch(groupVarsFileContent)
	if len(extracted) == 0 {
		return fmt.Errorf("error when extracting k3s version: %q key not found in %q", k3sVersionKey, req.Repo.Path)
	}

	currentVersion := extracted[0][len(k3sVersionKey)+2:]
	releases, _, err := c.client.GetRepositoryReleases(ctx, legacy.CommonRequest{
		Owner: req.ReleaseRepo.Owner,
		Repo:  req.ReleaseRepo.Name,
	})
	if err != nil {
		return fmt.Errorf(
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
		if strings.Index(*r.Name, "rc") != -1 {
			continue
		}
		if len(latestStableVersions) < 3 {
			latestStableVersions = append(latestStableVersions, r)
		}
	}

	versionToUpdateTo := &github.RepositoryRelease{}
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
			versionToUpdateTo = v
			break
		}
	}

	// No update required in this case
	if versionToUpdateTo.Name == nil {
		logger.Infof("Current version %q is the latest version available for k3s, therefore not updating", currentVersion)
		return nil
	}

	// Proceed to make the update
	//
	// Step 1: Create a new branch
	// Step 2: Update file content locally
	// Step 3: Update file content on github repo, on a new branch
	// Step 4: Open pull request with new release details.

	branch, _, err := c.client.GetBranch(ctx, legacy.GetBranchRequest{
		Owner:      req.Repo.Owner,
		Repo:       req.Repo.Name,
		BranchName: fmt.Sprintf("refs/heads/%s", req.Repo.Branch),
	})
	if err != nil {
		return fmt.Errorf("error when fetching branch %q: %s", req.Repo.Branch, err)
	}

	branchName := fmt.Sprintf("minor/k3s-%s-update", *versionToUpdateTo.Name)
	branch, _, err = c.client.CreateBranch(ctx, legacy.CreateBranchRequest{
		Owner: req.Repo.Owner,
		Repo:  req.Repo.Name,
		Reference: &github.Reference{
			Ref:    github.String(fmt.Sprintf("refs/heads/%s", branchName)),
			Object: branch.Object,
		},
	})
	if err != nil {
		return fmt.Errorf("error when creating branch %q: %s", branchName, err)
	}

	now := time.Now()
	newGroupVarsFileContent := strings.ReplaceAll(
		groupVarsFileContent,
		fmt.Sprintf("%s: %s", k3sVersionKey, currentVersion),
		fmt.Sprintf("%s: %s", k3sVersionKey, *versionToUpdateTo.Name),
	)

	_, _, err = c.client.UpdateFile(ctx, legacy.UpdateFileRequest{
		Owner:    req.Repo.Owner,
		Repo:     req.Repo.Name,
		FilePath: req.Repo.Path,
		RepositoryContentFileOptions: &github.RepositoryContentFileOptions{
			Content: []byte(newGroupVarsFileContent),
			Branch:  github.String(branchName),
			Committer: &github.CommitAuthor{
				Name:  github.String("k3supdater-bot"),
				Email: github.String("k3supdater-bot@k3s.io"),
				Date:  &now,
			},
			Message: github.String(
				fmt.Sprintf("Updated k3s version %q to %q.", currentVersion, *versionToUpdateTo.Name),
			),
			SHA: contents.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("error when updating file %q: %s", req.Repo.Path, err)
	}

	_, _, err = c.client.CreatePullRequest(ctx, legacy.CreatePRRequest{
		Owner: req.Repo.Owner,
		Repo:  req.Repo.Name,
		NewPullRequest: &github.NewPullRequest{
			Base: github.String(req.Repo.Branch),
			Head: github.String(branchName),
			Body: versionToUpdateTo.Body,
			Title: github.String(
				fmt.Sprintf(
					"Minor: k3s update from %s to %s",
					currentVersion,
					*versionToUpdateTo.Name,
				),
			),
		},
	})
	if err != nil {
		return fmt.Errorf("error when opening pull request on repository: %s", err)
	}

	return nil
}
