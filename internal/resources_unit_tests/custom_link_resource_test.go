package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCustomLinkResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomLinkResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestCustomLinkResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomLinkResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "object_types", "link_text", "link_url"},

		Optional: []string{"enabled", "weight", "group_name", "button_class", "new_window"},

		Computed: []string{"id"},
	})

}

func TestCustomLinkResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomLinkResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_custom_link")

}

func TestCustomLinkResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomLinkResource()

	testutil.ValidateResourceConfigure(t, r)

}
