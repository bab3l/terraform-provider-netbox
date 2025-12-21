package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestVirtualChassisResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	if r == nil {

		t.Fatal("Expected non-nil VirtualChassis resource")

	}

}

func TestVirtualChassisResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualChassisResource()

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

		Optional: []string{"domain", "master", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id", "member_count"},
	})

}

func TestVirtualChassisResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_virtual_chassis")

}

func TestVirtualChassisResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	testutil.ValidateResourceConfigure(t, r)

}
