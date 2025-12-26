package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIPSecProfileResource(t *testing.T) {

	t.Parallel()

	r := resources.NewIPSecProfileResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPSecProfileResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewIPSecProfileResource()
	schemaRequest := &resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}
	r.Schema(context.Background(), *schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"name", "mode", "ike_policy", "ipsec_policy"},
		Optional: []string{"description", "comments"},
		Computed: []string{"id"},
	})
}

func TestIPSecProfileResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewIPSecProfileResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ipsec_profile")
}

func TestIPSecProfileResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewIPSecProfileResource()
	testutil.ValidateResourceConfigure(t, r)
}
