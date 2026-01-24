// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"reflect"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschemapkg "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

	case schema.Int32Attribute:

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

	case schema.Int32Attribute:

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

// ValidateStringAttributeHasValidatorType checks that a resource string attribute includes a validator of the specified type.
func ValidateStringAttributeHasValidatorType(t *testing.T, attr schema.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	stringAttr, ok := attr.(schema.StringAttribute)
	if !ok {
		t.Errorf("Field %s should be a StringAttribute", fieldName)
		return
	}

	for _, v := range stringAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateDataSourceStringAttributeHasValidatorType checks that a data source string attribute includes a validator of the specified type.
func ValidateDataSourceStringAttributeHasValidatorType(t *testing.T, attr dsschemapkg.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	stringAttr, ok := attr.(dsschemapkg.StringAttribute)
	if !ok {
		t.Errorf("Field %s should be a StringAttribute", fieldName)
		return
	}

	for _, v := range stringAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateFloat64AttributeHasValidatorType checks that a resource float64 attribute includes a validator of the specified type.
func ValidateFloat64AttributeHasValidatorType(t *testing.T, attr schema.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	floatAttr, ok := attr.(schema.Float64Attribute)
	if !ok {
		t.Errorf("Field %s should be a Float64Attribute", fieldName)
		return
	}

	for _, v := range floatAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateInt64AttributeHasValidatorType checks that a resource int64 attribute includes a validator of the specified type.
func ValidateInt64AttributeHasValidatorType(t *testing.T, attr schema.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	intAttr, ok := attr.(schema.Int64Attribute)
	if !ok {
		t.Errorf("Field %s should be an Int64Attribute", fieldName)
		return
	}

	for _, v := range intAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateInt32AttributeHasValidatorType checks that a resource int32 attribute includes a validator of the specified type.
func ValidateInt32AttributeHasValidatorType(t *testing.T, attr schema.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	intAttr, ok := attr.(schema.Int32Attribute)
	if !ok {
		t.Errorf("Field %s should be an Int32Attribute", fieldName)
		return
	}

	for _, v := range intAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateDataSourceInt32AttributeHasValidatorType checks that a data source int32 attribute includes a validator of the specified type.
func ValidateDataSourceInt32AttributeHasValidatorType(t *testing.T, attr dsschemapkg.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	intAttr, ok := attr.(dsschemapkg.Int32Attribute)
	if !ok {
		t.Errorf("Field %s should be an Int32Attribute", fieldName)
		return
	}

	for _, v := range intAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
}

// ValidateDataSourceInt64AttributeHasValidatorType checks that a data source int64 attribute includes a validator of the specified type.
func ValidateDataSourceInt64AttributeHasValidatorType(t *testing.T, attr dsschemapkg.Attribute, fieldName string, expectedType reflect.Type) {
	t.Helper()

	intAttr, ok := attr.(dsschemapkg.Int64Attribute)
	if !ok {
		t.Errorf("Field %s should be an Int64Attribute", fieldName)
		return
	}

	for _, v := range intAttr.Validators {
		if reflect.TypeOf(v) == expectedType {
			return
		}
	}

	t.Errorf("Field %s is missing validator of type %s", fieldName, expectedType.String())
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

// DataSourceValidation defines the expected attributes for a datasource schema.
type DataSourceValidation struct {
	LookupAttrs   []string // Attributes that should be optional (used for lookups)
	ComputedAttrs []string // Attributes that should be computed
}

// ValidateDataSourceSchema checks that a datasource schema matches the expected structure.
// For datasources, all attributes are typically optional (for lookups) or computed.
func ValidateDataSourceSchema(t *testing.T, schemaAttrs map[string]dsschemapkg.Attribute, validation DataSourceValidation) {
	t.Helper()

	// Check lookup attributes (should be optional)
	for _, attr := range validation.LookupAttrs {
		fieldAttr, exists := schemaAttrs[attr]
		if !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
			continue
		}

		if !isDataSourceOptional(fieldAttr) {
			t.Errorf("Expected attribute %s to be optional for lookup", attr)
		}
	}

	// Check computed attributes
	for _, attr := range validation.ComputedAttrs {
		fieldAttr, exists := schemaAttrs[attr]
		if !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
			continue
		}

		if !isDataSourceComputed(fieldAttr) {
			t.Errorf("Expected attribute %s to be computed", attr)
		}
	}
}

// isDataSourceOptional checks if a datasource schema attribute is marked as optional.
func isDataSourceOptional(attr dsschemapkg.Attribute) bool {
	switch a := attr.(type) {
	case dsschemapkg.StringAttribute:
		return a.Optional
	case dsschemapkg.Int32Attribute:
		return a.Optional
	case dsschemapkg.Int64Attribute:
		return a.Optional
	case dsschemapkg.BoolAttribute:
		return a.Optional
	case dsschemapkg.SetAttribute:
		return a.Optional
	case dsschemapkg.ListAttribute:
		return a.Optional
	case dsschemapkg.MapAttribute:
		return a.Optional
	case dsschemapkg.SetNestedAttribute:
		return a.Optional
	case dsschemapkg.ListNestedAttribute:
		return a.Optional
	case dsschemapkg.SingleNestedAttribute:
		return a.Optional
	default:
		return false
	}
}

// isDataSourceComputed checks if a datasource schema attribute is marked as computed.
func isDataSourceComputed(attr dsschemapkg.Attribute) bool {
	switch a := attr.(type) {
	case dsschemapkg.StringAttribute:
		return a.Computed
	case dsschemapkg.Int32Attribute:
		return a.Computed
	case dsschemapkg.Int64Attribute:
		return a.Computed
	case dsschemapkg.BoolAttribute:
		return a.Computed
	case dsschemapkg.SetAttribute:
		return a.Computed
	case dsschemapkg.ListAttribute:
		return a.Computed
	case dsschemapkg.MapAttribute:
		return a.Computed
	case dsschemapkg.SetNestedAttribute:
		return a.Computed
	case dsschemapkg.ListNestedAttribute:
		return a.Computed
	case dsschemapkg.SingleNestedAttribute:
		return a.Computed
	default:
		return false
	}
}

// ValidateDataSourceMetadata checks that a datasource's metadata is correctly configured.
// It verifies that the datasource type name matches the expected name.
func ValidateDataSourceMetadata(t *testing.T, d datasource.DataSource, providerTypeName, expectedTypeName string) {
	t.Helper()

	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: providerTypeName,
	}

	metadataResponse := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	if metadataResponse.TypeName != expectedTypeName {
		t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResponse.TypeName)
	}
}

// ValidateDataSourceConfigure checks that a datasource's Configure method handles provider data correctly.
// It verifies that:
// 1. Configure succeeds with nil provider data (backwards compatibility)
// 2. Configure succeeds with a valid APIClient
// 3. Configure fails with invalid provider data (type assertion check).
func ValidateDataSourceConfigure(t *testing.T, d datasource.DataSource) {
	t.Helper()

	// Type assert to access Configure method
	configurable, ok := d.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("DataSource does not implement DataSourceWithConfigure")
	}

	// Test 1: Configure with nil provider data (should succeed)
	configureRequest := datasource.ConfigureRequest{
		ProviderData: nil,
	}

	configureResponse := &datasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test 2: Configure with valid APIClient (should succeed)
	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &datasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)
	}

	// Test 3: Configure with invalid provider data (should fail)
	configureRequest.ProviderData = InvalidProviderData

	configureResponse = &datasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}
