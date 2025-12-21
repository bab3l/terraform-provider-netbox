package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestClusterTypeResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterTypeResource()

	if r == nil {

		t.Fatal("Expected non-nil ClusterType resource")

	}

}

func TestClusterTypeResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterTypeResource()

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

		Optional: []string{"description"},
	})

}

func TestClusterTypeResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterTypeResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_cluster_type")

}

func TestClusterTypeResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterTypeResource()

	testutil.ValidateResourceConfigure(t, r)

}
