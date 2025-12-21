package datasources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCircuitGroupDataSource(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupDataSource()
	if ds == nil {
		t.Fatal("Expected non-nil CircuitGroup data source")
	}
}

func TestCircuitGroupDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupDataSource()
	schemaRequest := datasource.SchemaRequest{}
	schemaResponse := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check that required lookup attributes exist
	lookupAttrs := []string{"id", "name", "slug"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Check that computed attributes exist
	computedAttrs := []string{"description", "tenant_id", "circuit_count"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestCircuitGroupDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupDataSource()
	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_circuit_group"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCircuitGroupDataSourceConfigure(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupDataSource()

	// Type assert to access Configure method
	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("Data source does not implement DataSourceWithConfigure")
	}

	configureRequest := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &datasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)
	}
}

// Acceptance Tests

func TestAccCircuitGroupDataSource_byID(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-cg-ds-id")
	slug := testutil.RandomSlug("tf-test-cg-ds-id")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_byID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupDataSource_byName(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-cg-ds-name")
	slug := testutil.RandomSlug("tf-test-cg-ds-name")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_byName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupDataSource_bySlug(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-cg-ds-slug")
	slug := testutil.RandomSlug("tf-test-cg-ds-slug")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_bySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccCircuitGroupDataSourceConfig_byID(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  id = netbox_circuit_group.test.id
}
`, name, slug)
}

func testAccCircuitGroupDataSourceConfig_byName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  name = netbox_circuit_group.test.name
}
`, name, slug)
}

func testAccCircuitGroupDataSourceConfig_bySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

data "netbox_circuit_group" "test" {
  slug = netbox_circuit_group.test.slug
}
`, name, slug)
}
