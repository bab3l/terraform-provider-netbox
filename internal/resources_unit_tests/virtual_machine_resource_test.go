package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestVirtualMachineResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	if r == nil {

		t.Fatal("Expected non-nil Virtual Machine resource")

	}

}

func TestVirtualMachineResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualMachineResource()

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

		Required: []string{"name"},

		Optional: []string{"status", "cluster", "vcpus", "memory", "disk", "description", "comments"},

		Computed: []string{"id"},
	})

}

func TestVirtualMachineResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_virtual_machine")

}

func TestVirtualMachineResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	testutil.ValidateResourceConfigure(t, r)

}
