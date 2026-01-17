package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func findTFLintConfigPath(startDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ".tflint.hcl")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("config file not found searching parents of %s", startDir)
		}
		dir = parent
	}
}

func shouldSkipDir(path string) bool {
	base := filepath.Base(path)
	if base == ".terraform" || base == ".git" || base == "vendor" || base == "node_modules" {
		return true
	}
	return false
}

func hasTerraformFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tf") {
			return true, nil
		}
	}
	return false, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	configPath, err := findTFLintConfigPath(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding TFLint config: %v\n", err)
		os.Exit(1)
	}
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving config path: %v\n", err)
		os.Exit(1)
	}

	// Initialize tflint plugins once
	fmt.Println("Initializing TFLint plugins...")

	//nolint:gosec // trusted config path
	initCmd := exec.Command("tflint", "--init", "--config", configPath)
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	if err := initCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing tflint: %v\n", err)
		os.Exit(1)
	}
	var dirsToLint []string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDir(path) {
				return filepath.SkipDir
			}
			hasTf, err := hasTerraformFiles(path)
			if err != nil {
				return err
			}
			if hasTf {
				dirsToLint = append(dirsToLint, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directories: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Linting %d directories...\n", len(dirsToLint))

	// Worker pool
	numWorkers := runtime.NumCPU()
	if numWorkers > 4 {
		numWorkers = 4
	}

	jobs := make(chan string, len(dirsToLint))
	results := make(chan error, len(dirsToLint))
	var wg sync.WaitGroup
	var outputMutex sync.Mutex
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dir := range jobs {
				args := []string{"--config", configPath, "--chdir", dir}
				args = append(args, os.Args[1:]...)
				// #nosec G204
				cmd := exec.Command("tflint", args...)
				output, err := cmd.CombinedOutput()
				if err != nil {
					outputMutex.Lock()
					fmt.Printf("Failure in %s:\n%s\n", dir, string(output))
					outputMutex.Unlock()
					results <- err
				} else {
					results <- nil
				}
			}
		}()
	}

	for _, dir := range dirsToLint {
		jobs <- dir
	}
	close(jobs)
	wg.Wait()
	close(results)
	failed := false
	for err := range results {
		if err != nil {
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
