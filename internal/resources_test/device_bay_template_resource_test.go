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

func TestDeviceBayTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil DeviceBayTemplate resource")

	}

}

func TestDeviceBayTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"device_type", "name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"label", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceBayTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device_bay_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceBayTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayTemplateResource()

	// Type assert to access Configure method

	configurable, ok := r.(fwresource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")

	}

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)

	}

}

// Acceptance Tests.

func TestAccDeviceBayTemplateResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-dbt")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})

}

func TestAccDeviceBayTemplateResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-dbt-full")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", "Test Label"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", "Test device bay template with full options"),

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})

}

func TestAccDeviceBayTemplateResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-dbt-upd")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
				),
			},

			{

				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", "Updated Label"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", "Updated description"),
				),
			},
		},
	})

}

func testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}



resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

func testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}



resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  label       = "Test Label"

  description = "Test device bay template with full options"

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

func testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}



resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  label       = "Updated Label"

  description = "Updated description"

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}

func TestAccDeviceBayTemplateResource_import(t *testing.T) {

	name := testutil.RandomName("tf-test-dbt-import")

	manufacturerName := testutil.RandomName("tf-test-mfr-import")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")

	deviceTypeName := testutil.RandomName("tf-test-dt-import")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-import")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_device_bay_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

// TestAccConsistency_DeviceBayTemplate_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_DeviceBayTemplate_LiteralNames(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	bayName := testutil.RandomName("bay")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", bayName),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "device_type", deviceTypeSlug),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName),
			},
		},
	})

}

func testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.id

  subdevice_role = "parent"

}



resource "netbox_device_bay_template" "test" {

  # Use literal string slug to mimic existing user state

  device_type = %q

  name = %q



  depends_on = [netbox_device_type.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, bayName)

}
