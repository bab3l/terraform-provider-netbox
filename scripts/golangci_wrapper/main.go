package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const gotmpDirName = ".gotmp"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: golangci_wrapper <golangci-lint args...>")
		os.Exit(2)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	repoRoot := findRepoRoot(cwd)

	gotmpDirAbs, err := filepath.Abs(filepath.Join(repoRoot, gotmpDirName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve GOTMPDIR path: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(gotmpDirAbs, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create GOTMPDIR %q: %v\n", gotmpDirAbs, err)
		os.Exit(1)
	}

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Fprintln(os.Stderr, "golangci-lint not found on PATH")
		fmt.Fprintln(os.Stderr, "install: https://golangci-lint.run/")
		os.Exit(127)
	}

	args := os.Args[1:]

	cmd := exec.Command("golangci-lint", args...) // #nosec G204 -- fixed binary name; args are from developer-invoked tooling
	cmd.Env = withEnv(os.Environ(), "GOTMPDIR="+gotmpDirAbs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitCode(exitErr))
		}

		fmt.Fprintf(os.Stderr, "failed to run golangci-lint: %v\n", err)
		os.Exit(1)
	}
}

func findRepoRoot(startDir string) string {
	dir := startDir
	for {
		gitDir := filepath.Join(dir, ".git")
		if fi, err := os.Stat(gitDir); err == nil && fi.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return startDir
		}
		dir = parent
	}
}

func withEnv(base []string, kv string) []string {
	key := strings.SplitN(kv, "=", 2)[0]
	out := make([]string, 0, len(base)+1)
	for _, e := range base {
		if strings.HasPrefix(e, key+"=") {
			continue
		}
		out = append(out, e)
	}
	out = append(out, kv)
	return out
}

func exitCode(exitErr *exec.ExitError) int {
	code := exitErr.ExitCode()
	if code != -1 {
		return code
	}

	if runtime.GOOS == "windows" {
		return 1
	}

	return 1
}
