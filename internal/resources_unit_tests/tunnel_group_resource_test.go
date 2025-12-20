package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTunnelGroupResource(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelGroupResource()
	if r == nil {
		t.Fatal("Expected non-nil tunnel group resource")
	}
}

func TestTunnelGroupResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelGroupResource()
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
		Required: []string{"name", "slug"},
		Optional: []string{"description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestTunnelGroupResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelGroupResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_tunnel_group")
}

func TestTunnelGroupResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelGroupResource()
	testutil.ValidateResourceConfigure(t, r)
}
