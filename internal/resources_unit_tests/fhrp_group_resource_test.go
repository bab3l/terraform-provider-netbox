package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestFHRPGroupResource(t *testing.T) {
	t.Parallel()
	r := resources.NewFHRPGroupResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestFHRPGroupResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewFHRPGroupResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"protocol", "group_id"},
		Optional: []string{"name", "auth_type", "auth_key", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestFHRPGroupResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewFHRPGroupResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_fhrp_group")
}

func TestFHRPGroupResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewFHRPGroupResource()
	testutil.ValidateResourceConfigure(t, r)
}
