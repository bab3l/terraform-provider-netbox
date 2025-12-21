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

func TestCircuitGroupAssignmentDataSource(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupAssignmentDataSource()
	if ds == nil {
		t.Fatal("Expected non-nil CircuitGroupAssignment data source")
	}
}

func TestCircuitGroupAssignmentDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupAssignmentDataSource()
	schemaRequest := datasource.SchemaRequest{}
	schemaResponse := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check that lookup attribute exists
	if _, exists := schemaResponse.Schema.Attributes["id"]; !exists {
		t.Error("Expected lookup attribute 'id' to exist in schema")
	}

	// Check that computed attributes exist
	computedAttrs := []string{"group_id", "circuit_id", "priority"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestCircuitGroupAssignmentDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupAssignmentDataSource()
	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_circuit_group_assignment"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCircuitGroupAssignmentDataSourceConfigure(t *testing.T) {
	t.Parallel()

	ds := datasources.NewCircuitGroupAssignmentDataSource()

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

func TestAccCircuitGroupAssignmentDataSource_byID(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	groupName := testutil.RandomName("tf-test-cga-ds-group")
	groupSlug := testutil.RandomSlug("tf-test-cga-ds-grp")
	providerName := testutil.RandomName("tf-test-cga-ds-prov")
	providerSlug := testutil.RandomSlug("tf-test-cga-ds-prov")
	circuitTypeName := testutil.RandomName("tf-test-cga-ds-type")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-ds-type")
	circuitCid := testutil.RandomSlug("tf-test-cga-ds-ckt")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentDataSourceConfig_byID(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "circuit_id"),
				),
			},
		},
	})
}

func TestAccCircuitGroupAssignmentDataSource_withPriority(t *testing.T) {
	// Generate unique names
	groupName := testutil.RandomName("tf-test-cga-ds-grp-pri")
	groupSlug := testutil.RandomSlug("tf-test-cga-ds-grp-pri")
	providerName := testutil.RandomName("tf-test-cga-ds-prov-pri")
	providerSlug := testutil.RandomSlug("tf-test-cga-ds-prov-pri")
	circuitTypeName := testutil.RandomName("tf-test-cga-ds-type-pri")
	circuitTypeSlug := testutil.RandomSlug("tf-test-cga-ds-type-pri")
	circuitCid := testutil.RandomSlug("tf-test-cga-ds-ckt-pri")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupAssignmentCleanup(groupName)
	cleanup.RegisterCircuitGroupCleanup(groupName)
	cleanup.RegisterCircuitCleanup(circuitCid)
	cleanup.RegisterProviderCleanup(providerName)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitGroupAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupAssignmentDataSourceConfig_withPriority(
					groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, "primary",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_circuit_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group_assignment.test", "priority", "primary"),
				),
			},
		},
	})
}

func testAccCircuitGroupAssignmentDataSourceConfig_byID(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_type" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_circuit" "test" {
  cid              = %[7]q
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
}

data "netbox_circuit_group_assignment" "test" {
  id = netbox_circuit_group_assignment.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid)
}

func testAccCircuitGroupAssignmentDataSourceConfig_withPriority(groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_type" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_circuit" "test" {
  cid              = %[7]q
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit_group_assignment" "test" {
  group_id   = netbox_circuit_group.test.id
  circuit_id = netbox_circuit.test.id
  priority   = %[8]q
}

data "netbox_circuit_group_assignment" "test" {
  id = netbox_circuit_group_assignment.test.id
}
`, groupName, groupSlug, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, priority)
}
