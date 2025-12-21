package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestL2VPNTerminationResource(t *testing.T) {

	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
	if r == nil {
		t.Fatal("Expected non-nil L2VPN Termination resource")
	}
}

func TestL2VPNTerminationResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
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
		Required: []string{"l2vpn", "assigned_object_type", "assigned_object_id"},
		Computed: []string{"id"},
	})
}

func TestL2VPNTerminationResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_l2vpn_termination")
}

func TestL2VPNTerminationResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
	testutil.ValidateResourceConfigure(t, r)
}
