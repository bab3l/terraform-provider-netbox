package resources_unit_tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestIPAddressResource(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestIPAddressResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()
	schemaRequest := &resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}
	r.Schema(context.Background(), *schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{
		Required: []string{"address"},
		Optional: []string{"vrf", "tenant", "status", "role", "assigned_object_type", "assigned_object_id", "nat_inside", "dns_name", "description", "comments"},
		Computed: []string{"id"},
	})

	testutil.ValidateStringAttributeHasValidatorType(
		t,
		schemaResponse.Schema.Attributes["address"],
		"address",
		reflect.TypeOf(validators.IPAddressWithPrefixValidator{}),
	)
}

func TestIPAddressResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_ip_address")
}

func TestIPAddressResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()
	testutil.ValidateResourceConfigure(t, r)
}
