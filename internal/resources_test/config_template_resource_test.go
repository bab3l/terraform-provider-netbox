package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfigTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil config template resource")

	}

}

func TestConfigTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check required attributes

	requiredAttrs := []string{"name", "template_code"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Check computed attributes

	computedAttrs := []string{"id", "data_path"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Check optional attributes

	optionalAttrs := []string{"description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestConfigTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_config_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestConfigTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConfigTemplateResource().(*resources.ConfigTemplateResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct provider data

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

}

// testAccConfigTemplateResourceBasic creates a basic config template with required fields only.

func testAccConfigTemplateResourceBasic(name, templateCode string) string {

	return fmt.Sprintf(`



resource "netbox_config_template" "test" {



  name          = %q



  template_code = %q



}



`, name, templateCode)

}

// testAccConfigTemplateResourceFull creates a config template with all optional fields.

func testAccConfigTemplateResourceFull(name, templateCode, description string) string {

	return fmt.Sprintf(`



resource "netbox_config_template" "test" {



  name          = %q



  template_code = %q



  description   = %q



}



`, name, templateCode, description)

}

func TestAccConfigTemplateResource_basic(t *testing.T) {

	name := testutil.RandomName("config-tmpl")

	templateCode := "hostname {{ device.name }}"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateResourceBasic(name, templateCode),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
				),
			},

			{

				ResourceName: "netbox_config_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConfigTemplateResource_full(t *testing.T) {

	name := testutil.RandomName("config-tmpl")

	templateCode := "hostname {{ device.name }}"

	description := "Test config template"

	updatedName := testutil.RandomName("config-tmpl-updated")

	updatedTemplateCode := "hostname {{ device.name }}\ninterface {{ interface.name }}"

	updatedDescription := "Updated test config template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateResourceFull(name, templateCode, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),

					resource.TestCheckResourceAttr("netbox_config_template.test", "description", description),
				),
			},

			{

				Config: testAccConfigTemplateResourceFull(updatedName, updatedTemplateCode, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", updatedTemplateCode),

					resource.TestCheckResourceAttr("netbox_config_template.test", "description", updatedDescription),
				),
			},
		},
	})

}
