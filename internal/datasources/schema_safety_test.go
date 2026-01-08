package datasources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// dataSourceInfo contains metadata about a data source for testing.
type dataSourceInfo struct {
	name           string
	factory        func() datasource.DataSource
	searchFields   []string // fields that can be used to search (should be Optional)
	returnedFields []string // fields returned by the data source (should be Computed)
}

// allDataSources returns all data sources to test.
func allDataSources() []dataSourceInfo {
	return []dataSourceInfo{
		{
			name:           "netbox_tenant",
			factory:        datasources.NewTenantDataSource,
			searchFields:   []string{"id", "name", "slug"},
			returnedFields: []string{"description", "comments", "group", "group_id", "tags", "custom_fields"},
		},
		{
			name:           "netbox_tenant_group",
			factory:        datasources.NewTenantGroupDataSource,
			searchFields:   []string{"id", "name", "slug"},
			returnedFields: []string{"description", "parent", "parent_id", "tags", "custom_fields"},
		},
		{
			name:           "netbox_site",
			factory:        datasources.NewSiteDataSource,
			searchFields:   []string{"id", "name", "slug"},
			returnedFields: []string{"status", "description", "comments", "facility", "tenant", "tenant_id", "region", "region_id", "group", "group_id", "tags", "custom_fields"},
		},
		{
			name:           "netbox_site_group",
			factory:        datasources.NewSiteGroupDataSource,
			searchFields:   []string{"id", "name", "slug"},
			returnedFields: []string{"description", "parent", "tags", "custom_fields"},
		},
	}
}

// TestAllDataSourcesHaveSearchFields verifies all data sources have id, name, slug for searching.
func TestAllDataSourcesHaveSearchFields(t *testing.T) {
	t.Parallel()
	for _, dsi := range allDataSources() {
		t.Run(dsi.name, func(t *testing.T) {
			t.Parallel()
			ds := dsi.factory()
			schemaResp := &datasource.SchemaResponse{}
			ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
			if schemaResp.Diagnostics.HasError() {
				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)
			}
			for _, fieldName := range dsi.searchFields {
				attr, exists := schemaResp.Schema.Attributes[fieldName]
				if !exists {
					t.Errorf("Expected search field %s to exist", fieldName)
					continue
				}
				// Search fields should be Optional (can search by any one)
				stringAttr, ok := attr.(schema.StringAttribute)
				if !ok {
					t.Errorf("Search field %s should be a StringAttribute", fieldName)
					continue
				}
				if !stringAttr.Optional {
					t.Errorf("Search field %s should be Optional", fieldName)
				}

				// Search fields should also be computed (populated by the lookup)
				if !stringAttr.Computed {
					t.Errorf("Search field %s should be Computed (read-back after lookup)", fieldName)
				}
			}
		})
	}
}

// TestDataSourceReturnedFieldsAreComputed verifies returned fields are marked Computed.
func TestDataSourceReturnedFieldsAreComputed(t *testing.T) {
	t.Parallel()

	for _, dsi := range allDataSources() {
		t.Run(dsi.name, func(t *testing.T) {
			t.Parallel()
			ds := dsi.factory()
			schemaResp := &datasource.SchemaResponse{}
			ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
			if schemaResp.Diagnostics.HasError() {
				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)
			}
			for _, fieldName := range dsi.returnedFields {
				attr, exists := schemaResp.Schema.Attributes[fieldName]
				if !exists {
					// This might be expected for some data sources with fewer fields
					t.Logf("Optional field %s not present in data source", fieldName)
					continue
				}
				if !isDataSourceComputed(attr) {
					t.Errorf("Returned field %s should be Computed", fieldName)
				}
			}
		})
	}
}

// TestDataSourceMetadataPrefix verifies all data sources have the correct type name prefix.

func TestDataSourceMetadataPrefix(t *testing.T) {
	t.Parallel()

	for _, dsi := range allDataSources() {
		t.Run(dsi.name, func(t *testing.T) {
			t.Parallel()
			ds := dsi.factory()
			metaResp := &datasource.MetadataResponse{}
			ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "netbox"}, metaResp)
			if metaResp.TypeName != dsi.name {
				t.Errorf("Expected type name %s, got %s", dsi.name, metaResp.TypeName)
			}
		})
	}
}

// TestDataSourceSearchFieldsAreOptionalAndComputed verifies the pattern for search fields.
// Data sources should allow searching by any of id/name/slug (all Optional),
// and populate all of them after lookup (all Computed).
func TestDataSourceSearchFieldsAreOptionalAndComputed(t *testing.T) {
	t.Parallel()

	for _, dsi := range allDataSources() {
		t.Run(dsi.name, func(t *testing.T) {
			t.Parallel()
			ds := dsi.factory()
			schemaResp := &datasource.SchemaResponse{}
			ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
			if schemaResp.Diagnostics.HasError() {
				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)
			}
			for _, fieldName := range []string{"id", "name", "slug"} {
				attr, exists := schemaResp.Schema.Attributes[fieldName]
				if !exists {
					t.Errorf("Expected standard search field %s to exist", fieldName)
					continue
				}
				stringAttr, ok := attr.(schema.StringAttribute)
				if !ok {
					t.Errorf("Search field %s should be a StringAttribute", fieldName)
					continue
				}
				// Standard pattern for data source search fields
				if !stringAttr.Optional {
					t.Errorf("Search field %s should be Optional", fieldName)
				}
				if !stringAttr.Computed {
					t.Errorf("Search field %s should be Computed", fieldName)
				}
				if stringAttr.Required {
					t.Errorf("Search field %s should NOT be Required (conflicts with Optional)", fieldName)
				}
			}
		})
	}
}

// Helper to check if a data source attribute is computed.
func isDataSourceComputed(attr schema.Attribute) bool {
	switch a := attr.(type) {
	case schema.StringAttribute:
		return a.Computed
	case schema.Int64Attribute:
		return a.Computed
	case schema.BoolAttribute:
		return a.Computed
	case schema.SetAttribute:
		return a.Computed
	case schema.ListAttribute:
		return a.Computed
	case schema.MapAttribute:
		return a.Computed
	case schema.SetNestedAttribute:
		return a.Computed
	default:
		return false
	}
}
