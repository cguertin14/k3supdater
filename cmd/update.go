package cmd

import (
	"context"
	"fmt"

	"github.com/cguertin14/k3supdater/pkg/updater"
	"github.com/cguertin14/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	githubAccessToken string = "GITHUB_ACCESS_TOKEN"
	repoOwner         string = "repo-owner"
	repoName          string = "repo-name"
	repoBranch        string = "repo-branch"
	groupVarsFilepath string = "group-vars-filepath"
	releaseRepoOwner  string = "release-repo-owner"
	releaseRepoName   string = "release-repo-name"
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

	// init logger
	ctxLogger := logger.Initialize(logger.Config{Level: "info"})
	ctx = context.WithValue(ctx, logger.CtxKey, ctxLogger)

	// create business logic client here
	client := updater.NewClient(ctx, updater.Dependencies{
		AccessToken: v.GetString(githubAccessToken),
	})

	if err = client.UpdateK3sRelease(ctx, updater.UpdateReleaseReq{
		Repo: updater.Repository{
			Owner:  v.GetString(repoOwner),
			Name:   v.GetString(repoName),
			Path:   v.GetString(groupVarsFilepath),
			Branch: v.GetString(repoBranch),
		},
		ReleaseRepo: updater.Repository{
			Owner: v.GetString(releaseRepoOwner),
			Name:  v.GetString(releaseRepoName),
		},
	}); err != nil {
		return fmt.Errorf("error when updating k3s version: %s", err)
	}

	return
}

func init() {
	updateCmd.Flags().String(repoOwner, "", "The github owner of the repository (i.e.: cguertin14, some-other-user, etc.)")
	updateCmd.Flags().String(repoName, "", "The github repository name minus the user/org part (i.e.: k3s-ansible-ha, some-other-repo, etc.)")
	updateCmd.Flags().String(repoBranch, "main", "The branch of your github repo to edit (i.e.: main)")
	updateCmd.Flags().String(groupVarsFilepath, "inventory/pi-cluster/group_vars/all.yml", "The path of the 'inventory/<YOUR_MACHINE>/group_vars/<YOUR_FILE>.yml' file in your github repo to edit.")
	updateCmd.Flags().String(releaseRepoOwner, "k3s-io", "The github owner of the release repository (i.e.: k3s-io, some-other-org, etc.).")
	updateCmd.Flags().String(releaseRepoName, "k3s", "The github release repository name minus the user/org part (i.e.: k3s, some-other-repo, etc.)")

	// Required flags
	updateCmd.MarkFlagRequired(repoOwner)
	updateCmd.MarkFlagRequired(repoName)
}
