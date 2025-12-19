package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConsoleServerPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil console server port template resource")

	}

}

func TestConsoleServerPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

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

	requiredAttrs := []string{"name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Check computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Check optional attributes

	optionalAttrs := []string{"device_type", "module_type", "label", "type", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestConsoleServerPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_console_server_port_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestConsoleServerPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortTemplateResource().(*resources.ConsoleServerPortTemplateResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

}

// Acceptance test configurations

// testAccConsoleServerPortTemplateResourceBasic creates a console server port template with minimum required fields.

func testAccConsoleServerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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



resource "netbox_console_server_port_template" "test" {



  device_type = netbox_device_type.test.id

  name        = %q

}



`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

// testAccConsoleServerPortTemplateResourceFull creates a console server port template with all fields.

func testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string) string {

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



resource "netbox_console_server_port_template" "test" {



  device_type = netbox_device_type.test.id

  name        = %q

  label       = %q

  type        = %q

  description = %q

}



`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description)

}

func TestAccConsoleServerPortTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "consoleserver0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_console_server_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccConsoleServerPortTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "consoleserver0"

	label := "Console Server 0"

	portType := "rj-45"

	description := "Test console server port template"

	updatedName := "consoleserver1"

	updatedLabel := "Console Server 1"

	updatedDescription := "Updated console server port template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "description", description),
				),
			},

			// Update test

			{

				Config: testAccConsoleServerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_console_server_port_template.test", "description", updatedDescription),
				),
			},
		},
	})

}
