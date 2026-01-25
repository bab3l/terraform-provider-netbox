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

func TestIPRangeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPRangeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required:         []string{"start_address", "end_address"},
		Optional:         []string{"vrf", "tenant", "role", "description", "comments"},
		Computed:         []string{"id", "size"},
		OptionalComputed: []string{"status", "mark_utilized"},
	})

	testutil.ValidateStringAttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["start_address"],
		"start_address",
		reflect.TypeOf(validators.IPAddressValidator{}),
	)
	testutil.ValidateStringAttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["end_address"],
		"end_address",
		reflect.TypeOf(validators.IPAddressValidator{}),
	)
}

func TestIPRangeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ip_range")
}

func TestIPRangeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()
	testutil.ValidateResourceConfigure(t, r)
}
