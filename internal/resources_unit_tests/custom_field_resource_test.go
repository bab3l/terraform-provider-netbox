package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCustomFieldResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestCustomFieldResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"object_types", "type", "name"},

		Optional: []string{"related_object_type", "label", "group_name", "description", "required", "ui_visible", "ui_editable", "is_cloneable", "default", "weight", "validation_minimum", "validation_maximum", "validation_regex", "choice_set", "comments"},

		Computed: []string{"id"},
	})

}

func TestCustomFieldResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_custom_field")

}

func TestCustomFieldResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	testutil.ValidateResourceConfigure(t, r)

}
