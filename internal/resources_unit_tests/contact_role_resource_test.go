package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewContactRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestContactRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewContactRoleResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "slug"},

		Optional: []string{"description", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestContactRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewContactRoleResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_contact_role")

}

func TestContactRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewContactRoleResource()

	testutil.ValidateResourceConfigure(t, r)

}
