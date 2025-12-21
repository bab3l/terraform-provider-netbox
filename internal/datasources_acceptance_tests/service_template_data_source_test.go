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

func TestServiceTemplateDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewServiceTemplateDataSource()
	if d == nil {
		t.Fatal("Expected non-nil Service Template data source")
	}
}

func TestServiceTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewServiceTemplateDataSource()
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
	lookupAttrs := []string{"id", "name"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"protocol", "ports", "description", "comments"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestServiceTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewServiceTemplateDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_service_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestServiceTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewServiceTemplateDataSource().(*datasources.ServiceTemplateDataSource)

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

func TestAccServiceTemplateDataSource_byID(t *testing.T) {
	name := acctest.RandomWithPrefix("test-svc-tmpl-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "ports.0", "80"),
				),
			},
		},
	})
}

func TestAccServiceTemplateDataSource_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("test-svc-tmpl-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateDataSourceConfig_byName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccServiceTemplateDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}

data "netbox_service_template" "test" {
  id = netbox_service_template.test.id
}
`, name)
}

func testAccServiceTemplateDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}

data "netbox_service_template" "test" {
  name = netbox_service_template.test.name

  depends_on = [netbox_service_template.test]
}
`, name)
}
