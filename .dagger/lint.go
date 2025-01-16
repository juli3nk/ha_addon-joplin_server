package main

import (
	"context"

	"dagger/ha-addon-joplin-server/internal/dagger"
)

// Job: commit-msg
func (m *HaAddonJoplinServer) LintCommitMsg(ctx context.Context) (string, error) {
	return dag.Commitlint().Lint(m.Worktree, dagger.CommitlintLintOpts{Args: []string{"-l"}}).Stdout(ctx)
}

// Job: jsonfile
func (m *HaAddonJoplinServer) LintJsonFile(ctx context.Context) (string, error) {
	return dag.Jsonfile().Lint(ctx, m.Worktree)
}

// Job: dockerfile
func (m *HaAddonJoplinServer) LintDockerfile(ctx context.Context) (string, error) {
	opts := dagger.DockerLintOpts{
		Ignore: []string{
			"DL3008",
			"DL4006",
		},
	}

	return dag.Docker().Lint(ctx, m.Worktree, opts)
}
