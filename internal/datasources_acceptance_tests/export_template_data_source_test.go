package datasources_acceptance_tests

import (
	"context"
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

func TestExportTemplateDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewExportTemplateDataSource()
	if d == nil {
		t.Fatal("Expected non-nil Export Template data source")
	}
}

func TestExportTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewExportTemplateDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes
	requiredAttrs := []string{"id"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"name", "object_types", "template_code", "description", "mime_type", "file_extension", "as_attachment"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestExportTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewExportTemplateDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_export_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestExportTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewExportTemplateDataSource().(*datasources.ExportTemplateDataSource)

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

func TestAccExportTemplateDataSource_byID(t *testing.T) {

	t.Parallel()
	name := acctest.RandomWithPrefix("test-export-tmpl-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "object_types.#", "1"),
					resource.TestCheckResourceAttrSet("data.netbox_export_template.test", "template_code"),
				),
			},
		},
	})
}

func testAccExportTemplateDataSourceConfig_byID(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}

data "netbox_export_template" "test" {
  id = netbox_export_template.test.id
}
`
}
