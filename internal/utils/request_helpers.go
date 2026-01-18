// Package utils provides utility functions for working with Netbox provider data structures.
package utils

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// =====================================================
// REQUEST BUILDER INTERFACES
// =====================================================
// These interfaces match the generated go-netbox request types.
// All request types that have these fields implement these interfaces,
// allowing us to write generic helpers that work across all resources.

// DescriptionSetter is implemented by request types that have a Description field.
type DescriptionSetter interface {
	SetDescription(v string)
}

// CommentsSetter is implemented by request types that have a Comments field.
type CommentsSetter interface {
	SetComments(v string)
}

// LabelSetter is implemented by request types that have a Label field.
// Includes both SetLabel for setting values and direct field access for clearing.
type LabelSetter interface {
	SetLabel(v string)
}

// TagsSetter is implemented by request types that have a Tags field.
type TagsSetter interface {
	SetTags(v []netbox.NestedTagRequest)
}

// CustomFieldsSetter is implemented by request types that have a CustomFields field.
type CustomFieldsSetter interface {
	SetCustomFields(v map[string]interface{})
}

// CommonDescriptiveSetter combines Description and Comments setters.
type CommonDescriptiveSetter interface {
	DescriptionSetter
	CommentsSetter
}

// CommonMetadataSetter combines Tags and CustomFields setters.
type CommonMetadataSetter interface {
	TagsSetter
	CustomFieldsSetter
}

// FullCommonFieldsSetter combines all common field setters.
type FullCommonFieldsSetter interface {
	DescriptionSetter
	CommentsSetter
	TagsSetter
	CustomFieldsSetter
}

// =====================================================
// REQUEST BUILDER HELPERS
// =====================================================
// ApplyDescription sets the Description field on a request if the value is set.
// Works with any request type that implements DescriptionSetter.
func ApplyDescription[T DescriptionSetter](request T, description types.String) {
	if IsSet(description) {
		request.SetDescription(description.ValueString())
	} else if description.IsNull() {
		request.SetDescription("")
	}
}

// ApplyComments sets the Comments field on a request if the value is set.
// Works with any request type that implements CommentsSetter.
func ApplyComments[T CommentsSetter](request T, comments types.String) {
	if IsSet(comments) {
		request.SetComments(comments.ValueString())
	} else if comments.IsNull() {
		request.SetComments("")
	}
}

// ApplyLabel sets the Label field on a request if the value is set, or clears it if null.
// Works with any request type that implements LabelSetter.
// Note: We explicitly set empty string to clear the field, as NetBox interprets omitted
// fields as "keep current value" but accepts empty string as "clear the value".
func ApplyLabel[T LabelSetter](request T, label types.String) {
	if IsSet(label) {
		request.SetLabel(label.ValueString())
	} else if label.IsNull() {
		// Set to empty string to clear the field
		// Note: Setting to nil (omit from JSON) causes NetBox to keep the old value
		request.SetLabel("")
	}
}

// ApplyTags sets the Tags field on a request if the value is set.
// Works with any request type that implements TagsSetter.
// Returns any diagnostics from converting the tag models.
func ApplyTags[T TagsSetter](ctx context.Context, request T, tags types.Set, diags *diag.Diagnostics) {
	// If tags is null or empty, explicitly send empty array to clear tags
	if !IsSet(tags) {
		request.SetTags([]netbox.NestedTagRequest{})
		return
	}

	tagRequests, tagDiags := TagModelsToNestedTagRequests(ctx, tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return
	}

	request.SetTags(tagRequests)
}

// ApplyTagsFromSlugs sets the Tags field on a request using tag slugs.
// Looks up tag names by slug to build NestedTagRequests.
func ApplyTagsFromSlugs[T TagsSetter](ctx context.Context, client *netbox.APIClient, request T, tags types.Set, diags *diag.Diagnostics) {
	if !IsSet(tags) {
		request.SetTags([]netbox.NestedTagRequest{})
		return
	}

	tagSlugs := SetToStringSlice(ctx, tags)
	if len(tagSlugs) == 0 {
		request.SetTags([]netbox.NestedTagRequest{})
		return
	}

	tagRequests := TagSlugsToNestedTagRequests(ctx, client, tagSlugs, diags)
	if diags.HasError() {
		return
	}

	request.SetTags(tagRequests)
}

