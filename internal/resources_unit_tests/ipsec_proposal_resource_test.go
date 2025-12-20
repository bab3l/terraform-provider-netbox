package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIPSecProposalResource(t *testing.T) {
	t.Parallel()
	r := resources.NewIPSecProposalResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPSecProposalResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewIPSecProposalResource()
	schemaRequest := &resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}
	r.Schema(context.Background(), *schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name"},
		Optional: []string{"description", "encryption_algorithm", "authentication_algorithm", "sa_lifetime_seconds", "sa_lifetime_data", "comments"},
		Computed: []string{"id"},
	})
}

func TestIPSecProposalResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewIPSecProposalResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ipsec_proposal")
}

func TestIPSecProposalResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewIPSecProposalResource()
	testutil.ValidateResourceConfigure(t, r)
}
