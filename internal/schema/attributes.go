// Package schema provides reusable schema attribute factories for Terraform resources and data sources.
//

// This package reduces boilerplate by providing pre-configured schema attributes
// for common patterns used throughout the Netbox provider.

package schema

import (
	"regexp"

	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// =====================================================
// RESOURCE SCHEMA ATTRIBUTES

// =====================================================
// IDAttribute returns the standard ID attribute for resources.

// This is a computed string field that uses UseStateForUnknown.

func IDAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,

		MarkdownDescription: "Unique identifier for the " + resourceName + " (assigned by Netbox).",

		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// ComputedIDAttribute returns a computed ID attribute for reference fields.
// This is used to store the resolved ID of a referenced resource.

func ComputedIDAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "The numeric ID of the " + resourceName + ".",

		Computed: true,

		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// NameAttribute returns a required name attribute with standard validation.

func NameAttribute(resourceName string, maxLength int) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Name of the " + resourceName + ".",

		Required: true,

		Validators: []validator.String{
			stringvalidator.LengthBetween(1, maxLength),
		},
	}
}

// OptionalNameAttribute returns an optional name attribute with standard validation.

func OptionalNameAttribute(resourceName string, maxLength int) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Name of the " + resourceName + ".",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(maxLength),
		},
	}
}

// ModelAttribute returns a required model name attribute with standard validation.
// Used for device types where the model field identifies the device.

func ModelAttribute(resourceName string, maxLength int) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Model name/number of the " + resourceName + ".",

		Required: true,

		Validators: []validator.String{
			stringvalidator.LengthBetween(1, maxLength),
		},
	}
}

// SlugAttribute returns a required slug attribute with standard validation.

func SlugAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "URL-friendly identifier for the " + resourceName + ". Must be unique and contain only lowercase letters, numbers, hyphens, and underscores.",

		Required: true,

		Validators: []validator.String{
			stringvalidator.LengthBetween(1, 100),

			validators.ValidSlug(),
		},
	}
}

// DescriptionAttribute returns an optional description attribute.

func DescriptionAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Description of the " + resourceName + ".",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(200),
		},
	}
}

// CommentsAttribute returns an optional comments attribute.

func CommentsAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Additional comments or notes about the " + resourceName + ". Supports Markdown formatting.",

		Optional: true,
	}
}

// CommentsAttributeWithLimit returns an optional comments attribute with a length limit.

func CommentsAttributeWithLimit(resourceName string, maxLength int) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Additional comments or notes about the " + resourceName + ". Supports Markdown formatting.",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(maxLength),
		},
	}
}

// ReferenceAttribute returns an optional reference attribute (ID or slug lookup).
// Use this for foreign key relationships like tenant, site, location, etc.

func ReferenceAttribute(targetResource string, description string) schema.StringAttribute {
	if description == "" {
		description = "ID or slug of the " + targetResource + "."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,

		Optional: true,
	}
}

// RequiredReferenceAttribute returns a required reference attribute.

func RequiredReferenceAttribute(targetResource string, description string) schema.StringAttribute {
	if description == "" {
		description = "ID or slug of the " + targetResource + ". Required."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,

		Required: true,
	}
}

// IDOnlyReferenceAttribute returns an optional reference that only accepts integer IDs.

func IDOnlyReferenceAttribute(targetResource string, description string) schema.StringAttribute {
	if description == "" {
		description = "ID of the " + targetResource + "."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,

		Optional: true,

		Validators: []validator.String{
			stringvalidator.RegexMatches(

				validators.IntegerRegex(),

				"must be a valid integer ID",
			),
		},
	}
}

// StatusAttribute returns a status enum attribute with the given valid values.
// The first value in the list is used as the default.

func StatusAttribute(validValues []string, description string) schema.StringAttribute {
	defaultValue := "active"

	if len(validValues) > 0 {
		// Find "active" in list, or use first value

		for _, v := range validValues {
			if v == "active" {
				defaultValue = v

				break
			}
		}

		if defaultValue == "" {
			defaultValue = validValues[0]
		}
	}

	return schema.StringAttribute{
		MarkdownDescription: description,

		Optional: true,

		Computed: true,

		Default: stringdefault.StaticString(defaultValue),

		Validators: []validator.String{
			stringvalidator.OneOf(validValues...),
		},
	}
}

// EnumAttribute returns an optional enum attribute with the given valid values.

func EnumAttribute(description string, validValues []string) schema.StringAttribute {
	// Allow empty string for optional clearing

	valuesWithEmpty := append([]string{""}, validValues...)

	return schema.StringAttribute{
		MarkdownDescription: description,

		Optional: true,

		Validators: []validator.String{
			stringvalidator.OneOf(valuesWithEmpty...),
		},
	}
}

