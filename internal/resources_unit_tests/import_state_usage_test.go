package resources_unit_tests

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResourcesUseImportStateHelper(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "resources")

	var offenders []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, readErr := os.ReadFile(path) // #nosec G304 -- test scans repo files under known root
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), "ImportStatePassthroughID(") {
			offenders = append(offenders, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("failed to scan resource files: %v", err)
	}

	if len(offenders) > 0 {
		t.Fatalf("resources still use ImportStatePassthroughID: %s", strings.Join(offenders, ", "))
	}
}
