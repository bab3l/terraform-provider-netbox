// Package utils provides utility functions for working with Netbox provider data structures.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormatAPIError formats an API error with response body details for better diagnostics
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

// TagModel represents a tag in Terraform schema
type TagModel struct {
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
}

// TagsToNestedTagRequests converts Terraform tag models to go-netbox NestedTagRequest slice
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

// NestedTagsToTagModels converts go-netbox NestedTag slice to Terraform tag models
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

// CustomFieldModel represents a custom field in Terraform schema
type CustomFieldModel struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// CustomFieldsToMap converts Terraform custom field models to go-netbox map format
func CustomFieldsToMap(customFields []CustomFieldModel) map[string]interface{} {
	if len(customFields) == 0 {
		return nil
	}

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

// MapToCustomFieldModels converts go-netbox custom fields map to Terraform models
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
						stringValues = append(stringValues, strings.TrimSpace(v.(string)))
					}
					cf.Value = types.StringValue(strings.Join(stringValues, ","))
				} else {
					cf.Value = types.StringValue("")
				}
			default:
				cf.Value = types.StringValue(strings.TrimSpace(value.(string)))
			}
		} else {
			cf.Value = types.StringValue("")
		}

		result = append(result, cf)
	}

	return result
}

// GetTagsAttributeType returns the attribute type for tags
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

// GetCustomFieldsAttributeType returns the attribute type for custom fields
func GetCustomFieldsAttributeType() types.SetType {
	return types.SetType{
		ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"type":  types.StringType,
				"value": types.StringType,
			},
		},
	}
}
