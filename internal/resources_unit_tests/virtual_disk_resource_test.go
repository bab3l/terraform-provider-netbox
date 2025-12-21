package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestVirtualDiskResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	if r == nil {

		t.Fatal("Expected non-nil VirtualDisk resource")

	}

}

func TestVirtualDiskResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

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

		Required: []string{"virtual_machine", "name", "size"},

		Optional: []string{"description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestVirtualDiskResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_virtual_disk")

}

func TestVirtualDiskResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	testutil.ValidateResourceConfigure(t, r)

}
