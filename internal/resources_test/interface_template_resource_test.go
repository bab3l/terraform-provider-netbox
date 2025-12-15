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

func TestInterfaceTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil interface template resource")

	}

}

func TestInterfaceTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceTemplateResource()

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

	requiredAttrs := []string{"name", "type"}

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

	optionalAttrs := []string{"device_type", "module_type", "label", "enabled", "mgmt_only", "description", "bridge", "poe_mode", "poe_type", "rf_role"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestInterfaceTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_interface_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestInterfaceTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceTemplateResource().(*resources.InterfaceTemplateResource)

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

// testAccInterfaceTemplateResourcePrereqs creates a manufacturer and device type for interface template tests.

func testAccInterfaceTemplateResourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

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







`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug)

}

// testAccInterfaceTemplateResourceBasic creates a basic interface template with required fields only.

func testAccInterfaceTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, interfaceType string) string {

	return testAccInterfaceTemplateResourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug) + fmt.Sprintf(`







resource "netbox_interface_template" "test" {







  device_type = netbox_device_type.test.id







  name        = %q







  type        = %q







}







`, name, interfaceType)

}

// testAccInterfaceTemplateResourceFull creates an interface template with all optional fields.

func testAccInterfaceTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, interfaceType, label, description string) string {

	return testAccInterfaceTemplateResourcePrereqs(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug) + fmt.Sprintf(`







resource "netbox_interface_template" "test" {







  device_type = netbox_device_type.test.id







  name        = %q







  type        = %q







  label       = %q







  enabled     = true







  mgmt_only   = false







  description = %q







}







`, name, interfaceType, label, description)

}

func TestAccInterfaceTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := interfaceName

	interfaceType := "1000base-t"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, interfaceType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", interfaceType),
				),
			},

			{

				ResourceName: "netbox_interface_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccInterfaceTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := interfaceName

	interfaceType := "1000base-t"

	label := "Ethernet 0"

	description := "Test interface template"

	updatedName := "eth1"

	updatedLabel := "Ethernet 1"

	updatedDescription := "Updated interface template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, interfaceType, label, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "type", interfaceType),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "enabled", "true"),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "mgmt_only", "false"),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "description", description),
				),
			},

			{

				Config: testAccInterfaceTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, interfaceType, updatedLabel, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_interface_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_interface_template.test", "description", updatedDescription),
				),
			},
		},
	})

}
