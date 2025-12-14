package datasources_test

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

func TestTunnelGroupDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelGroupDataSource()

	if d == nil {

		t.Fatal("Expected non-nil tunnel group data source")

	}

}

func TestTunnelGroupDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelGroupDataSource()

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

	requiredAttrs := []string{"id", "name", "slug", "description", "tags", "custom_fields"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestTunnelGroupDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelGroupDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTunnelGroupDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewTunnelGroupDataSource().(*datasources.TunnelGroupDataSource)

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

func TestAccTunnelGroupDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceByID(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func TestAccTunnelGroupDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceByName(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func TestAccTunnelGroupDataSource_bySlug(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceBySlug(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func testAccTunnelGroupDataSourceByID(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_tunnel_group" "test" {



  name = %[1]q



  slug = %[2]q



}







data "netbox_tunnel_group" "test" {



  id = netbox_tunnel_group.test.id



}



`, name, slug)

}

func testAccTunnelGroupDataSourceByName(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_tunnel_group" "test" {



  name = %[1]q



  slug = %[2]q



}







data "netbox_tunnel_group" "test" {



  name = netbox_tunnel_group.test.name



}



`, name, slug)

}

func testAccTunnelGroupDataSourceBySlug(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_tunnel_group" "test" {



  name = %[1]q



  slug = %[2]q



}







data "netbox_tunnel_group" "test" {



  slug = netbox_tunnel_group.test.slug



}



`, name, slug)

}
