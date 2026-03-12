package resources_unit_tests

import (
	"context"
	"strings"
	"testing"

	netboxprovider "github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAllRegisteredResourcesHaveBaselineUnitCoverage(t *testing.T) {
	t.Parallel()

	p := netboxprovider.New("test")().(*netboxprovider.NetboxProvider)
	factories := p.Resources(context.Background())
	if len(factories) == 0 {
		t.Fatal("expected registered resources")
	}

	seen := make(map[string]struct{}, len(factories))

	for _, factory := range factories {
		r := factory()
		if r == nil {
			t.Fatal("expected non-nil resource instance")
		}

		metadataResp := &fwresource.MetadataResponse{}
		r.Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "netbox"}, metadataResp)
		typeName := metadataResp.TypeName
		if typeName == "" {
			t.Fatal("expected resource metadata type name")
		}
		if !strings.HasPrefix(typeName, "netbox_") {
			t.Fatalf("expected resource type name %q to use netbox_ prefix", typeName)
		}
		if _, exists := seen[typeName]; exists {
			t.Fatalf("duplicate registered resource metadata name %q", typeName)
		}
		seen[typeName] = struct{}{}

		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			schemaResp := &fwresource.SchemaResponse{}
			r.Schema(context.Background(), fwresource.SchemaRequest{}, schemaResp)
			if schemaResp.Diagnostics.HasError() {
				t.Fatalf("schema returned diagnostics: %+v", schemaResp.Diagnostics)
			}
			if len(schemaResp.Schema.Attributes) == 0 {
				t.Fatal("expected schema attributes")
			}

			testutil.ValidateResourceMetadata(t, r, "netbox", typeName)
			testutil.ValidateResourceConfigure(t, r)
		})
	}
}
