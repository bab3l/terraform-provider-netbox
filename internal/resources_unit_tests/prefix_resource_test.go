package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPrefixResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource()
	if r == nil {
		t.Fatal("Expected non-nil Prefix resource")
	}
}

func TestPrefixResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"prefix"},
		Optional: []string{"status", "site", "vrf", "tenant", "vlan", "role", "is_pool", "mark_utilized", "description", "comments"},
		Computed: []string{"id"},
	})
}

func TestPrefixResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewPrefixResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_prefix")
}

func TestPrefixResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewPrefixResource().(*resources.PrefixResource)
	testutil.ValidateResourceConfigure(t, r)
}
