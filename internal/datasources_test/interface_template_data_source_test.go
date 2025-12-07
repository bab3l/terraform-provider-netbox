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

func TestInterfaceTemplateDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewInterfaceTemplateDataSource()
	if d == nil {
		t.Fatal("Expected non-nil interface template data source")
	}
}

func TestInterfaceTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewInterfaceTemplateDataSource()
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
	requiredAttrs := []string{"id", "name", "type", "device_type", "module_type", "label", "enabled", "mgmt_only", "description", "bridge", "poe_mode", "poe_type", "rf_role"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected attribute %s to exist in schema", attr)
		}
	}
}

func TestInterfaceTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewInterfaceTemplateDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_interface_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestInterfaceTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewInterfaceTemplateDataSource().(*datasources.InterfaceTemplateDataSource)

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

// testAccInterfaceTemplateDataSourcePrereqs creates prerequisites for interface template data source tests
func testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_interface_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType)
}

// testAccInterfaceTemplateDataSourceByID looks up an interface template by ID
func testAccInterfaceTemplateDataSourceByID(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {
	return testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType) + `
data "netbox_interface_template" "test" {
  id = netbox_interface_template.test.id
}
`
}

// testAccInterfaceTemplateDataSourceByName looks up an interface template by name and device type
func testAccInterfaceTemplateDataSourceByName(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType string) string {
	return testAccInterfaceTemplateDataSourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType) + fmt.Sprintf(`
data "netbox_interface_template" "test" {
  name        = %q
  device_type = netbox_device_type.test.id

  depends_on = [netbox_interface_template.test]
}
`, templateName)
}

func TestAccInterfaceTemplateDataSource_byID(t *testing.T) {
	testutil.TestAccPreCheck(t)

	manufacturerName := testutil.RandomName("mfr-ds")
	manufacturerSlug := testutil.RandomSlug("mfr-ds")
	deviceTypeName := testutil.RandomName("dt-ds")
	deviceTypeSlug := testutil.RandomSlug("dt-ds")
	templateName := "eth0"
	templateType := "1000base-t"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateDataSourceByID(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "name", templateName),
					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "type", templateType),
					resource.TestCheckResourceAttrSet("data.netbox_interface_template.test", "id"),
				),
			},
		},
	})
}

func TestAccInterfaceTemplateDataSource_byName(t *testing.T) {
	testutil.TestAccPreCheck(t)

	manufacturerName := testutil.RandomName("mfr-ds")
	manufacturerSlug := testutil.RandomSlug("mfr-ds")
	deviceTypeName := testutil.RandomName("dt-ds")
	deviceTypeSlug := testutil.RandomSlug("dt-ds")
	templateName := "eth0"
	templateType := "1000base-t"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceTemplateDataSourceByName(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, templateName, templateType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "name", templateName),
					resource.TestCheckResourceAttr("data.netbox_interface_template.test", "type", templateType),
					resource.TestCheckResourceAttrSet("data.netbox_interface_template.test", "id"),
				),
			},
		},
	})
}
