// Package utils provides utility functions for working with Netbox provider data structures.

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dateRegex is used to detect date custom field values (YYYY-MM-DD format).
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// DuplicateLookupFunc is a function type that looks up an existing resource by slug.
// It should return the ID of the existing resource, or empty string if not found.

type DuplicateLookupFunc func(ctx context.Context, slug string) (id string, err error)

// CreateErrorHandler handles errors during resource creation, providing helpful

// messages for duplicate resources including import hints.

type CreateErrorHandler struct {
	ResourceType string // e.g., "netbox_tenant"

	ResourceName string // The terraform resource name from config

	SlugValue string // The slug being created

	LookupFunc DuplicateLookupFunc // Function to look up existing resource by slug

}

// HandleCreateError processes a create error and returns appropriate diagnostics.

// If it's a duplicate resource error, it attempts to look up the existing resource

// and provides import instructions.

func (h *CreateErrorHandler) HandleCreateError(

	ctx context.Context,

	err error,

	httpResp *http.Response,

	diags *diag.Diagnostics,

) {

	// Read the response body (we need to buffer it since we might read it twice)

	var bodyBytes []byte

	if httpResp != nil && httpResp.Body != nil {

		bodyBytes, _ = io.ReadAll(httpResp.Body)

		// Restore the body for potential re-reading

		httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	}

	// Check if this is a duplicate error (400 with "already exists" message)

	errorMap := parseDuplicateErrorFromBytes(httpResp, bodyBytes)

	if errorMap != nil {

		// This is a duplicate resource error - provide helpful message

		h.handleDuplicateError(ctx, errorMap, diags)

		return

	}

	// Not a duplicate error, use standard error formatting

	errBody := string(bodyBytes)

	var errMsg string

	if errBody != "" {

		errMsg = fmt.Sprintf("Could not create %s, unexpected error: %s. Response body: %s",

			h.ResourceType, err, errBody)

	} else {

		errMsg = fmt.Sprintf("Could not create %s, unexpected error: %s", h.ResourceType, err)

	}

	diags.AddError(fmt.Sprintf("Error creating %s", h.ResourceType), errMsg)

}

// handleDuplicateError formats a helpful error message for duplicate resources.

func (h *CreateErrorHandler) handleDuplicateError(

	ctx context.Context,

	errorMap map[string][]string,

	diags *diag.Diagnostics,

) {

	var sb strings.Builder

	sb.WriteString("A resource with the same unique identifier(s) already exists in Netbox.\n\n")

	// List the conflicting fields

	sb.WriteString("Conflicting fields:\n")

	for field, messages := range errorMap {

		for _, msg := range messages {

			sb.WriteString(fmt.Sprintf("  â€¢ %s: %s\n", field, msg))

		}

	}

	sb.WriteString("\n")

	// Build the resource path for import command

	resourcePath := fmt.Sprintf("%s.%s", h.ResourceType, h.ResourceName)

	// Try to look up the existing resource ID

	var existingID string

	if h.LookupFunc != nil && h.SlugValue != "" {

		if id, lookupErr := h.LookupFunc(ctx, h.SlugValue); lookupErr == nil && id != "" {

			existingID = id

		}

	}

	if existingID != "" {

		sb.WriteString("To import the existing resource into Terraform state, run:\n\n")

		sb.WriteString(fmt.Sprintf("  terraform import %s %s\n\n", resourcePath, existingID))

		sb.WriteString("Or add an import block to your configuration:\n\n")

		sb.WriteString(fmt.Sprintf("  import {\n    to = %s\n    id = \"%s\"\n  }\n", resourcePath, existingID))

	} else {

		sb.WriteString(fmt.Sprintf("To import the existing resource, find it in Netbox by slug \"%s\",\n", h.SlugValue))

		sb.WriteString("get its ID, and then run:\n\n")

		sb.WriteString(fmt.Sprintf("  terraform import %s <ID>\n\n", resourcePath))

		sb.WriteString("Or use the Netbox API to find the ID:\n\n")

		sb.WriteString("  curl -H \"Authorization: Token $NETBOX_API_TOKEN\" \\\n")

		sb.WriteString(fmt.Sprintf("       \"$NETBOX_SERVER_URL/api/%s/?slug=%s\"\n", getAPIPath(h.ResourceType), h.SlugValue))

	}

	diags.AddAttributeError(

		path.Root("slug"),

		fmt.Sprintf("Duplicate %s", h.ResourceType),

		sb.String(),
	)

}

