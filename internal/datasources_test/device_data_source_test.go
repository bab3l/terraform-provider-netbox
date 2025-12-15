package datasources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDeviceDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewDeviceDataSource()

	if d == nil {

		t.Fatal("Expected non-nil device data source")

	}

}

func TestDeviceDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewDeviceDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check that key attributes exist

	requiredAttrs := []string{"id", "name", "device_type", "role", "site"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewDeviceDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewDeviceDataSource().(*datasources.DeviceDataSource)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccDeviceDataSource_byName(t *testing.T) {

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

				Config: testAccDeviceDataSourceConfig_byName(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("data.netbox_device.test", "device_type", deviceTypeSlug),

					resource.TestCheckResourceAttr("data.netbox_device.test", "role", deviceRoleSlug),

					resource.TestCheckResourceAttr("data.netbox_device.test", "site", siteSlug),

					resource.TestCheckResourceAttr("data.netbox_device.test", "status", "active"),
				),
			},
		},
	})

}

func TestAccDeviceDataSource_bySerial(t *testing.T) {

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

				Config: testAccDeviceDataSourceConfig_bySerial(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("data.netbox_device.test", "serial", serial),
				),
			},
		},
	})

}

// Helper functions to generate test configurations

func testAccDeviceDataSourceConfig_byName(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {

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















data "netbox_device" "test" {







  name = netbox_device.test.name







}







`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)

}

func testAccDeviceDataSourceConfig_bySerial(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial string) string {

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







  serial      = %[10]q







}















data "netbox_device" "test" {







  serial = netbox_device.test.serial







}







`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial)

}
