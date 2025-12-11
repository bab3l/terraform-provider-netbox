package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestExportTemplateResource(t *testing.T) {
	t.Parallel()

	r := resources.NewExportTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil Export Template resource")
	}
}

func TestExportTemplateResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewExportTemplateResource()
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
	requiredAttrs := []string{"name", "object_types", "template_code"}
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
	optionalAttrs := []string{"description", "mime_type", "file_extension", "as_attachment"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestExportTemplateResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewExportTemplateResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_export_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestExportTemplateResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewExportTemplateResource().(*resources.ExportTemplateResource)

	// Test with nil provider data
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with nil provider data should not error: %+v", configureResponse.Diagnostics)
	}

	// Test with valid API client
	configureRequest = fwresource.ConfigureRequest{
		ProviderData: netbox.NewAPIClient(netbox.NewConfiguration()),
	}
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with valid provider data should not error: %+v", configureResponse.Diagnostics)
	}
}

func TestAccExportTemplateResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-export-template")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "object_types.#", "1"),
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
				),
			},
			// Test update
			{
				Config: testAccExportTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", "Updated description"),
				),
			},
			// Test import
			{
				ResourceName:      "netbox_export_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExportTemplateResource_full(t *testing.T) {
	name := acctest.RandomWithPrefix("test-export-template")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "object_types.#", "2"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "mime_type", "text/csv"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "file_extension", "csv"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "as_attachment", "true"),
				),
			},
		},
	})
}

func testAccExportTemplateResourceConfig_basic(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}
`
}

func testAccExportTemplateResourceConfig_updated(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `-updated"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }},{{ site.slug }}\n{% endfor %}"
  description   = "Updated description"
}
`
}

func testAccExportTemplateResourceConfig_full(name string) string {
	return `
resource "netbox_export_template" "test" {
  name           = "` + name + `"
  object_types   = ["dcim.site", "dcim.device"]
  template_code  = "name,slug\n{% for obj in queryset %}{{ obj.name }},{{ obj.slug }}\n{% endfor %}"
  description    = "Test description"
  mime_type      = "text/csv"
  file_extension = "csv"
  as_attachment  = true
}
`
}