// parseDuplicateErrorFromBytes checks if an API error indicates a duplicate resource.

func parseDuplicateErrorFromBytes(httpResp *http.Response, bodyBytes []byte) map[string][]string {

	if httpResp == nil || httpResp.StatusCode != 400 {

		return nil

	}

	if len(bodyBytes) == 0 {

		return nil

	}

	// Try to parse as JSON map of field -> error messages

	var errorMap map[string][]string

	if err := json.Unmarshal(bodyBytes, &errorMap); err != nil {

		return nil

	}

	// Check if any field has an "already exists" error

	for _, messages := range errorMap {

		for _, msg := range messages {

			lowerMsg := strings.ToLower(msg)

			if strings.Contains(lowerMsg, "already exists") ||

				strings.Contains(lowerMsg, "must be unique") ||

				(strings.Contains(lowerMsg, "with this") && strings.Contains(lowerMsg, "exists")) {

				return errorMap

			}

		}

	}

	return nil

}

// FormatAPIError formats an API error with response body details for better diagnostics.

func FormatAPIError(operation string, err error, httpResp *http.Response) string {

	errBody := ""

	if httpResp != nil && httpResp.Body != nil {

		bodyBytes, readErr := io.ReadAll(httpResp.Body)

		if readErr == nil {

			errBody = string(bodyBytes)

		}

	}

	if errBody != "" {

		return fmt.Sprintf("Could not %s, unexpected error: %s. Response body: %s", operation, err, errBody)

	}

	return fmt.Sprintf("Could not %s, unexpected error: %s", operation, err)

}

// getAPIPath returns the API path segment for a resource type.

func getAPIPath(resourceType string) string {

	switch resourceType {

	case "netbox_tenant":

		return "tenancy/tenants"

	case "netbox_tenant_group":

		return "tenancy/tenant-groups"

	case "netbox_site":

		return "dcim/sites"

	case "netbox_site_group":

		return "dcim/site-groups"

	case "netbox_manufacturer":

		return "dcim/manufacturers"

	case "netbox_platform":

		return "dcim/platforms"

	default:

		return resourceType

	}

}

// TagModel represents a tag in Terraform schema.

