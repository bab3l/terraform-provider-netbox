package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportIdentityWithCustomFieldsSchema returns an identity schema with an ID and optional custom field list.
// Custom field entries should be provided as "name:type" strings.
func ImportIdentityWithCustomFieldsSchema() identityschema.Schema {
	return identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "Numeric ID of the resource to import.",
			},
			"custom_fields": identityschema.ListAttribute{
				ElementType:       types.StringType,
				OptionalForImport: true,
				Description:       "Custom field names to seed during import, formatted as name:type.",
			},
		},
	}
}
