package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestClusterResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterResource()

	if r == nil {

		t.Fatal("Expected non-nil Cluster resource")

	}

}

func TestClusterResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterResource()

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

		Required: []string{"name", "type"},

		Computed: []string{"id"},

		Optional: []string{"status", "group", "description", "comments", "tags", "custom_fields"},
	})

}

func TestClusterResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_cluster")

}

func TestClusterResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewClusterResource()

	testutil.ValidateResourceConfigure(t, r)

}
