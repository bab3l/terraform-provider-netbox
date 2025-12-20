package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIKEPolicyResource(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEPolicyResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIKEPolicyResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEPolicyResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name"},
		Optional: []string{"description", "version", "mode", "proposals", "preshared_key", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestIKEPolicyResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEPolicyResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ike_policy")
}

func TestIKEPolicyResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEPolicyResource()
	testutil.ValidateResourceConfigure(t, r)
}
