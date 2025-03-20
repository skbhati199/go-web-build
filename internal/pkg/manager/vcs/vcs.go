package vcs

import (
	"context"
	"os/exec"
)

type VersionControl interface {
	Init(ctx context.Context) error
	Add(ctx context.Context, files ...string) error
	Commit(ctx context.Context, message string) error
	Push(ctx context.Context, remote, branch string) error
	Pull(ctx context.Context, remote, branch string) error
}

type Git struct {
	workDir string
}

func NewGit(workDir string) *Git {
	return &Git{workDir: workDir}
}

func (g *Git) Init(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "init")
	cmd.Dir = g.workDir
	return cmd.Run()
}

func (g *Git) Add(ctx context.Context, files ...string) error {
	args := append([]string{"add"}, files...)
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = g.workDir
	return cmd.Run()
}

func (g *Git) Commit(ctx context.Context, message string) error {
	cmd := exec.CommandContext(ctx, "git", "commit", "-m", message)
	cmd.Dir = g.workDir
	return cmd.Run()
}

func (g *Git) Push(ctx context.Context, remote, branch string) error {
	cmd := exec.CommandContext(ctx, "git", "push", remote, branch)
	cmd.Dir = g.workDir
	return cmd.Run()
}

func (g *Git) Pull(ctx context.Context, remote, branch string) error {
	cmd := exec.CommandContext(ctx, "git", "pull", remote, branch)
	cmd.Dir = g.workDir
	return cmd.Run()
}
