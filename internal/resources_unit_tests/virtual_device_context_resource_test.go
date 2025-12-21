package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestVirtualDeviceContextResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualDeviceContextResource()

	if r == nil {

		t.Fatal("Expected non-nil VirtualDeviceContext resource")

	}

}

func TestVirtualDeviceContextResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualDeviceContextResource()

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

		Required: []string{"name", "device", "status"},

		Optional: []string{"identifier", "tenant", "primary_ip4", "primary_ip6", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestVirtualDeviceContextResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualDeviceContextResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_virtual_device_context")

}

func TestVirtualDeviceContextResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualDeviceContextResource()

	testutil.ValidateResourceConfigure(t, r)

}
