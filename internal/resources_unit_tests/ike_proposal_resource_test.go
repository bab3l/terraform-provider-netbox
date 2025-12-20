package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIKEProposalResource(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEProposalResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIKEProposalResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEProposalResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "authentication_method", "encryption_algorithm", "group"},
		Optional: []string{"description", "authentication_algorithm", "sa_lifetime", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestIKEProposalResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEProposalResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ike_proposal")
}

func TestIKEProposalResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewIKEProposalResource()
	testutil.ValidateResourceConfigure(t, r)
}
