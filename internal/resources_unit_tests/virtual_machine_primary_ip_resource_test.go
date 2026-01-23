package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestVirtualMachinePrimaryIPResource(t *testing.T) {
	t.Parallel()

	r := resources.NewVirtualMachinePrimaryIPResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestVirtualMachinePrimaryIPResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewVirtualMachinePrimaryIPResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"virtual_machine"},
		Optional: []string{"primary_ip4", "primary_ip6"},
		Computed: []string{"id"},
	})
}

func TestVirtualMachinePrimaryIPResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewVirtualMachinePrimaryIPResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_virtual_machine_primary_ip")
}

func TestVirtualMachinePrimaryIPResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewVirtualMachinePrimaryIPResource()
	testutil.ValidateResourceConfigure(t, r)
}
