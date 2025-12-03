// Package schema provides examples of using schema attribute factories.
//
// This file demonstrates the before/after comparison of schema definitions.
// It is not used in production - it's for documentation purposes only.
package schema

/*
=============================================================================
BEFORE: Traditional Schema Definition (tenant_resource.go - 125 lines)
=============================================================================

func (r *TenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant in Netbox...",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the tenant (assigned by Netbox).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the tenant. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the tenant. Must be unique...",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant group that this tenant belongs to.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the tenant...",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the tenant...",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this tenant...",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
							},
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
								validators.ValidSlug(),
							},
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				// ... 40+ more lines for custom_fields definition
			},
		},
	}
}

=============================================================================
AFTER: Using Schema Factories (~15 lines)
=============================================================================

import nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"

func (r *TenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant in Netbox. Tenants represent individual customers or organizational units in multi-tenancy scenarios.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.IDAttribute("tenant"),
			"name":          nbschema.NameAttribute("tenant", 100),
			"slug":          nbschema.SlugAttribute("tenant"),
			"group":         nbschema.IDOnlyReferenceAttribute("tenant group", ""),
			"description":   nbschema.DescriptionAttribute("tenant"),
			"comments":      nbschema.CommentsAttributeWithLimit("tenant", 1000),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

=============================================================================
SUMMARY: Code Reduction
=============================================================================

| Resource        | Before (lines) | After (lines) | Reduction |
|-----------------|----------------|---------------|-----------|
| tenant          | ~125           | ~15           | 88%       |
| site            | ~180           | ~25           | 86%       |
| device          | ~230           | ~35           | 85%       |
| rack            | ~280           | ~40           | 86%       |

The factories provide:
- Consistent validation across all resources
- Automatic markdown descriptions
- Reduced chance of copy-paste errors
- Single place to update common patterns

=============================================================================
DATA SOURCE EXAMPLE
=============================================================================

import nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"

func (d *TenantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		MarkdownDescription: "Use this data source to get information about a tenant in Netbox.",

		Attributes: map[string]dsschema.Attribute{
			"id":            nbschema.DSIDAttribute("tenant"),
			"name":          nbschema.DSNameAttribute("tenant"),
			"slug":          nbschema.DSSlugAttribute("tenant"),
			"group":         nbschema.DSComputedStringAttribute("Name of the tenant group."),
			"group_id":      nbschema.DSComputedStringAttribute("ID of the tenant group."),
			"description":   nbschema.DSComputedStringAttribute("Description of the tenant."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments about the tenant."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

*/
