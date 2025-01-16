package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger/ha-addon-joplin-server/internal/dagger"

	cplatforms "github.com/containerd/platforms"
)

// Build container images
func (m *HaAddonJoplinServer) Build(
	// +optional
	version string,
) (*HaAddonJoplinServer, error) {
	specifiers := []string{
		"linux/amd64",
		"linux/arm64/v8",
	}
	platforms, err := cplatforms.ParseAll(specifiers)
	if err != nil {
		return nil, err
	}

	// Addon
	git := dag.Gitlocal(m.Worktree)

	gitCommit, err := git.GetLatestCommit(context.TODO())
	if err != nil {
		return nil, err
	}

	gitTag, err := git.GetLatestTag(context.TODO())
	if err != nil {
		return nil, err
	}

	gitUncommited, err := git.Uncommited(context.TODO())
	if err != nil {
		return nil, err
	}

	addonVersion := getVersion(version, gitTag, gitCommit, gitUncommited)

	// App
	appRepo := dag.Git(appSourceUrl).Tag(joplinVersion)

	appGitCommit, err := appRepo.Commit(context.TODO())
	if err != nil {
		return nil, err
	}

	appRepoWorktree := appRepo.Tree()

	// Build
	tsNow := time.Now()

	appBuildArgs := []string{
		fmt.Sprintf("BUILD_DATE=%s", tsNow.Format("2006-01-02T15:04:05 -0700")),
		fmt.Sprintf("VERSION=%s", joplinVersion),
		fmt.Sprintf("REVISION=%s", appGitCommit),
	}

	appImageDir := dag.Docker().Build(
		appRepoWorktree,
		dagger.Platform(cplatforms.Format(platforms[0])),
		dagger.DockerBuildOpts{Dockerfile: "Dockerfile.server", BuildArgs: appBuildArgs},
	).Directory("/home/joplin/packages")

	for _, platform := range platforms {
		addonBuildArgs := []string{
			fmt.Sprintf("APP_VERSION=%s", joplinVersion),
			fmt.Sprintf("APP_REVISION=%s", appGitCommit),
			fmt.Sprintf("BUILD_ARCH=%s", platform.Architecture),
			fmt.Sprintf("BUILD_DATE=%s", tsNow.Format("2006-01-02T15:04:05 -0700")),
			fmt.Sprintf("BUILD_VERSION=%s", addonVersion),
			fmt.Sprintf("BASHIO_VERSION=%s", bashioVersion),
		}

		addonImage := dag.Docker().Build(
			m.Worktree,
			dagger.Platform(cplatforms.Format(platform)),
			dagger.DockerBuildOpts{BuildArgs: addonBuildArgs},
		).WithDirectory("/home/joplin/packages", appImageDir, dagger.ContainerWithDirectoryOpts{Owner: "joplin:joplin"})

		m.Containers = append(m.Containers, addonImage)
	}

	return m, nil
}

func (m *HaAddonJoplinServer) Stdout(ctx context.Context) (string, error) {
	var outputs []string

	for _, ctr := range m.Containers {
		out, err := ctr.WithExec([]string{"cat", "/etc/os-release"}).Stdout(ctx)
		if err != nil {
			return "", err
		}

		outputs = append(outputs, out)
	}

	return strings.Join(outputs, "\n"), nil
}
