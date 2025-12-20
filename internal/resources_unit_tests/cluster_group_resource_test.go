package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestClusterGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil ClusterGroup resource")

	}

}

func TestClusterGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterGroupResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"name", "slug"},

		Computed: []string{"id"},

		Optional: []string{"description", "tags", "custom_fields"},
	})

}

func TestClusterGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterGroupResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_cluster_group")

}

func TestClusterGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterGroupResource()

	testutil.ValidateResourceConfigure(t, r)

}
