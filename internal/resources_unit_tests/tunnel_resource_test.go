package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTunnelResource(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelResource()
	if r == nil {
		t.Fatal("Expected non-nil tunnel resource")
	}
}

func TestTunnelResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelResource()
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
		Required: []string{"name", "encapsulation"},
		Optional: []string{"status", "group", "ipsec_profile", "tenant", "tunnel_id", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestTunnelResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tunnel")
}

func TestTunnelResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()

	r := resources.NewTunnelResource()
	testutil.ValidateResourceConfigure(t, r)
}