// TagSlugsToNestedTagRequests resolves tag slugs to NestedTagRequests via the API.
func TagSlugsToNestedTagRequests(ctx context.Context, client *netbox.APIClient, tagSlugs []string, diags *diag.Diagnostics) []netbox.NestedTagRequest {
	if len(tagSlugs) == 0 {
		return []netbox.NestedTagRequest{}
	}

	requests := make([]netbox.NestedTagRequest, 0, len(tagSlugs))
	for _, slug := range tagSlugs {
		tags, httpResp, err := client.ExtrasAPI.ExtrasTagsList(ctx).Slug([]string{slug}).Execute()
		CloseResponseBody(httpResp)
		if err != nil {
			diags.AddError("Error reading tag", FormatAPIError("read tag by slug", err, httpResp))
			return nil
		}
		if tags == nil || len(tags.GetResults()) == 0 {
			diags.AddError("Tag Not Found", fmt.Sprintf("No tag found with slug: %s", slug))
			return nil
		}
		tag := tags.GetResults()[0]
		requests = append(requests, netbox.NestedTagRequest{
			Name: tag.GetName(),
			Slug: tag.GetSlug(),
		})
	}

	return requests
}

// ApplyCustomFields sets the CustomFields field on a request if the value is set.
// Works with any request type that implements CustomFieldsSetter.
// Returns any diagnostics from converting the custom field models.
func ApplyCustomFields[T CustomFieldsSetter](ctx context.Context, request T, customFields types.Set, diags *diag.Diagnostics) {
	if !IsSet(customFields) {
		tflog.Debug(ctx, "ApplyCustomFields: customFields is not set (null or unknown), skipping")
		return
	}

	var models []CustomFieldModel
	cfDiags := customFields.ElementsAs(ctx, &models, false)
	diags.Append(cfDiags...)
	if diags.HasError() {
		tflog.Error(ctx, "ApplyCustomFields: Error extracting models", map[string]interface{}{
			"error": cfDiags,
		})
		return
	}

	cfMap := CustomFieldModelsToMap(models)
	tflog.Debug(ctx, "ApplyCustomFields: Setting custom fields", map[string]interface{}{
		"models_count": len(models),
		"models":       models,
		"map":          cfMap,
	})
	request.SetCustomFields(cfMap)
}

// ApplyDescriptiveFields sets both Description and Comments fields on a request.
// Works with any request type that implements CommonDescriptiveSetter.
func ApplyDescriptiveFields[T CommonDescriptiveSetter](request T, description, comments types.String) {
	ApplyDescription(request, description)
	ApplyComments(request, comments)
}

// ApplyMetadataFields sets both Tags and CustomFields on a request.
// Works with any request type that implements CommonMetadataSetter.
func ApplyMetadataFields[T CommonMetadataSetter](ctx context.Context, request T, tags, customFields types.Set, diags *diag.Diagnostics) {
	ApplyTags(ctx, request, tags, diags)
	if diags.HasError() {
		return
	}
	ApplyCustomFields(ctx, request, customFields, diags)
}

// ApplyCommonFields sets Description, Comments, Tags, and CustomFields on a request.
// Works with any request type that implements FullCommonFieldsSetter.
// This is the most common case for resources with full metadata support.
func ApplyCommonFields[T FullCommonFieldsSetter](ctx context.Context, request T, description, comments types.String, tags, customFields types.Set, diags *diag.Diagnostics) {
	ApplyDescriptiveFields(request, description, comments)
	ApplyMetadataFields(ctx, request, tags, customFields, diags)
}

// =====================================================
// CONVENIENCE TYPES FOR COMMON FIELD STRUCTS
// =====================================================

