// Package resources_test contains safety tests to prevent common Terraform provider issues.

//

// These tests verify schema correctness and catch patterns that lead to runtime errors:

// 1. Optional fields must handle null vs empty string correctly

// 2. Required fields must always be present

// 3. Computed fields must be properly marked

// 4. All resources must have consistent schema patterns

package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// resourceInfo contains metadata about a resource for testing.

type resourceInfo struct {
	name string

	factory func() resource.Resource

	requiredFields []string

	optionalFields []string

	computedFields []string

	optionalComputed []string // fields that are both optional and computed

}

// allResources returns all resources to test with their expected schema.

func allResources() []resourceInfo {

	return []resourceInfo{

		{

			name: "netbox_tenant",

			factory: resources.NewTenantResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"group", "description", "comments", "tags", "custom_fields"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},

		{

			name: "netbox_tenant_group",

			factory: resources.NewTenantGroupResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"parent", "description", "tags", "custom_fields"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},

		{

			name: "netbox_site",

			factory: resources.NewSiteResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"status", "region", "group", "tenant", "facility", "description", "comments", "tags", "custom_fields"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},

		{

			name: "netbox_site_group",

			factory: resources.NewSiteGroupResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"parent", "description", "tags", "custom_fields"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},

		{

			name: "netbox_manufacturer",

			factory: resources.NewManufacturerResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"description"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},

		{

			name: "netbox_platform",

			factory: resources.NewPlatformResource,

			requiredFields: []string{"name", "slug"},

			optionalFields: []string{"manufacturer", "description"},

			computedFields: []string{"id"},

			optionalComputed: []string{},
		},
	}

}

// TestAllResourcesHaveIDField verifies all resources have a computed "id" field.

// This is a Terraform requirement for resource state management.

func TestAllResourcesHaveIDField(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			idAttr, exists := schemaResp.Schema.Attributes["id"]

			if !exists {

				t.Fatal("Resource must have an 'id' attribute")

			}

			// Verify id is computed (read-only, set by provider)

			if stringAttr, ok := idAttr.(schema.StringAttribute); ok {

				if !stringAttr.Computed {

					t.Error("id attribute must be Computed")

				}

			} else {

				t.Error("id attribute must be a StringAttribute")

			}

		})

	}

}

// TestOptionalFieldsAreNotRequired verifies optional fields are correctly marked.

// Incorrectly marking optional fields as required causes plan failures.

func TestOptionalFieldsAreNotRequired(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			for _, fieldName := range ri.optionalFields {

				attr, exists := schemaResp.Schema.Attributes[fieldName]

				if !exists {

					t.Errorf("Expected optional field %s to exist", fieldName)

					continue

				}

				// Check that the attribute is Optional (not Required)

				if isRequired(attr) {

					t.Errorf("Field %s should be Optional, not Required", fieldName)

				}

			}

		})

	}

}

// TestRequiredFieldsAreMarkedRequired verifies required fields are correctly marked.

// Missing required field markers causes Terraform to accept incomplete configs.

func TestRequiredFieldsAreMarkedRequired(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			for _, fieldName := range ri.requiredFields {

				attr, exists := schemaResp.Schema.Attributes[fieldName]

				if !exists {

					t.Errorf("Expected required field %s to exist", fieldName)

					continue

				}

				if !isRequired(attr) {

					t.Errorf("Field %s should be Required", fieldName)

				}

			}

		})

	}

}

// TestComputedFieldsAreMarkedComputed verifies computed fields are correctly marked.

// Fields set by the provider (like id) must be Computed.

func TestComputedFieldsAreMarkedComputed(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			for _, fieldName := range ri.computedFields {

				attr, exists := schemaResp.Schema.Attributes[fieldName]

				if !exists {

					t.Errorf("Expected computed field %s to exist", fieldName)

					continue

				}

				if !isComputed(attr) {

					t.Errorf("Field %s should be Computed", fieldName)

				}

			}

		})

	}

}

// TestOptionalStringFieldsAllowNull verifies optional string fields can be null.

