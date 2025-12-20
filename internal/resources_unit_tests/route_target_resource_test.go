package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestRouteTargetResource(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	if r == nil {
		t.Fatal("Expected non-nil RouteTarget resource")
	}
}

func TestRouteTargetResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
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
		Required: []string{"name"},
		Optional: []string{"tenant", "description", "comments", "tags", "custom_fields"},
		Computed: []string{"id"},
	})
}

func TestRouteTargetResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_route_target")
}

func TestRouteTargetResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	testutil.ValidateResourceConfigure(t, r)
}
