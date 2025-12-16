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

func TestFrontPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil front port template resource")

	}

}

func TestFrontPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()

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

	requiredAttrs := []string{"name", "type", "rear_port"}

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

	optionalAttrs := []string{"device_type", "module_type", "label", "color", "rear_port_position", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestFrontPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_front_port_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestFrontPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortTemplateResource().(*resources.FrontPortTemplateResource)

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

	// Test with incorrect provider data

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

// Acceptance test configurations

func testAccFrontPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType string) string {

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































resource "netbox_rear_port_template" "test" {















  device_type = netbox_device_type.test.id















  name        = %q















  type        = %q















  positions   = 2















}































resource "netbox_front_port_template" "test" {















  device_type = netbox_device_type.test.id















  name        = %q















  type        = %q















  rear_port   = netbox_rear_port_template.test.name















}















`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, portType, frontPortName, portType)

}

func testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType, label, color, description string, rearPortPosition int32) string {

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































resource "netbox_rear_port_template" "test" {















  device_type = netbox_device_type.test.id















  name        = %q















  type        = %q















  positions   = 4















}































resource "netbox_front_port_template" "test" {















  device_type        = netbox_device_type.test.id















  name               = %q















  type               = %q















  rear_port          = netbox_rear_port_template.test.name















  rear_port_position = %d















  label              = %q















  color              = %q















  description        = %q















}















`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, portType, frontPortName, portType, rearPortPosition, label, color, description)

}

func TestAccFrontPortTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	frontPortName := "front0"

	portType := "8p8c"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccFrontPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_front_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port", rearPortName),
				),
			},

			{

				ResourceName: "netbox_front_port_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccFrontPortTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	frontPortName := "front0"

	portType := "lc"

	label := "Front Port 0"

	description := "Test front port template"

	rearPortPosition := int32(1)

	updatedFrontPortName := "front1"

	updatedLabel := "Front Port 1"

	updatedDescription := "Updated front port template"

	updatedRearPortPosition := int32(2)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType, label, color, description, rearPortPosition),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_front_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "color", color),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", fmt.Sprintf("%d", rearPortPosition)),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", description),
				),
			},

			{

				Config: testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, updatedFrontPortName, portType, updatedLabel, color, updatedDescription, updatedRearPortPosition),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", updatedFrontPortName),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", fmt.Sprintf("%d", updatedRearPortPosition)),

					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", updatedDescription),
				),
			},

			{

				ResourceName: "netbox_front_port_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}
