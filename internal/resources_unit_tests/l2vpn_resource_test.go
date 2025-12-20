package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestL2VPNResource(t *testing.T) {
	t.Parallel()
	r := resources.NewL2VPNResource()
	if r == nil {
		t.Fatal("Expected non-nil L2VPN resource")
	}
}

func TestL2VPNResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewL2VPNResource()
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
		Required: []string{"name", "slug", "type"},
		Optional: []string{"identifier", "tenant", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestL2VPNResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewL2VPNResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_l2vpn")
}

func TestL2VPNResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewL2VPNResource()
	testutil.ValidateResourceConfigure(t, r)
}
