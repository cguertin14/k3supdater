package updater

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	github_mocks "github.com/cguertin14/k3supdater/pkg/github/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v43/github"
)

func TestUpdateK3sRelease(t *testing.T) {
	// handle working + failing cases for each
	// function call, that's it. edge cases
	// will be handled individually in each
	// function test.

	cases := map[string]struct {
		currentVersion            string
		groupVarsFileContentError error
		latestReleaseError        error
		createNewBranchError      error
		updateFileError           error
		createPRError             error

		expectError bool
	}{
		"success case with no error": {},
		"success case with no new version": {
			currentVersion: "v1.25.3",
		},
		"error case with groupVarsFile error": {
			groupVarsFileContentError: errors.New("some error"),
			expectError:               true,
		},
		"error case with latest release error": {
			latestReleaseError: errors.New("some error"),
			expectError:        true,
		},
		"error case with create new branch error": {
			createNewBranchError: errors.New("some error"),
			expectError:          true,
		},
		"error case with update file error": {
			updateFileError: errors.New("some error"),
			expectError:     true,
		},
		"error case with create PR error": {
			createPRError: errors.New("some error"),
			expectError:   true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.currentVersion == "" {
				c.currentVersion = "v1.23.3"
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().GetRepositoryContents(gomock.Any(), gomock.Any()).
				Times(1).
				Return(&github.RepositoryContent{
					SHA: github.String("some sha"),
					Content: github.String(
						base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s: %s", k3sVersionKey, c.currentVersion))),
					),
				}, nil, nil, c.groupVarsFileContentError)
			githubMockClient.EXPECT().GetRepositoryReleases(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return([]*github.RepositoryRelease{
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.23.5-rc1"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.23.4"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.22.3"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.21.5"),
					},
				}, nil, c.latestReleaseError)
			githubMockClient.EXPECT().GetBranch(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(&github.Reference{
					Object: &github.GitObject{},
				}, nil, nil)
			githubMockClient.EXPECT().CreateBranch(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(nil, nil, c.createNewBranchError)
			githubMockClient.EXPECT().UpdateFile(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(nil, nil, c.updateFileError)
			githubMockClient.EXPECT().CreatePullRequest(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(nil, nil, c.createPRError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			err := client.UpdateK3sRelease(context.Background(), UpdateReleaseReq{
				Repo: Repository{
					Owner:  "some owner",
					Name:   "some name",
					Path:   "/some/existing/path",
					Branch: "main",
				},
				ReleaseRepo: Repository{
					Owner: "k3s-io",
					Name:  "k3s",
				},
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}

func TestGetGroupVarsFileContent(t *testing.T) {
	cases := map[string]struct {
		fileContent string

		expectedError error
		expectError   bool
	}{
		"success case with no error": {
			fileContent: base64.StdEncoding.EncodeToString([]byte("some text")),
		},
		"error case with response error": {
			expectedError: errors.New("some error"),
			expectError:   true,
		},
		"error case with base64 error": {
			fileContent: "some unencoded text",
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().GetRepositoryContents(gomock.Any(), gomock.Any()).
				Times(1).
				Return(&github.RepositoryContent{
					SHA:     github.String("some sha"),
					Content: github.String(c.fileContent),
				}, nil, nil, c.expectedError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			_, _, err := client.getGroupVarsFileContent(context.Background(), UpdateReleaseReq{
				Repo: Repository{
					Owner:  "some owner",
					Name:   "some name",
					Path:   "/some/existing/path",
					Branch: "main",
				},
				ReleaseRepo: Repository{
					Owner: "k3s-io",
					Name:  "k3s",
				},
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}

func TestGetLatestK3sRelease(t *testing.T) {
	cases := map[string]struct {
		fileContent string

		expectedError error
		expectError   bool
	}{
		"success case with no error": {
			fileContent: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s: %s", k3sVersionKey, "v1.23.3"))),
		},
		"error case with response error": {
			expectedError: errors.New("some error"),
			expectError:   true,
		},
		"error case with base64 error": {
			fileContent: "some unencoded text",
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().GetRepositoryReleases(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return([]*github.RepositoryRelease{
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.23.5-rc1"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.23.4"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.22.3"),
					},
					{
						Body: github.String("some release notes"),
						Name: github.String("v1.21.5"),
					},
				}, nil, c.expectedError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			_, _, err := client.getLatestK3sRelease(context.Background(), getLatestK3sReleaseRequest{
				UpdateReleaseReq: UpdateReleaseReq{
					Repo: Repository{
						Owner:  "some owner",
						Name:   "some name",
						Path:   "/some/existing/path",
						Branch: "main",
					},
					ReleaseRepo: Repository{
						Owner: "k3s-io",
						Name:  "k3s",
					},
				},
				fileContent: c.fileContent,
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}

func TestCreateNewBranch(t *testing.T) {
	cases := map[string]struct {
		getBranchError    error
		createBranchError error
		expectError       bool
	}{
		"success case with no error": {},
		"error case with get branch error": {
			getBranchError: errors.New("some error"),
			expectError:    true,
		},
		"error case with create branch error": {
			createBranchError: errors.New("some error"),
			expectError:       true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().GetBranch(gomock.Any(), gomock.Any()).
				Times(1).
				Return(&github.Reference{
					Object: &github.GitObject{},
				}, nil, c.getBranchError)
			githubMockClient.EXPECT().CreateBranch(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(nil, nil, c.createBranchError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			_, err := client.createNewBranch(context.Background(), createNewBranchReq{
				latestRelease: &github.RepositoryRelease{
					Name: github.String("some branch name"),
				},
				UpdateReleaseReq: UpdateReleaseReq{
					Repo: Repository{
						Owner:  "some owner",
						Name:   "some name",
						Path:   "/some/existing/path",
						Branch: "main",
					},
					ReleaseRepo: Repository{
						Owner: "k3s-io",
						Name:  "k3s",
					},
				},
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}

func TestUpdateFile(t *testing.T) {
	cases := map[string]struct {
		updateFileError error
		expectError     bool
	}{
		"success case with no error": {},
		"error case with update file error": {
			updateFileError: errors.New("some error"),
			expectError:     true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().UpdateFile(gomock.Any(), gomock.Any()).
				MaxTimes(1).
				Return(nil, nil, c.updateFileError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			err := client.updateFile(context.Background(), updateFileReq{
				latestRelease: &github.RepositoryRelease{
					Name: github.String("some branch name"),
				},
				fileContent:    "some file content",
				currentVersion: "v1.23.4",
				branchName:     "main",
				repoContent: &github.RepositoryContent{
					SHA: github.String("some sha"),
				},
				UpdateReleaseReq: UpdateReleaseReq{
					Repo: Repository{
						Owner:  "some owner",
						Name:   "some name",
						Path:   "/some/existing/path",
						Branch: "main",
					},
					ReleaseRepo: Repository{
						Owner: "k3s-io",
						Name:  "k3s",
					},
				},
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}

func TestCreatePR(t *testing.T) {
	cases := map[string]struct {
		createPRError error
		expectError   bool
	}{
		"success case with no error": {},
		"error case with update file error": {
			createPRError: errors.New("some error"),
			expectError:   true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create new mock client instance
			githubMockClient := github_mocks.NewMockClient(ctrl)

			// define mock behavior
			githubMockClient.EXPECT().CreatePullRequest(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, nil, c.createPRError)

			// create mock updater client
			client := NewClient(context.Background(), Dependencies{
				Client: githubMockClient,
			})

			err := client.createPR(context.Background(), createPRRequest{
				latestRelease: &github.RepositoryRelease{
					Name: github.String("some branch name"),
				},
				currentVersion: "v1.23.4",
				branchName:     "main",
				UpdateReleaseReq: UpdateReleaseReq{
					Repo: Repository{
						Owner:  "some owner",
						Name:   "some name",
						Path:   "/some/existing/path",
						Branch: "main",
					},
					ReleaseRepo: Repository{
						Owner: "k3s-io",
						Name:  "k3s",
					},
				},
			})
			if c.expectError && err == nil {
				t.FailNow()
			}
		})
	}
}
