package datasources_unit_tests

import (
	"context"
	"strings"
	"testing"

	netboxprovider "github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestAllRegisteredDataSourcesHaveBaselineUnitCoverage(t *testing.T) {
	t.Parallel()

	p := netboxprovider.New("test")().(*netboxprovider.NetboxProvider)
	factories := p.DataSources(context.Background())
	if len(factories) == 0 {
		t.Fatal("expected registered data sources")
	}

	seen := make(map[string]struct{}, len(factories))

	for _, factory := range factories {
		d := factory()
		if d == nil {
			t.Fatal("expected non-nil data source instance")
		}

		metadataResp := &datasource.MetadataResponse{}
		d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "netbox"}, metadataResp)
		typeName := metadataResp.TypeName
		if typeName == "" {
			t.Fatal("expected data source metadata type name")
		}
		if !strings.HasPrefix(typeName, "netbox_") {
			t.Fatalf("expected data source type name %q to use netbox_ prefix", typeName)
		}
		if _, exists := seen[typeName]; exists {
			t.Fatalf("duplicate registered data source metadata name %q", typeName)
		}
		seen[typeName] = struct{}{}

		t.Run(typeName, func(t *testing.T) {
			t.Parallel()

			schemaResp := &datasource.SchemaResponse{}
			d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
			if schemaResp.Diagnostics.HasError() {
				t.Fatalf("schema returned diagnostics: %+v", schemaResp.Diagnostics)
			}
			if len(schemaResp.Schema.Attributes) == 0 {
				t.Fatal("expected schema attributes")
			}

			testutil.ValidateDataSourceMetadata(t, d, "netbox", typeName)
			testutil.ValidateDataSourceConfigure(t, d)
		})
	}
}
