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

func TestDeviceResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()

	if r == nil {

		t.Fatal("Expected non-nil device resource")

	}

}

func TestDeviceResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"device_type", "role", "site"}

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

	optionalAttrs := []string{"name", "tenant", "platform", "serial", "asset_tag", "location", "rack", "position", "face", "latitude", "longitude", "status", "airflow", "vc_position", "vc_priority", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceResource().(*resources.DeviceResource)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccDeviceResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	deviceName := testutil.RandomName("tf-test-device")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeModel := testutil.RandomName("tf-test-device-type")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	deviceRoleName := testutil.RandomName("tf-test-device-role")

	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceCleanup(deviceName)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),

					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("netbox_device.test", "device_type", deviceTypeSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "role", deviceRoleSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "site", siteSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},

			{

				ResourceName: "netbox_device.test",

				ImportState: true,

				ImportStateVerify: true,

				// Note: some fields may use slugs in config but IDs in state after import

				ImportStateVerifyIgnore: []string{"device_type", "role", "site"},
			},
		},
	})

}

func TestAccDeviceResource_full(t *testing.T) {

	// Generate unique names

	deviceName := testutil.RandomName("tf-test-device")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeModel := testutil.RandomName("tf-test-device-type")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	deviceRoleName := testutil.RandomName("tf-test-device-role")

	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	serial := testutil.RandomName("SN")

	assetTag := testutil.RandomName("AT")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceCleanup(deviceName)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),

					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("netbox_device.test", "device_type", deviceTypeSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "role", deviceRoleSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "site", siteSlug),

					resource.TestCheckResourceAttr("netbox_device.test", "status", "planned"),

					resource.TestCheckResourceAttr("netbox_device.test", "serial", serial),

					resource.TestCheckResourceAttr("netbox_device.test", "asset_tag", assetTag),

					resource.TestCheckResourceAttr("netbox_device.test", "description", "Test device description"),

					resource.TestCheckResourceAttr("netbox_device.test", "comments", "Test device comments"),
				),
			},
		},
	})

}

func TestAccDeviceResource_update(t *testing.T) {

	// Generate unique names

	deviceName := testutil.RandomName("tf-test-device")

	deviceNameUpdated := testutil.RandomName("tf-test-device-updated")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeModel := testutil.RandomName("tf-test-device-type")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	deviceRoleName := testutil.RandomName("tf-test-device-role")

	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceCleanup(deviceName)

	cleanup.RegisterDeviceCleanup(deviceNameUpdated)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},

			{

				Config: testAccDeviceResourceConfig_updated(deviceNameUpdated, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceNameUpdated),

					resource.TestCheckResourceAttr("netbox_device.test", "status", "staged"),

					resource.TestCheckResourceAttr("netbox_device.test", "description", "Updated description"),
				),
			},
		},
	})

}

// Helper functions to generate test configurations

func testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {



  name = %[1]q



  slug = %[2]q



}







resource "netbox_device_type" "test" {



  manufacturer = netbox_manufacturer.test.slug



  model        = %[3]q



  slug         = %[4]q



}







resource "netbox_device_role" "test" {



  name = %[5]q



  slug = %[6]q



}







resource "netbox_site" "test" {



  name   = %[7]q



  slug   = %[8]q



  status = "active"



}







resource "netbox_device" "test" {



  name        = %[9]q



  device_type = netbox_device_type.test.slug



  role        = netbox_device_role.test.slug



  site        = netbox_site.test.slug



}



`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)

}

func testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {



  name = %[1]q



  slug = %[2]q



}







resource "netbox_device_type" "test" {



  manufacturer = netbox_manufacturer.test.slug



  model        = %[3]q



  slug         = %[4]q



}







resource "netbox_device_role" "test" {



  name = %[5]q



  slug = %[6]q



}







resource "netbox_site" "test" {



  name   = %[7]q



  slug   = %[8]q



  status = "active"



}







resource "netbox_device" "test" {



  name        = %[9]q



  device_type = netbox_device_type.test.slug



  role        = netbox_device_role.test.slug



  site        = netbox_site.test.slug



  status      = "planned"



  serial      = %[10]q



  asset_tag   = %[11]q



  description = "Test device description"



  comments    = "Test device comments"



}



`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial, assetTag)

}

func testAccDeviceResourceConfig_updated(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {



  name = %[1]q



  slug = %[2]q



}







resource "netbox_device_type" "test" {



  manufacturer = netbox_manufacturer.test.slug



  model        = %[3]q



  slug         = %[4]q



}







resource "netbox_device_role" "test" {



  name = %[5]q



  slug = %[6]q



}







resource "netbox_site" "test" {



  name   = %[7]q



  slug   = %[8]q



  status = "active"



}







resource "netbox_device" "test" {



  name        = %[9]q



  device_type = netbox_device_type.test.slug



  role        = netbox_device_role.test.slug



  site        = netbox_site.test.slug



  status      = "staged"



  description = "Updated description"



}



`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)

}
