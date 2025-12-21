package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactGroupResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestContactGroupResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactGroupResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "slug"},

		Optional: []string{"parent", "description", "tags", "custom_fields"},

		Computed: []string{"id", "parent_id"},
	})

}

func TestContactGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactGroupResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_contact_group")

}

func TestContactGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactGroupResource()

	testutil.ValidateResourceConfigure(t, r)

}
