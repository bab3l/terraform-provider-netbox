package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestFHRPGroupAssignmentResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestFHRPGroupAssignmentResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"group_id", "interface_type", "interface_id", "priority"},
		Optional: []string{},
		Computed: []string{"id"},
	})
}

func TestFHRPGroupAssignmentResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_fhrp_group_assignment")
}

func TestFHRPGroupAssignmentResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()
	testutil.ValidateResourceConfigure(t, r)
}
