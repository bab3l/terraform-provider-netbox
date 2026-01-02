package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceResourceSchema_FieldConsistency(t *testing.T) {
	t.Parallel()

	r := resources.NewDeviceResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	// Test that the device schema now has Optional+Computed fields for consistency
	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required:         []string{"site", "device_type", "role"},
		Optional:         []string{"name", "description", "comments", "tenant", "platform", "serial", "asset_tag", "rack", "position", "face", "latitude", "longitude", "vc_position", "vc_priority"},
		Computed:         []string{"id"},
		OptionalComputed: []string{"status", "airflow", "tags", "custom_fields"},
	})
}

func TestTagsAttribute_OptionalComputed(t *testing.T) {
	t.Parallel()

	attr := nbschema.TagsAttribute()
	if !attr.Optional {
		t.Error("TagsAttribute should be Optional")
	}
	if !attr.Computed {
		t.Error("TagsAttribute should be Computed")
	}
}

func TestCustomFieldsAttribute_OptionalComputed(t *testing.T) {
	t.Parallel()

	attr := nbschema.CustomFieldsAttribute()
	if !attr.Optional {
		t.Error("CustomFieldsAttribute should be Optional")
	}
	if !attr.Computed {
		t.Error("CustomFieldsAttribute should be Computed")
	}
}
