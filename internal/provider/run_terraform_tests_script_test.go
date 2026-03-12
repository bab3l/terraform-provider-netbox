package provider

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTerraformIntegrationScriptResourceApiMapStaysAligned(t *testing.T) {
	t.Parallel()

	repoRoot := providerRepoRoot(t)
	scriptPath := filepath.Join(repoRoot, "scripts", "run-terraform-tests.ps1")
	fixtureRoot := filepath.Join(repoRoot, "test", "terraform", "resources")

	mappedResources := loadTerraformIntegrationScriptResourceMap(t, scriptPath)
	providerResources := registeredProviderResourceNames(t)
	fixtureResources := terraformIntegrationFixtureResources(t, fixtureRoot)

	for resourceName := range mappedResources {
		if _, ok := providerResources[resourceName]; !ok {
			t.Errorf("integration script map contains stale resource %q with no registered provider resource", resourceName)
		}
	}

	for resourceName := range fixtureResources {
		if _, ok := providerResources[resourceName]; !ok {
			t.Errorf("integration test fixture %q has no registered provider resource", resourceName)
		}
		if _, ok := mappedResources[resourceName]; !ok {
			t.Errorf("integration script map is missing resource %q used by test/terraform/resources", resourceName)
		}
	}
}

func providerRepoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine current filename")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func loadTerraformIntegrationScriptResourceMap(t *testing.T, scriptPath string) map[string]struct{} {
	t.Helper()

	// #nosec G304 -- scriptPath is derived from the repository root located via runtime.Caller.
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", scriptPath, err)
	}

	re := regexp.MustCompile(`(?m)^\s+"(netbox_[^"]+)"\s*=\s*@\{`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	if len(matches) == 0 {
		t.Fatalf("failed to parse any ResourceApiMap entries from %s", scriptPath)
	}

	entries := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		entries[match[1]] = struct{}{}
	}

	return entries
}

func registeredProviderResourceNames(t *testing.T) map[string]struct{} {
	t.Helper()

	p := New("test")().(*NetboxProvider)
	factories := p.Resources(context.Background())
	if len(factories) == 0 {
		t.Fatal("expected registered provider resources")
	}

	resourceNames := make(map[string]struct{}, len(factories))
	for _, factory := range factories {
		instance := factory()
		metadataResp := &resource.MetadataResponse{}
		instance.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "netbox"}, metadataResp)
		if metadataResp.TypeName == "" {
			t.Fatal("encountered empty provider resource type name")
		}
		resourceNames[metadataResp.TypeName] = struct{}{}
	}

	return resourceNames
}

func terraformIntegrationFixtureResources(t *testing.T, fixtureRoot string) map[string]struct{} {
	t.Helper()

	entries, err := os.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatalf("failed to read integration fixture root %s: %v", fixtureRoot, err)
	}

	resources := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		resources["netbox_"+entry.Name()] = struct{}{}
	}

	if len(resources) == 0 {
		t.Fatalf("expected integration resource fixtures under %s", fixtureRoot)
	}

	return resources
}
