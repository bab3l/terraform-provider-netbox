package datasources_acceptance_tests

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTunnelTerminationDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewTunnelTerminationDataSource()
	if d == nil {
		t.Fatal("Expected non-nil TunnelTermination data source")
	}
}

func TestTunnelTerminationDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewTunnelTerminationDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Lookup attributes
	lookupAttrs := []string{"id", "tunnel", "tunnel_name"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"role", "termination_type", "termination_id", "outside_ip", "tags", "custom_fields"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestTunnelTerminationDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewTunnelTerminationDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel_termination"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestTunnelTerminationDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewTunnelTerminationDataSource()

	configurable, ok := d.(fwdatasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("Data source does not implement DataSourceWithConfigure")
	}

	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

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

func TestAccTunnelTerminationDataSource_byID(t *testing.T) {

	t.Parallel()
	tunnelName := testutil.RandomName("tf-test-tunnel-term-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationDataSourceConfig_byID(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "id", "netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "tunnel", "netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "role", "peer"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationDataSource_byTunnel(t *testing.T) {

	t.Parallel()
	tunnelName := testutil.RandomName("tf-test-tunnel-term-ds2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationDataSourceConfig_byTunnel(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "tunnel", "netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "termination_type", "dcim.device"),
				),
			},
		},
	})
}

func testAccTunnelTerminationDataSourceConfig_byID(tunnelName string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}

data "netbox_tunnel_termination" "test" {
  id = netbox_tunnel_termination.test.id
}
`
}

func testAccTunnelTerminationDataSourceConfig_byTunnel(tunnelName string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}

data "netbox_tunnel_termination" "test" {
  tunnel = netbox_tunnel.test.id
  depends_on = [netbox_tunnel_termination.test]
}
`
}
