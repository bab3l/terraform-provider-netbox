// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SchemaValidation defines the expected attributes for a resource schema.

type SchemaValidation struct {
	Required []string // Attributes that must be required

	Optional []string // Attributes that must be optional

	Computed []string // Attributes that must be computed

	OptionalComputed []string // Attributes that must be both optional and computed

}

// ValidateResourceSchema checks that a resource schema matches the expected structure.

// It verifies that all required fields are present and marked as required,

// all optional fields are present and not marked as required,

// and all computed fields are present and marked as computed.

func ValidateResourceSchema(t *testing.T, schemaAttrs map[string]schema.Attribute, validation SchemaValidation) {

	t.Helper()

	// Check required attributes

	for _, attr := range validation.Required {

		fieldAttr, exists := schemaAttrs[attr]

		if !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

			continue

		}

		if !isRequired(fieldAttr) {

			t.Errorf("Expected attribute %s to be required", attr)

		}

	}

	// Check optional attributes

	for _, attr := range validation.Optional {

		fieldAttr, exists := schemaAttrs[attr]

		if !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

			continue

		}

		if isRequired(fieldAttr) {

			t.Errorf("Expected attribute %s to be optional, not required", attr)

		}

	}

	// Check computed attributes

	for _, attr := range validation.Computed {

		fieldAttr, exists := schemaAttrs[attr]

		if !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

			continue

		}

		if !isComputed(fieldAttr) {

			t.Errorf("Expected attribute %s to be computed", attr)

		}

	}

	// Check optional+computed attributes

	for _, attr := range validation.OptionalComputed {

		fieldAttr, exists := schemaAttrs[attr]

		if !exists {

			t.Errorf("Expected optional+computed attribute %s to exist in schema", attr)

			continue

		}

		if isRequired(fieldAttr) {

			t.Errorf("Expected attribute %s to be optional (not required)", attr)

		}

		if !isComputed(fieldAttr) {

			t.Errorf("Expected attribute %s to be computed", attr)

		}

	}

}

// isRequired checks if a schema attribute is marked as required.

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

	case schema.ListNestedAttribute:

		return a.Required

	case schema.SingleNestedAttribute:

		return a.Required

	default:

		return false

	}

}

// isComputed checks if a schema attribute is marked as computed.

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

	case schema.ListNestedAttribute:

		return a.Computed

	case schema.SingleNestedAttribute:

		return a.Computed

	default:

		return false

	}

}

// ValidateStringAttributeAllowsNull checks that an optional string attribute allows null values.

// This prevents the "was null, but now cty.StringVal("")" error.

func ValidateStringAttributeAllowsNull(t *testing.T, attr schema.Attribute, fieldName string) {

	t.Helper()

	stringAttr, ok := attr.(schema.StringAttribute)

	if !ok {

		t.Errorf("Field %s should be a StringAttribute", fieldName)

		return

	}

	if !stringAttr.Optional {

		t.Errorf("Field %s should be Optional to allow null values", fieldName)

	}

	if stringAttr.Required {

		t.Errorf("Field %s should not be Required (conflicts with Optional)", fieldName)

	}

}

// ValidateResourceMetadata checks that a resource's metadata is correctly configured.

// It verifies that the resource type name matches the expected name.

func ValidateResourceMetadata(t *testing.T, r resource.Resource, providerTypeName, expectedTypeName string) {

	t.Helper()

	metadataRequest := resource.MetadataRequest{

		ProviderTypeName: providerTypeName,
	}

	metadataResponse := &resource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	if metadataResponse.TypeName != expectedTypeName {

		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResponse.TypeName)

	}

}

// ValidateResourceConfigure checks that a resource's Configure method handles provider data correctly.

// It verifies that:

// 1. Configure succeeds with nil provider data (backwards compatibility)

// 2. Configure succeeds with a valid APIClient

// 3. Configure fails with invalid provider data (type assertion check)

func ValidateResourceConfigure(t *testing.T, r resource.Resource) {

	t.Helper()

	// Type assert to access Configure method

	configurable, ok := r.(resource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")

	}

	// Test 1: Configure with nil provider data (should succeed)

	configureRequest := resource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &resource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test 2: Configure with valid APIClient (should succeed)

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &resource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)

	}

	// Test 3: Configure with invalid provider data (should fail)

	configureRequest.ProviderData = InvalidProviderData

	configureResponse = &resource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}
