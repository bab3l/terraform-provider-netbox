package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestPlatformResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPlatformResource()
	if r == nil {
		t.Fatal("Expected non-nil platform resource")
	}
}

func TestPlatformResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPlatformResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required:         []string{"name", "slug"},
		Optional:         []string{"manufacturer", "description"},
		Computed:         []string{"id"},
		OptionalComputed: []string{},
	})
}

func TestPlatformResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPlatformResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_platform")

}

func TestPlatformResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPlatformResource()

	testutil.ValidateResourceConfigure(t, r)

}
