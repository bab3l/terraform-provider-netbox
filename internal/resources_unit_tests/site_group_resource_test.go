package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSiteGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteGroupResource()
	if r == nil {
		t.Fatal("Expected non-nil SiteGroup resource")
	}
}

func TestSiteGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteGroupResource()
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
		Optional: []string{"parent", "description", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestSiteGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteGroupResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_site_group")
}

func TestSiteGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewSiteGroupResource()
	testutil.ValidateResourceConfigure(t, r)
}