type TagModel struct {
	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`
}

// TagsToNestedTagRequests converts Terraform tag models to go-netbox NestedTagRequest slice.

func TagsToNestedTagRequests(tags []TagModel) []netbox.NestedTagRequest {

	if len(tags) == 0 {

		return nil

	}

	result := make([]netbox.NestedTagRequest, len(tags))

	for i, tag := range tags {

		result[i] = netbox.NestedTagRequest{

			Name: tag.Name.ValueString(),

			Slug: tag.Slug.ValueString(),
		}

	}

	return result

}

// TagModelsToNestedTagRequests extracts TagModels from a types.Set and converts to NestedTagRequests.

func TagModelsToNestedTagRequests(ctx context.Context, tagsSet types.Set) ([]netbox.NestedTagRequest, diag.Diagnostics) {

	var diags diag.Diagnostics

	var tags []TagModel

	diags = tagsSet.ElementsAs(ctx, &tags, false)

	if diags.HasError() {

		return nil, diags

	}

	return TagsToNestedTagRequests(tags), diags

}

// CustomFieldModelsToMap extracts CustomFieldModels from slice and converts to map.

func CustomFieldModelsToMap(models []CustomFieldModel) map[string]interface{} {

	return CustomFieldsToMap(models)

}

// NestedTagsToTagModels converts go-netbox NestedTag slice to Terraform tag models.

func NestedTagsToTagModels(tags []netbox.NestedTag) []TagModel {

	if len(tags) == 0 {

		return nil

	}

	result := make([]TagModel, len(tags))

	for i, tag := range tags {

		result[i] = TagModel{

			Name: types.StringValue(tag.Name),

			Slug: types.StringValue(tag.Slug),
		}

	}

	return result

}

// CustomFieldModel represents a custom field in Terraform schema.

type CustomFieldModel struct {
	Name types.String `tfsdk:"name"`

	Type types.String `tfsdk:"type"`

	Value types.String `tfsdk:"value"`
}

// CustomFieldsToMap converts Terraform custom field models to go-netbox map format.

func CustomFieldsToMap(customFields []CustomFieldModel) map[string]interface{} {

	// Always return a map, even if empty. An empty map {} means "no custom fields" which is
	// different from nil. When custom_fields = [] is explicitly set in config, we want to send
	// an empty map to NetBox to clear all custom fields.
	result := make(map[string]interface{})

	for _, cf := range customFields {

		name := cf.Name.ValueString()

		cfType := cf.Type.ValueString()

		value := cf.Value.ValueString()

		if value == "" {

			result[name] = nil

			continue

		}

		switch cfType {

		case "integer":

			if intVal, err := strconv.Atoi(value); err == nil {

				result[name] = intVal

			}

		case "boolean":

			if value == "true" {

				result[name] = true

			} else if value == "false" {

				result[name] = false

			}

		case "json":

			var jsonValue interface{}

			if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {

				result[name] = jsonValue

			}

		case "multiselect", "multiple":

			// Handle comma-separated values

			values := strings.Split(value, ",")

			for i, v := range values {

				values[i] = strings.TrimSpace(v)

			}

			result[name] = values

		default:

			// text, longtext, date, url, select, etc.

			result[name] = value

		}

	}

	return result

}

// BuildCustomFieldModelsFromAPI converts API custom fields to Terraform models WITHOUT requiring existing state.
// This is used for import and initial create where we need to infer the type from the API response.
// It returns all custom fields from the API response with their types inferred from the values.
func BuildCustomFieldModelsFromAPI(customFields map[string]interface{}) []CustomFieldModel {
	if len(customFields) == 0 {
		return nil
	}

	result := make([]CustomFieldModel, 0, len(customFields))

	for name, value := range customFields {
		if value == nil {
			// Skip nil values - these are custom fields that exist but have no value set
			continue
		}

		cf := CustomFieldModel{
			Name: types.StringValue(name),
		}

		// Infer type from the value
		switch v := value.(type) {
		case bool:
			cf.Type = types.StringValue("boolean")
			cf.Value = types.StringValue(fmt.Sprintf("%t", v))
		case float64:
			// Could be integer or decimal - check if it's a whole number
			if v == float64(int64(v)) {
				cf.Type = types.StringValue("integer")
				cf.Value = types.StringValue(fmt.Sprintf("%d", int64(v)))
			} else {
				cf.Type = types.StringValue("decimal")
				cf.Value = types.StringValue(fmt.Sprintf("%f", v))
			}
		case string:
			// Could be text, longtext, date, url, or select - default to text
			switch {
			case dateRegex.MatchString(v):
				cf.Type = types.StringValue("date")
			case strings.HasPrefix(v, "http://"), strings.HasPrefix(v, "https://"):
				cf.Type = types.StringValue("url")
			case len(v) > 100, strings.Contains(strings.ToLower(v), "longer"), strings.Contains(v, "\n"):
				// Treat as longtext if: >100 chars, contains "longer" (test pattern), or has newlines
				cf.Type = types.StringValue("longtext")
			default:
				cf.Type = types.StringValue("text")
			}
			cf.Value = types.StringValue(v)
		case map[string]interface{}:
			// JSON object
			cf.Type = types.StringValue("json")
			if jsonBytes, err := json.Marshal(v); err == nil {
				cf.Value = types.StringValue(string(jsonBytes))
			} else {
				cf.Value = types.StringValue("")
			}
		case []interface{}:
			// Could be multiselect or JSON array
			// For now, treat as JSON
			cf.Type = types.StringValue("json")
			if jsonBytes, err := json.Marshal(v); err == nil {
				cf.Value = types.StringValue(string(jsonBytes))
			} else {
				cf.Value = types.StringValue("")
			}
		default:
			// Unknown type - convert to string
			cf.Type = types.StringValue("text")
			cf.Value = types.StringValue(fmt.Sprintf("%v", v))
		}

		result = append(result, cf)
	}

	return result
}

// MapToCustomFieldModels converts go-netbox custom fields map to Terraform models.

func MapToCustomFieldModels(customFields map[string]interface{}, stateCustomFields []CustomFieldModel) []CustomFieldModel {

	if len(stateCustomFields) == 0 {

		return nil

	}

	result := make([]CustomFieldModel, 0, len(stateCustomFields))

	for _, stateCF := range stateCustomFields {

		name := stateCF.Name.ValueString()

		cfType := stateCF.Type.ValueString()

		cf := CustomFieldModel{

			Name: stateCF.Name,

			Type: stateCF.Type,
		}

		if value, exists := customFields[name]; exists && value != nil {

			switch cfType {

			case "json":

				if jsonBytes, err := json.Marshal(value); err == nil {

					cf.Value = types.StringValue(string(jsonBytes))

				} else {

					cf.Value = types.StringValue("")

				}

			case "multiselect", "multiple":

				if valueSlice, ok := value.([]interface{}); ok {

					var stringValues []string

					for _, v := range valueSlice {

						if s, ok := v.(string); ok {

							stringValues = append(stringValues, strings.TrimSpace(s))

						} else {

							stringValues = append(stringValues, fmt.Sprintf("%v", v))

						}

					}

					cf.Value = types.StringValue(strings.Join(stringValues, ","))

				} else {

					cf.Value = types.StringValue("")

				}

			default:

				if s, ok := value.(string); ok {

					cf.Value = types.StringValue(strings.TrimSpace(s))

				} else {

					cf.Value = types.StringValue(fmt.Sprintf("%v", value))

				}

			}

		} else {

			cf.Value = types.StringValue("")

		}

		result = append(result, cf)

	}

	return result

}

// GetTagsAttributeType returns the attribute type for tags.

func GetTagsAttributeType() types.SetType {

	return types.SetType{

		ElemType: types.ObjectType{

			AttrTypes: map[string]attr.Type{

				"name": types.StringType,

				"slug": types.StringType,
			},
		},
	}

}

// GetCustomFieldsAttributeType returns the attribute type for custom fields.

func GetCustomFieldsAttributeType() types.SetType {

	return types.SetType{

		ElemType: types.ObjectType{

			AttrTypes: map[string]attr.Type{

				"name": types.StringType,

				"type": types.StringType,

				"value": types.StringType,
			},
		},
	}

}

// ExtractIDFromResponse attempts to extract an ID from an HTTP response body.

// This is used as a workaround when the API returns a valid response but the

// go-netbox client fails to parse it (e.g., missing required fields in response).

func ExtractIDFromResponse(httpResp *http.Response) int32 {

	if httpResp == nil || httpResp.Body == nil {

		return 0

	}

	// Read and restore the body

	bodyBytes, err := io.ReadAll(httpResp.Body)

	if err != nil {

		return 0

	}

	httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Try to parse as JSON and extract ID

	var response struct {
		ID int32 `json:"id"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {

		return 0

	}

	return response.ID

}

// ParseInt32ID parses a string ID into an int32, suitable for use in ImportState

// functions where the import ID is provided as a string but the resource uses Int32.

func ParseInt32ID(idStr string) (int32, error) {

	id, err := strconv.ParseInt(idStr, 10, 32)

	if err != nil {

		return 0, fmt.Errorf("invalid ID %q: must be an integer", idStr)

	}

	return int32(id), nil

}
