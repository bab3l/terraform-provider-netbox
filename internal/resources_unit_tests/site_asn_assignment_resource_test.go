package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestSiteASNAssignmentResource(t *testing.T) {
	t.Parallel()

	r := resources.NewSiteASNAssignmentResource()
	if r == nil {
		t.Fatal("Expected non-nil Site ASN Assignment resource")
	}
}

func TestSiteASNAssignmentResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewSiteASNAssignmentResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"site", "asn"},
		Computed: []string{"id"},
	})
}

func TestSiteASNAssignmentResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewSiteASNAssignmentResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_site_asn_assignment")
}

func TestSiteASNAssignmentResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewSiteASNAssignmentResource()
	testutil.ValidateResourceConfigure(t, r)
}
