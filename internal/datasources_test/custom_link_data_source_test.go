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

func TestCustomLinkDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomLinkDataSource()
	if d == nil {
		t.Fatal("Expected non-nil custom link data source")
	}
}

func TestCustomLinkDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomLinkDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	lookupAttrs := []string{"id", "name"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}
}

func TestCustomLinkDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomLinkDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_custom_link"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCustomLinkDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomLinkDataSource()

	// Type assert to access Configure method
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
	configureResponse = &fwdatasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestAccCustomLinkDataSource_byID(t *testing.T) {
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "link_text", "View Details"),
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "object_types.#", "1"),
				),
			},
		},
	})
}

func TestAccCustomLinkDataSource_byName(t *testing.T) {
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkDataSourceConfig_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_custom_link.test", "id"),
				),
			},
		},
	})
}

func testAccCustomLinkDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
}

data "netbox_custom_link" "test" {
  id = netbox_custom_link.test.id
}
`, name)
}

func testAccCustomLinkDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
}

data "netbox_custom_link" "test" {
  name = netbox_custom_link.test.name
}
`, name)
}
