package resources_test

import (
	"context"
	"fmt"
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

func TestServiceTemplateResource(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil Service Template resource")
	}
}

func TestServiceTemplateResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceTemplateResource()
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
	requiredAttrs := []string{"name", "ports"}
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
	optionalAttrs := []string{"protocol", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestServiceTemplateResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceTemplateResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_service_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestServiceTemplateResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewServiceTemplateResource().(*resources.ServiceTemplateResource)

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

func TestAccServiceTemplateResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-service-template")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.0", "80"),
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
				),
			},
			// Test update
			{
				Config: testAccServiceTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated description"),
				),
			},
			// Test import
			{
				ResourceName:      "netbox_service_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_full(t *testing.T) {
	name := acctest.RandomWithPrefix("test-service-template")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "3"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}
`, name)
}

func testAccServiceTemplateResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "udp"
  ports       = [53, 123]
  description = "Updated description"
}
`, name+"-updated")
}

func testAccServiceTemplateResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "tcp"
  ports       = [80, 443, 8080]
  description = "Test description"
  comments    = "Test comments"
}
`, name)
}