// RequiredEnumAttribute returns a required enum attribute with the given valid values.

func RequiredEnumAttribute(description string, validValues []string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: description,

		Required: true,

		Validators: []validator.String{
			stringvalidator.OneOf(validValues...),
		},
	}
}

// ColorAttribute returns an optional color attribute (6-character hex without #).

func ColorAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Color for the " + resourceName + " in 6-character hexadecimal format (without #). Example: 'aa1409'.",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthBetween(6, 6),

			stringvalidator.RegexMatches(

				regexp.MustCompile(`^[0-9a-fA-F]{6}$`),

				"must be exactly 6 hexadecimal characters (0-9, a-f)",
			),
		},
	}
}

// ComputedColorAttribute returns a computed color attribute (API assigns default).

func ComputedColorAttribute(resourceName string) schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Color for the " + resourceName + " in 6-character hexadecimal format (without #). Example: 'aa1409'. If not specified, Netbox assigns a default.",

		Optional: true,

		Computed: true,

		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},

		Validators: []validator.String{
			stringvalidator.LengthBetween(6, 6),

			stringvalidator.RegexMatches(

				regexp.MustCompile(`^[0-9a-fA-F]{6}$`),

				"must be exactly 6 hexadecimal characters (0-9, a-f)",
			),
		},
	}
}

// BoolAttributeWithDefault returns an optional bool attribute with a default value.
// This is useful for fields like vm_role that have a boolean with a default.

func BoolAttributeWithDefault(description string, defaultValue bool) schema.BoolAttribute {
	return schema.BoolAttribute{
		MarkdownDescription: description,

		Optional: true,

		Computed: true,

		Default: booldefault.StaticBool(defaultValue),
	}
}

// SerialAttribute returns an optional serial number attribute.

func SerialAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Serial number, assigned by the manufacturer.",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(50),
		},
	}
}

// AssetTagAttribute returns an optional asset tag attribute.

func AssetTagAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "A unique tag used for asset tracking.",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(50),
		},
	}
}

// FacilityAttribute returns an optional facility identifier attribute.

func FacilityAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Local facility ID or description.",

		Optional: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(50),
		},
	}
}

// TagsAttribute returns the standard tags set attribute for resources.

func TagsAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Tags assigned to this resource. Tags must already exist in Netbox.",

		Optional: true,

		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Name of the existing tag.",

					Required: true,

					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 100),
					},
				},

				"slug": schema.StringAttribute{
					MarkdownDescription: "Slug of the existing tag.",

					Required: true,

					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 100),

						validators.ValidSlug(),
					},
				},
			},
		},
	}
}

// CustomFieldsAttribute returns the standard custom fields set attribute for resources.

func CustomFieldsAttribute() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Custom fields assigned to this resource. Custom fields must be defined in Netbox before use.",

		Optional: true,

		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Name of the custom field.",

					Required: true,

					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 50),

						stringvalidator.RegexMatches(

							regexp.MustCompile(`^[a-z0-9_]+$`),

							"must contain only lowercase letters, numbers, and underscores",
						),
					},
				},

				"type": schema.StringAttribute{
					MarkdownDescription: "Type of the custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject).",

					Required: true,

					Validators: []validator.String{
						stringvalidator.OneOf(

							"text",

							"longtext",

							"integer",

							"boolean",

							"date",

							"url",

							"json",

							"select",

							"multiselect",

							"object",

							"multiobject",

							"multiple", // legacy

							"selection", // legacy

						),
					},
				},

				"value": schema.StringAttribute{
					MarkdownDescription: "Value of the custom field.",

					Required: true,

					Validators: []validator.String{
						stringvalidator.LengthAtMost(1000),
					},
				},
			},
		},
	}
}

// =====================================================
// DATA SOURCE SCHEMA ATTRIBUTES

// =====================================================
// DSIDAttribute returns the standard ID attribute for data sources (optional/computed).

func DSIDAttribute(resourceName string) dsschema.StringAttribute {
	return dsschema.StringAttribute{
		MarkdownDescription: "Unique identifier for the " + resourceName + ". Use to look up by ID.",

		Optional: true,

		Computed: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(50),
		},
	}
}

// DSNameAttribute returns a name attribute for data sources (optional/computed for lookup).

func DSNameAttribute(resourceName string) dsschema.StringAttribute {
	return dsschema.StringAttribute{
		MarkdownDescription: "Name of the " + resourceName + ". Use to look up by name.",

		Optional: true,

		Computed: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(100),
		},
	}
}

