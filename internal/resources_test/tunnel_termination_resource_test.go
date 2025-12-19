package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTunnelTerminationResource(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	if r == nil {
		t.Fatal("Expected non-nil TunnelTermination resource")
	}
}

func TestTunnelTerminationResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes
	requiredAttrs := []string{"tunnel", "termination_type"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Optional attributes
	optionalAttrs := []string{"role", "termination_id", "outside_ip", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestTunnelTerminationResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelTerminationResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel_termination"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestTunnelTerminationResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewTunnelTerminationResource()

	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("Resource does not implement ResourceWithConfigure")
	}

	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

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

func TestAccTunnelTerminationResource_basic(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-for-term")

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
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_update(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-for-term-upd")

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
				Config: testAccTunnelTerminationResourceConfig_withRole(tunnelName, "peer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
				),
			},
			{
				Config: testAccTunnelTerminationResourceConfig_withRole(tunnelName, "hub"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "hub"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_import(t *testing.T) {
	tunnelName := testutil.RandomName("tf-test-tunnel-for-term-imp")

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
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
			},
			{
				ResourceName:            "netbox_tunnel_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tunnel"},
			},
		},
	})
}

func testAccTunnelTerminationResourceConfig_basic(tunnelName string) string {
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
`
}

func testAccTunnelTerminationResourceConfig_withRole(tunnelName, role string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "` + role + `"
}
`
}
