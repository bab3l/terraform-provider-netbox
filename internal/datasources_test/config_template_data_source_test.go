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

func TestConfigTemplateDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewConfigTemplateDataSource()
	if d == nil {
		t.Fatal("Expected non-nil config template data source")
	}
}

func TestConfigTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewConfigTemplateDataSource()
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
	requiredAttrs := []string{"id", "name", "description", "template_code", "data_source", "data_path"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected attribute %s to exist in schema", attr)
		}
	}
}

func TestConfigTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewConfigTemplateDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_config_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestConfigTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewConfigTemplateDataSource().(*datasources.ConfigTemplateDataSource)

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

	configureRequest.ProviderData = "invalid"
	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with invalid provider data")
	}
}

// testAccConfigTemplateDataSourcePrereqs creates prerequisites for config template data source tests
func testAccConfigTemplateDataSourcePrereqs(name, templateCode string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
  name          = %q
  template_code = %q
}
`, name, templateCode)
}

// testAccConfigTemplateDataSourceByID looks up a config template by ID
func testAccConfigTemplateDataSourceByID(name, templateCode string) string {
	return testAccConfigTemplateDataSourcePrereqs(name, templateCode) + `
data "netbox_config_template" "test" {
  id = netbox_config_template.test.id
}
`
}

// testAccConfigTemplateDataSourceByName looks up a config template by name
func testAccConfigTemplateDataSourceByName(name, templateCode string) string {
	return testAccConfigTemplateDataSourcePrereqs(name, templateCode) + fmt.Sprintf(`
data "netbox_config_template" "test" {
  name = %q

  depends_on = [netbox_config_template.test]
}
`, name)
}

func TestAccConfigTemplateDataSource_byID(t *testing.T) {
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("config-tmpl-ds")
	templateCode := "hostname {{ device.name }}"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateDataSourceByID(name, templateCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_config_template.test", "template_code", templateCode),
					resource.TestCheckResourceAttrSet("data.netbox_config_template.test", "id"),
				),
			},
		},
	})
}

func TestAccConfigTemplateDataSource_byName(t *testing.T) {
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("config-tmpl-ds")
	templateCode := "hostname {{ device.name }}"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateDataSourceByName(name, templateCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_config_template.test", "template_code", templateCode),
					resource.TestCheckResourceAttrSet("data.netbox_config_template.test", "id"),
				),
			},
		},
	})
}
