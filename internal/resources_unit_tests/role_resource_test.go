package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "slug"},

		Optional: []string{"weight", "description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_role")

}

func TestRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	testutil.ValidateResourceConfigure(t, r)

}
