package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTunnelTerminationResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	if r == nil {
		t.Fatal("Expected non-nil tunnel termination resource")
	}
}

func TestTunnelTerminationResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
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
		Required: []string{"tunnel", "termination_type"},
		Optional: []string{"role", "termination_id", "outside_ip", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestTunnelTerminationResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tunnel_termination")
}

func TestTunnelTerminationResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	testutil.ValidateResourceConfigure(t, r)
}
