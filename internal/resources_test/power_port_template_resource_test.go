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

func TestPowerPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil power port template resource")

	}

}

func TestPowerPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()

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

	optionalAttrs := []string{"device_type", "module_type", "label", "type", "maximum_draw", "allocated_draw", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestPowerPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_power_port_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestPowerPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerPortTemplateResource().(*resources.PowerPortTemplateResource)

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

// testAccPowerPortTemplateResourceBasic creates a power port template with minimum required fields.

func testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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































resource "netbox_power_port_template" "test" {















  device_type = netbox_device_type.test.id















  name        = %q















}















`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

// testAccPowerPortTemplateResourceFull creates a power port template with all fields.

func testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string, maxDraw, allocDraw int) string {

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































resource "netbox_power_port_template" "test" {















  device_type     = netbox_device_type.test.id















  name            = %q















  label           = %q















  type            = %q















  maximum_draw    = %d















  allocated_draw  = %d















  description     = %q















}















`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, maxDraw, allocDraw, description)

}

func TestAccPowerPortTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "psu0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_power_port_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccPowerPortTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "psu0"

	label := "PSU 0"

	portType := "iec-60320-c14"

	description := "Test power port template"

	maxDraw := 500

	allocDraw := 400

	updatedName := "psu1"

	updatedLabel := "PSU 1"

	updatedDescription := "Updated power port template"

	updatedMaxDraw := 600

	updatedAllocDraw := 450

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description, maxDraw, allocDraw),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "500"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "400"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", description),
				),
			},

			// Update test

			{

				Config: testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription, updatedMaxDraw, updatedAllocDraw),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "600"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "450"),

					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", updatedDescription),
				),
			},
		},
	})

}