// CommonDescriptiveData holds the common descriptive fields (description + comments).
// Embed this in your resource model or pass to helpers.
type CommonDescriptiveData struct {
	Description types.String
	Comments    types.String
}

// CommonMetadataData holds the common metadata fields (tags + custom_fields).
// Embed this in your resource model or pass to helpers.
type CommonMetadataData struct {
	Tags         types.Set
	CustomFields types.Set
}

// CommonFieldsData holds all common fields.
// Embed this in your resource model or pass to helpers.
type CommonFieldsData struct {
	CommonDescriptiveData
	CommonMetadataData
}

// ApplyTo applies descriptive fields to a request implementing CommonDescriptiveSetter.
func (d CommonDescriptiveData) ApplyTo(request CommonDescriptiveSetter) {
	ApplyDescriptiveFields(request, d.Description, d.Comments)
}

// ApplyTo applies metadata fields to a request implementing CommonMetadataSetter.
func (d CommonMetadataData) ApplyTo(ctx context.Context, request CommonMetadataSetter, diags *diag.Diagnostics) {
	ApplyMetadataFields(ctx, request, d.Tags, d.CustomFields, diags)
}

// ApplyTo applies all common fields to a request implementing FullCommonFieldsSetter.
func (d CommonFieldsData) ApplyTo(ctx context.Context, request FullCommonFieldsSetter, diags *diag.Diagnostics) {
	ApplyCommonFields(ctx, request, d.Description, d.Comments, d.Tags, d.CustomFields, diags)
}

// =====================================================
// MERGE-AWARE REQUEST BUILDERS (FOR UPDATE OPERATIONS)
// =====================================================
// These functions are designed for Update operations where we need to preserve
// unmanaged custom fields that exist in state but aren't in the plan.

// ApplyCustomFieldsWithMerge handles custom fields during Update operations with merge logic.
// This is the merge-aware version of ApplyCustomFields for Update operations.
//
// Behavior:
//   - If plan is null/unknown: Preserve ALL state values (don't manage)
//   - If plan is empty set []: Remove all custom fields (explicit removal)
//   - If plan has values: Merge plan + state, send merged map (partial management)
//
// This enables partial custom field management where users can manage some fields
// in Terraform while preserving others managed externally (NetBox UI, automation, etc.)
//
// Example Update() usage:
//
//	var state, plan DeviceResourceModel
//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
//	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
//	utils.ApplyCustomFieldsWithMerge(ctx, &request, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
func ApplyCustomFieldsWithMerge[T CustomFieldsSetter](
	ctx context.Context,
	request T,
	planCustomFields types.Set,
	stateCustomFields types.Set,
	diags *diag.Diagnostics,
) {
	// User omitted custom_fields (null/unknown) - preserve ALL existing state values
	if !IsSet(planCustomFields) {
		if IsSet(stateCustomFields) {
			tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Preserving state (plan is null)")
			// Reuse existing conversion helper
			var stateModels []CustomFieldModel
			cfDiags := stateCustomFields.ElementsAs(ctx, &stateModels, false)
			diags.Append(cfDiags...)
			if diags.HasError() {
				return
			}
			// Reuse existing helper to convert models to map
			stateMap := CustomFieldModelsToMap(stateModels)
			tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Preserved state custom fields", map[string]interface{}{
				"count": len(stateModels),
			})
			request.SetCustomFields(stateMap)
			return
		}
		// Both null - send empty map (safe default, means "no custom fields")
		tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Both plan and state null, sending empty map")
		request.SetCustomFields(map[string]interface{}{})
		return
	}

	// User specified custom_fields in config (including empty set)
	// Check if plan is an empty set - this means explicit removal of ALL custom fields
	var planModels []CustomFieldModel
	planDiags := planCustomFields.ElementsAs(ctx, &planModels, false)
	diags.Append(planDiags...)
	if diags.HasError() {
		return
	}

	if len(planModels) == 0 {
		// Empty set in plan = explicit request to remove all custom fields
		tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Plan has empty set, removing all custom fields")
		request.SetCustomFields(map[string]interface{}{})
		return
	}

	// Plan has values - merge with existing state
	tflog.Debug(ctx, "ApplyCustomFieldsWithMerge: Merging plan with state")
	merged := MergeCustomFieldSets(ctx, planCustomFields, stateCustomFields, diags)
	if diags.HasError() {
		return
	}
	request.SetCustomFields(merged)
}

