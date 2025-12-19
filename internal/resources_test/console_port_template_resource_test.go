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

func TestConsolePortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil console port template resource")
	}
}

func TestConsolePortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

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

func TestConsolePortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_console_port_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestConsolePortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsolePortTemplateResource().(*resources.ConsolePortTemplateResource)

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

// testAccConsolePortTemplateResourceBasic creates a console port template with minimum required fields.

func testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.id
  name        = %q
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

// testAccConsolePortTemplateResourceFull creates a console port template with all fields.

func testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string) string {

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

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.id
  name        = %q
  label       = %q
  type        = %q
  description = %q
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description)
}

func TestAccConsolePortTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "console0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_console_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccConsolePortTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "console0"

	label := "Console 0"

	portType := "rj-45"

	description := "Test console port template"

	updatedName := "console1"

	updatedLabel := "Console 1"

	updatedDescription := "Updated console port template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", description),
				),
			},

			// Update test

			{

				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_ConsolePortTemplate(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("devicetype")

	deviceTypeSlug := testutil.RandomSlug("devicetype")

	const portName = "console0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_device_type" "test" {
  model = "%[3]s"
  slug = "%[4]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_port_template" "test" {

  device_type = netbox_device_type.test.model
  name = "%[5]s"
  type = "rj-45"
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)
}
