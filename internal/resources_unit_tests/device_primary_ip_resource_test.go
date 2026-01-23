package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDevicePrimaryIPResource(t *testing.T) {
	t.Parallel()

	r := resources.NewDevicePrimaryIPResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDevicePrimaryIPResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewDevicePrimaryIPResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device"},
		Optional: []string{"primary_ip4", "primary_ip6", "oob_ip"},
		Computed: []string{"id"},
	})
}

func TestDevicePrimaryIPResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewDevicePrimaryIPResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device_primary_ip")
}

func TestDevicePrimaryIPResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewDevicePrimaryIPResource()
	testutil.ValidateResourceConfigure(t, r)
}