// MergeCustomFieldSets merges plan custom fields with state custom fields.
// Uses existing conversion helpers to maintain consistency.
//
// Merge Logic:
//   - Start with ALL state custom fields (preserves unmanaged fields)
//   - Overlay plan custom fields (manages specified fields)
//   - Empty value in plan removes that field from result
//
// Returns a map ready to send to the NetBox API.
func MergeCustomFieldSets(
	ctx context.Context,
	plan types.Set,
	state types.Set,
	diags *diag.Diagnostics,
) map[string]interface{} {
	result := make(map[string]interface{})

	// Start with ALL state custom fields (preserves unmanaged fields)
	if IsSet(state) {
		var stateModels []CustomFieldModel
		stateDiags := state.ElementsAs(ctx, &stateModels, false)
		diags.Append(stateDiags...)
		if diags.HasError() {
			return result
		}

		// Use existing helper to convert state to map
		stateMap := CustomFieldModelsToMap(stateModels)
		for k, v := range stateMap {
			result[k] = v
		}
		tflog.Debug(ctx, "MergeCustomFieldSets: Loaded state fields", map[string]interface{}{
			"count": len(stateModels),
		})
	}

	// Overlay plan custom fields (manages specified fields)
	if IsSet(plan) {
		var planModels []CustomFieldModel
		planDiags := plan.ElementsAs(ctx, &planModels, false)
		diags.Append(planDiags...)
		if diags.HasError() {
			return result
		}

		// Use existing helper to convert plan to map
		planMap := CustomFieldModelsToMap(planModels)
		for k, v := range planMap {
			if v == nil {
				// Empty value in plan removes field from result
				delete(result, k)
				tflog.Debug(ctx, "MergeCustomFieldSets: Removing field (empty value in plan)", map[string]interface{}{
					"field": k,
				})
			} else {
				// Non-empty value in plan overwrites state
				result[k] = v
				tflog.Debug(ctx, "MergeCustomFieldSets: Updating field from plan", map[string]interface{}{
					"field": k,
					"value": v,
				})
			}
		}
	}

	tflog.Debug(ctx, "MergeCustomFieldSets: Merge complete", map[string]interface{}{
		"result_count": len(result),
	})
	return result
}

// ApplyCommonFieldsWithMerge is the merge-aware version of ApplyCommonFields for Update operations.
// Use this in Update() methods where you need to preserve unmanaged custom fields.
//
// Example Update() usage:
//
//	var state, plan DeviceResourceModel
//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
//	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
//	utils.ApplyCommonFieldsWithMerge(ctx, &request,
//	    plan.Description, plan.Comments,
//	    plan.Tags, state.Tags,
//	    plan.CustomFields, state.CustomFields,
//	    &resp.Diagnostics)
func ApplyCommonFieldsWithMerge[T FullCommonFieldsSetter](
	ctx context.Context,
	request T,
	description, comments types.String,
	planTags, stateTags types.Set,
	planCustomFields, stateCustomFields types.Set,
	diags *diag.Diagnostics,
) {
	// Apply simple fields (description, comments) - these don't need merge logic
	ApplyDescriptiveFields(request, description, comments)

	// Apply tags - always use plan. When tags are null/empty in plan, we clear them.
	// Unlike custom fields, tags don't have merge semantics - the plan is the source of truth.
	ApplyTags(ctx, request, planTags, diags)
	if diags.HasError() {
		return
	}

	// Apply custom fields with merge logic
	ApplyCustomFieldsWithMerge(ctx, request, planCustomFields, stateCustomFields, diags)
}
