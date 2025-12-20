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

func TestPowerOutletTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil power outlet template resource")

	}

}

func TestPowerOutletTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()

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

	optionalAttrs := []string{"device_type", "module_type", "label", "type", "power_port", "feed_leg", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestPowerOutletTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_power_outlet_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestPowerOutletTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletTemplateResource().(*resources.PowerOutletTemplateResource)

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

// testAccPowerOutletTemplateResourceBasic creates a power outlet template with minimum required fields.

func testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {

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



resource "netbox_power_outlet_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

// testAccPowerOutletTemplateResourceFull creates a power outlet template with all fields.

func testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description string) string {

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



resource "netbox_power_outlet_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  label       = %q

  type        = %q

  description = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description)

}

func TestAccPowerOutletTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "outlet0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_power_outlet_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccPowerOutletTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := "outlet0"

	label := "Outlet 0"

	outletType := "iec-60320-c13"

	description := "Test power outlet template"

	updatedName := "outlet1"

	updatedLabel := "Outlet 1"

	updatedDescription := "Updated power outlet template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", outletType),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", description),
				),
			},

			// Update test

			{

				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, outletType, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", updatedDescription),
				),
			},
		},
	})

}

// TestAccConsistency_PowerOutletTemplate_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_PowerOutletTemplate_LiteralNames(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	resourceName := testutil.RandomName("power_outlet")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", resourceName),

					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "device_type", deviceTypeSlug),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})

}

func testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model        = %q

  slug         = %q

  manufacturer = netbox_manufacturer.test.id

}



resource "netbox_power_outlet_template" "test" {

  # Use literal string slug to mimic existing user state

  device_type = %q

  name = %q

  type = "iec-60320-c13"



  depends_on = [netbox_device_type.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)

}
