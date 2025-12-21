package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestInterfaceResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInterfaceResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestInterfaceResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInterfaceResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device", "name", "type"},
		Optional: []string{"label", "enabled", "parent", "bridge", "lag", "mtu", "mac_address", "speed", "duplex", "wwn", "mgmt_only", "description", "mode", "mark_connected", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestInterfaceResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInterfaceResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_interface")
}

func TestInterfaceResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewInterfaceResource()
	testutil.ValidateResourceConfigure(t, r)
}
