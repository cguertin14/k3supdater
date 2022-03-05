package cmd

import (
	"fmt"
	"net/url"

	"github.com/cguertin14/k3supdater/pkg/updater"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	repoAccessToken     string = "REPO_ACCESS_TOKEN"
	repoURIFlag         string = "repo-uri"
	repoAccessTokenFlag string = "repo-access-token"
	releaseRepoURIFlag  string = "release-repo-uri"
)

var (
	updateCmd = &cobra.Command{
		Use:           "update",
		Short:         "Update a k3s ansible playbook version to newest one",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          update,
	}
)

func update(cmd *cobra.Command, args []string) (err error) {
	ctx := cmd.Context()

	v := viper.New()
	v.AutomaticEnv()
	if err = v.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("error when parsing flags: %s", err)
	}

	// validate that repo URI flag is a valid URI
	if _, err = url.ParseRequestURI(v.GetString(repoURIFlag)); err != nil {
		return fmt.Errorf("error when parsing URI %q: %s", repoURIFlag, err)
	}

	// validate that release repo URI flag is a valid URI
	if _, err = url.ParseRequestURI(v.GetString(releaseRepoURIFlag)); err != nil {
		return fmt.Errorf("error when parsing URI %q: %s", releaseRepoURIFlag, err)
	}

	// create business logic client here
	client := updater.NewClient(ctx, updater.Dependencies{
		AccessToken:    v.GetString(repoAccessToken),
		RepoURI:        v.GetString(repoURIFlag),
		ReleaseRepoURI: v.GetString(releaseRepoURIFlag),
	})

	if err = client.UpdateK3sRelease(ctx); err != nil {
		return fmt.Errorf("error when updating k3s version: %s", err)
	}

	return
}

func init() {
	updateCmd.Flags().String(repoURIFlag, "", "The repository URI which contains the ansible playbook to update (i.e.: https://github.com/cguertin14/k3s-ansible-ha.git)")
	updateCmd.Flags().String(
		repoAccessTokenFlag,
		"",
		fmt.Sprintf(
			"The github access token with write access to the repository to update (can also be passed from the %q env variable)",
			repoAccessToken,
		),
	)
	updateCmd.Flags().String(releaseRepoURIFlag, "https://github.com/k3s-io/k3s", "The base repository URI which contains release of k3s (defaults to https://github.com/k3s-io/k3s)")

	// Required flags
	updateCmd.MarkFlagRequired(repoURIFlag)
}
