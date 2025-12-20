package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestConfigContextResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigContextResource()

	if r == nil {

		t.Fatal("Expected non-nil resource")

	}

}

func TestConfigContextResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigContextResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"data", "name"},

		Optional: []string{

			"cluster_groups", "cluster_types", "clusters",

			"description", "device_types", "is_active",

			"locations", "platforms", "regions",

			"roles", "site_groups", "sites",

			"tags", "tenant_groups", "tenants", "weight",
		},

		Computed: []string{"id"},
	})

}

func TestConfigContextResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigContextResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_config_context")

}

func TestConfigContextResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigContextResource()

	testutil.ValidateResourceConfigure(t, r)

}
