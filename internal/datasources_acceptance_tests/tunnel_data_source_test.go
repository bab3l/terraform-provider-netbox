package datasources_acceptance_tests

import (
	"context"
	"fmt"
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

func TestTunnelDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelDataSource()

	if d == nil {

		t.Fatal("Expected non-nil tunnel data source")

	}

}

func TestTunnelDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check that key attributes exist

	requiredAttrs := []string{"id", "name", "status", "encapsulation", "group", "group_id", "ipsec_profile", "ipsec_profile_id", "tenant", "tenant_id", "tunnel_id", "description", "comments", "tags", "custom_fields"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestTunnelDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTunnelDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelDataSource().(*datasources.TunnelDataSource)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with invalid provider data")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccTunnelDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "status", "active"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
		},
	})

}

func TestAccTunnelDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "status", "active"),
				),
			},
		},
	})

}

func testAccTunnelDataSourceByID(name string) string {

	return fmt.Sprintf(`







resource "netbox_tunnel" "test" {



  name          = %[1]q



  status        = "active"







  encapsulation = "gre"



}







data "netbox_tunnel" "test" {







  id = netbox_tunnel.test.id



}







`, name)

}

func testAccTunnelDataSourceByName(name string) string {

	return fmt.Sprintf(`







resource "netbox_tunnel" "test" {



  name          = %[1]q



  status        = "active"







  encapsulation = "gre"



}







data "netbox_tunnel" "test" {



  name = netbox_tunnel.test.name



}







`, name)

}