// DSSlugAttribute returns a slug attribute for data sources (optional/computed for lookup).

func DSSlugAttribute(resourceName string) dsschema.StringAttribute {
	return dsschema.StringAttribute{
		MarkdownDescription: "URL-friendly identifier for the " + resourceName + ". Use to look up by slug.",

		Optional: true,

		Computed: true,

		Validators: []validator.String{
			stringvalidator.LengthAtMost(100),
		},
	}
}

// DSComputedStringAttribute returns a computed-only string attribute for data sources.

func DSComputedStringAttribute(description string) dsschema.StringAttribute {
	return dsschema.StringAttribute{
		MarkdownDescription: description,

		Computed: true,
	}
}

// DSComputedBoolAttribute returns a computed-only bool attribute for data sources.

func DSComputedBoolAttribute(description string) dsschema.BoolAttribute {
	return dsschema.BoolAttribute{
		MarkdownDescription: description,

		Computed: true,
	}
}

// DSComputedInt64Attribute returns a computed-only int64 attribute for data sources.

func DSComputedInt64Attribute(description string) dsschema.Int64Attribute {
	return dsschema.Int64Attribute{
		MarkdownDescription: description,

		Computed: true,
	}
}

// DSComputedFloat64Attribute returns a computed-only float64 attribute for data sources.

func DSComputedFloat64Attribute(description string) dsschema.Float64Attribute {
	return dsschema.Float64Attribute{
		MarkdownDescription: description,

		Computed: true,
	}
}

// DSTagsAttribute returns the standard tags set attribute for data sources (computed).

func DSTagsAttribute() dsschema.SetNestedAttribute {
	return dsschema.SetNestedAttribute{
		MarkdownDescription: "Tags assigned to this resource.",

		Computed: true,

		NestedObject: dsschema.NestedAttributeObject{
			Attributes: map[string]dsschema.Attribute{
				"name": dsschema.StringAttribute{
					MarkdownDescription: "Name of the tag.",

					Computed: true,
				},

				"slug": dsschema.StringAttribute{
					MarkdownDescription: "Slug of the tag.",

					Computed: true,
				},
			},
		},
	}
}

// DSCustomFieldsAttribute returns the standard custom fields set attribute for data sources (computed).

func DSCustomFieldsAttribute() dsschema.SetNestedAttribute {
	return dsschema.SetNestedAttribute{
		MarkdownDescription: "Custom fields assigned to this resource.",

		Computed: true,

		NestedObject: dsschema.NestedAttributeObject{
			Attributes: map[string]dsschema.Attribute{
				"name": dsschema.StringAttribute{
					MarkdownDescription: "Name of the custom field.",

					Computed: true,
				},

				"type": dsschema.StringAttribute{
					MarkdownDescription: "Type of the custom field.",

					Computed: true,
				},

				"value": dsschema.StringAttribute{
					MarkdownDescription: "Value of the custom field.",

					Computed: true,
				},
			},
		},
	}
}

// =====================================================
// SCHEMA COMPOSITION HELPERS

// =====================================================
// These helpers provide pre-composed attribute sets that are commonly used

// together across many resources, reducing schema definition boilerplate.
// DescriptionOnlyAttributes returns just the description attribute.
// Use this for resources that have description but not comments.
//
// Usage:
//
//	attrs := map[string]schema.Attribute{
//	    "id": IDAttribute("resource"),
//	    "name": NameAttribute("resource", 100),
//	}
//	maps.Copy(attrs, DescriptionOnlyAttributes("resource"))
func DescriptionOnlyAttributes(resourceName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"description": DescriptionAttribute(resourceName),
	}
}

// CommonDescriptiveAttributes returns the standard description and comments attributes.
// These are optional text fields that appear on most Netbox resources.

//
// Usage:

//
//	attrs := map[string]schema.Attribute{
//	    "id": IDAttribute("resource"),
//	    "name": NameAttribute("resource", 100),

//	}
//	maps.Copy(attrs, CommonDescriptiveAttributes("resource"))

func CommonDescriptiveAttributes(resourceName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"description": DescriptionAttribute(resourceName),

		"comments": CommentsAttribute(resourceName),
	}
}

// CommonMetadataAttributes returns the standard tags and custom_fields attributes.
// These are optional sets that appear on most Netbox resources for categorization

// and custom data storage.
//

// Usage:
//

//	attrs := map[string]schema.Attribute{
//	    "id": IDAttribute("resource"),
//	    "name": NameAttribute("resource", 100),

//	}
//	maps.Copy(attrs, CommonMetadataAttributes())

func CommonMetadataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"tags": TagsAttribute(),

		"custom_fields": CustomFieldsAttribute(),
	}
}
