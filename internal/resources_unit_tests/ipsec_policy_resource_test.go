package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIPSecPolicyResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPSecPolicyResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPSecPolicyResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPSecPolicyResource()
	schemaRequest := &resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}
	r.Schema(context.Background(), *schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name"},
		Optional: []string{"description", "proposals", "pfs_group", "comments"},
		Computed: []string{"id"},
	})
}

func TestIPSecPolicyResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPSecPolicyResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ipsec_policy")
}

func TestIPSecPolicyResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPSecPolicyResource()
	testutil.ValidateResourceConfigure(t, r)
}
