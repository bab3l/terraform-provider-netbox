package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// ReferenceAttributeWithDiffSuppress returns an optional reference attribute with diff suppression.
// This enhanced version suppresses diffs when values refer to the same NetBox object but use
// different formats (e.g., name vs slug vs ID).
func ReferenceAttributeWithDiffSuppress(targetResource string, description string) schema.StringAttribute {
	if description == "" {
		description = "ID or slug of the " + targetResource + "."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,
		Optional:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			ReferenceResolveToIDPlanModifier(targetResource),
			ReferenceEquivalencePlanModifier(),
		},
	}
}

// RequiredReferenceAttributeWithDiffSuppress returns a required reference attribute with diff suppression.
func RequiredReferenceAttributeWithDiffSuppress(targetResource string, description string) schema.StringAttribute {
	if description == "" {
		description = "ID or slug of the " + targetResource + ". Required."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			ReferenceResolveToIDPlanModifier(targetResource),
			ReferenceEquivalencePlanModifier(),
		},
	}
}
