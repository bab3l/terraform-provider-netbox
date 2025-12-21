package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestContactResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestContactResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name"},

		Optional: []string{"group", "description", "comments", "tags", "title", "phone", "email", "address", "link"},

		Computed: []string{"id"},
	})

}

func TestContactResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_contact")

}

func TestContactResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewContactResource()

	testutil.ValidateResourceConfigure(t, r)

}
