package datasources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestL2VPNDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewL2VPNDataSource()
	if d == nil {
		t.Fatal("Expected non-nil L2VPN data source")
	}
}

func TestL2VPNDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewL2VPNDataSource()
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
	lookupAttrs := []string{"id", "name", "slug"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"type", "identifier", "description", "comments"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestL2VPNDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewL2VPNDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_l2vpn"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestL2VPNDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewL2VPNDataSource().(*datasources.L2VPNDataSource)

	// Test with nil provider data
	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct client type
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestAccL2VPNDataSource_byID(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func TestAccL2VPNDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_byName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func testAccL2VPNDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

data "netbox_l2vpn" "test" {
  id = netbox_l2vpn.test.id
}
`, name, name)
}

func testAccL2VPNDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

data "netbox_l2vpn" "test" {
  name = netbox_l2vpn.test.name

  depends_on = [netbox_l2vpn.test]
}
`, name, name)
}
