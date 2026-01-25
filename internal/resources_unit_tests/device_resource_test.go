package resources_unit_tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestDeviceResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestDeviceResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"site", "device_type", "role"},
		Optional: []string{"name", "status", "description", "comments", "tenant", "platform", "cluster", "serial", "asset_tag", "rack", "position", "face", "latitude", "longitude", "config_template", "tags", "custom_fields"},
		Computed: []string{"id"},
	})

	testutil.ValidateFloat64AttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["latitude"],
		"latitude",
		reflect.TypeOf(validators.LatitudeValidator{}),
	)
	testutil.ValidateFloat64AttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["longitude"],
		"longitude",
		reflect.TypeOf(validators.LongitudeValidator{}),
	)
}

func TestDeviceResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_device")
}

func TestDeviceResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()
	testutil.ValidateResourceConfigure(t, r)
}
