package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCustomFieldChoiceSetResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestCustomFieldChoiceSetResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "extra_choices"},

		Optional: []string{"description", "base_choices", "order_alphabetically"},

		Computed: []string{"id"},
	})

}

func TestCustomFieldChoiceSetResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_custom_field_choice_set")

}

func TestCustomFieldChoiceSetResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()

	testutil.ValidateResourceConfigure(t, r)

}
