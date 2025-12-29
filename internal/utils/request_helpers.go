// Package utils provides utility functions for working with Netbox provider data structures.
package utils

import (
	"context"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	}
}

// ApplyComments sets the Comments field on a request if the value is set.
// Works with any request type that implements CommentsSetter.
func ApplyComments[T CommentsSetter](request T, comments types.String) {
	if IsSet(comments) {
		request.SetComments(comments.ValueString())
	}
}

// ApplyTags sets the Tags field on a request if the value is set.
// Works with any request type that implements TagsSetter.
// Returns any diagnostics from converting the tag models.
func ApplyTags[T TagsSetter](ctx context.Context, request T, tags types.Set, diags *diag.Diagnostics) {
	if !IsSet(tags) {
		return
	}

	tagRequests, tagDiags := TagModelsToNestedTagRequests(ctx, tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return
	}

	request.SetTags(tagRequests)
}

// ApplyCustomFields sets the CustomFields field on a request if the value is set.
// Works with any request type that implements CustomFieldsSetter.
// Returns any diagnostics from converting the custom field models.
func ApplyCustomFields[T CustomFieldsSetter](ctx context.Context, request T, customFields types.Set, diags *diag.Diagnostics) {
	if !IsSet(customFields) {
		return
	}

	var models []CustomFieldModel
	cfDiags := customFields.ElementsAs(ctx, &models, false)
	diags.Append(cfDiags...)
	if diags.HasError() {
		return
	}

	request.SetCustomFields(CustomFieldModelsToMap(models))
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
