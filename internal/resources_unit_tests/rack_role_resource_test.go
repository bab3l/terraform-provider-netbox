package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRackRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRackRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

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

		Optional: []string{"color", "description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestRackRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_rack_role")

}

func TestRackRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

	testutil.ValidateResourceConfigure(t, r)

}
