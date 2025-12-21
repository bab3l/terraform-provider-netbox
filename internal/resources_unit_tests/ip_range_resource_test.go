package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIPRangeResource(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPRangeResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPRangeResourceSchema(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPRangeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}
}

func TestIPRangeResourceMetadata(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPRangeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ip_range")
}

func TestIPRangeResourceConfigure(t *testing.T) {

	t.Parallel()
	t.Parallel()
	r := resources.NewIPRangeResource()
	testutil.ValidateResourceConfigure(t, r)
}
