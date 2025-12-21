package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConfigTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConfigTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "template_code"},

		Optional: []string{"description"},

		Computed: []string{"id", "data_path"},
	})

}

func TestConfigTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_config_template")

}

func TestConfigTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	testutil.ValidateResourceConfigure(t, r)

}
