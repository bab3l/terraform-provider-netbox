package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestPrefixesDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewPrefixesDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs:   []string{},
		ComputedAttrs: []string{"ids", "cidrs", "prefixes"},
	})

	block, ok := resp.Schema.Blocks["filter"]
	if !ok {
		t.Fatalf("Expected schema to define a 'filter' block")
	}

	setBlock, ok := block.(schema.SetNestedBlock)
	if !ok {
		t.Fatalf("Expected 'filter' to be schema.SetNestedBlock, got %T", block)
	}

	if _, ok := setBlock.NestedObject.Attributes["name"]; !ok {
		t.Fatalf("Expected filter block to include 'name' attribute")
	}
	if _, ok := setBlock.NestedObject.Attributes["values"]; !ok {
		t.Fatalf("Expected filter block to include 'values' attribute")
	}
}

func TestPrefixesDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewPrefixesDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_prefixes")
}

func TestPrefixesDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewPrefixesDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