// This prevents the "was null, but now cty.StringVal("")" error.

func TestOptionalStringFieldsAllowNull(t *testing.T) {

	t.Parallel()

	// These are the optional string fields that commonly cause null vs "" issues

	stringFieldsToCheck := map[string][]string{

		"netbox_tenant": {"description", "comments"},

		"netbox_tenant_group": {"description"},

		"netbox_site": {"description", "comments", "facility"},

		"netbox_site_group": {"description"},

		"netbox_manufacturer": {"description"},

		// netbox_platform has no optional string fields currently

	}

	for _, ri := range allResources() {

		fieldsToCheck, ok := stringFieldsToCheck[ri.name]

		if !ok {

			continue

		}

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			for _, fieldName := range fieldsToCheck {

				attr, exists := schemaResp.Schema.Attributes[fieldName]

				if !exists {

					t.Errorf("Expected field %s to exist", fieldName)

					continue

				}

				stringAttr, ok := attr.(schema.StringAttribute)

				if !ok {

					t.Errorf("Field %s should be a StringAttribute", fieldName)

					continue

				}

				// Optional fields should be Optional=true, Required=false

				if !stringAttr.Optional {

					t.Errorf("Field %s should be Optional to allow null values", fieldName)

				}

				if stringAttr.Required {

					t.Errorf("Field %s should not be Required (conflicts with Optional)", fieldName)

				}

			}

		})

	}

}

// TestResourceMetadataPrefix verifies all resources have the correct type name prefix.

func TestResourceMetadataPrefix(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			metaResp := &resource.MetadataResponse{}

			r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "netbox"}, metaResp)

			if metaResp.TypeName != ri.name {

				t.Errorf("Expected type name %s, got %s", ri.name, metaResp.TypeName)

			}

		})

	}

}

// TestAllDeclaredFieldsExistInSchema verifies all fields we expect are present.

// This catches accidental field removals.

func TestAllDeclaredFieldsExistInSchema(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			allFields := append(append(append(

				ri.requiredFields,

				ri.optionalFields...),

				ri.computedFields...),

				ri.optionalComputed...)

			for _, fieldName := range allFields {

				if _, exists := schemaResp.Schema.Attributes[fieldName]; !exists {

					t.Errorf("Expected field %s to exist in schema", fieldName)

				}

			}

		})

	}

}

// TestNoUnexpectedFieldsInSchema warns about fields not in our expected list.

// This helps catch new fields that need to be added to tests.

func TestNoUnexpectedFieldsInSchema(t *testing.T) {

	t.Parallel()

	for _, ri := range allResources() {

		t.Run(ri.name, func(t *testing.T) {

			t.Parallel()

			r := ri.factory()

			schemaResp := &resource.SchemaResponse{}

			r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

			if schemaResp.Diagnostics.HasError() {

				t.Fatalf("Schema error: %+v", schemaResp.Diagnostics)

			}

			expectedFields := make(map[string]bool)

			for _, f := range ri.requiredFields {

				expectedFields[f] = true

			}

			for _, f := range ri.optionalFields {

				expectedFields[f] = true

			}

			for _, f := range ri.computedFields {

				expectedFields[f] = true

			}

			for _, f := range ri.optionalComputed {

				expectedFields[f] = true

			}

			for fieldName := range schemaResp.Schema.Attributes {

				if !expectedFields[fieldName] {

					t.Logf("WARNING: Unexpected field %s in schema - add to test expectations", fieldName)

				}

			}

		})

	}

}

// Helper functions to check attribute properties.

func isRequired(attr schema.Attribute) bool {

	switch a := attr.(type) {

	case schema.StringAttribute:

		return a.Required

	case schema.Int64Attribute:

		return a.Required

	case schema.BoolAttribute:

		return a.Required

	case schema.SetAttribute:

		return a.Required

	case schema.ListAttribute:

		return a.Required

	case schema.MapAttribute:

		return a.Required

	case schema.SetNestedAttribute:

		return a.Required

	default:

		return false

	}

}

func isComputed(attr schema.Attribute) bool {

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
