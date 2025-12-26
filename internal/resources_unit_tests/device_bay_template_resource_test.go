package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceBayTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDeviceBayTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"device_type", "name"},
		Optional: []string{"label", "description"},
		Computed: []string{"id"},
	})
}

func TestDeviceBayTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device_bay_template")
}

func TestDeviceBayTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()
	testutil.ValidateResourceConfigure(t, r)
}
