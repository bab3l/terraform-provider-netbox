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

func TestInterfaceResource(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceResource()

	if r == nil {

		t.Fatal("Expected non-nil interface resource")

	}

}

func TestInterfaceResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required attributes

	requiredAttrs := []string{"device", "name", "type"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Optional attributes

	optionalAttrs := []string{

		"label", "enabled", "parent", "bridge", "lag", "mtu",

		"mac_address", "speed", "duplex", "wwn", "mgmt_only",

		"description", "mode", "mark_connected", "tags", "custom_fields",
	}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestInterfaceResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_interface"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestInterfaceResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewInterfaceResource().(*resources.InterfaceResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct provider data type

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with incorrect provider data type

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccInterfaceResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	const interfaceName = "eth0"

	deviceName := testutil.RandomName("tf-test-iface-device")

	manufacturerName := testutil.RandomName("tf-test-iface-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-iface-mfr")

	deviceTypeModel := testutil.RandomName("tf-test-iface-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-iface-dt")

	deviceRoleName := testutil.RandomName("tf-test-iface-role")

	deviceRoleSlug := testutil.RandomSlug("tf-test-iface-role")

	siteName := testutil.RandomName("tf-test-iface-site")

	siteSlug := testutil.RandomSlug("tf-test-iface-site")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

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

			testutil.CheckInterfaceDestroy,

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(

					manufacturerName, manufacturerSlug,

					deviceTypeModel, deviceTypeSlug,

					deviceRoleName, deviceRoleSlug,

					siteName, siteSlug,

					deviceName, interfaceName,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttrPair("netbox_interface.test", "device", "netbox_device.test", "id"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_full(t *testing.T) {

	// Generate unique names

	interfaceName := "GigabitEthernet0/0"

	deviceName := testutil.RandomName("tf-test-iface-dev-full")

	manufacturerName := testutil.RandomName("tf-test-iface-mfr-full")

	manufacturerSlug := testutil.RandomSlug("tf-test-iface-mfr-f")

	deviceTypeModel := testutil.RandomName("tf-test-iface-dt-full")

	deviceTypeSlug := testutil.RandomSlug("tf-test-iface-dt-f")

	deviceRoleName := testutil.RandomName("tf-test-iface-role-full")

	deviceRoleSlug := testutil.RandomSlug("tf-test-iface-role-f")

	siteName := testutil.RandomName("tf-test-iface-site-full")

	siteSlug := testutil.RandomSlug("tf-test-iface-site-f")

	description := "Test interface with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

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

			testutil.CheckInterfaceDestroy,

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_full(

					manufacturerName, manufacturerSlug,

					deviceTypeModel, deviceTypeSlug,

					deviceRoleName, deviceRoleSlug,

					siteName, siteSlug,

					deviceName, interfaceName, description,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttr("netbox_interface.test", "description", description),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "true"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "1500"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mgmt_only", "false"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_update(t *testing.T) {

	// Generate unique names

	const interfaceName = "eth0"

	const updatedInterfaceName = "eth1"

	deviceName := testutil.RandomName("tf-test-iface-dev-upd")

	manufacturerName := testutil.RandomName("tf-test-iface-mfr-upd")

	manufacturerSlug := testutil.RandomSlug("tf-test-iface-mfr-u")

	deviceTypeModel := testutil.RandomName("tf-test-iface-dt-upd")

	deviceTypeSlug := testutil.RandomSlug("tf-test-iface-dt-u")

	deviceRoleName := testutil.RandomName("tf-test-iface-role-upd")

	deviceRoleSlug := testutil.RandomSlug("tf-test-iface-role-u")

	siteName := testutil.RandomName("tf-test-iface-site-upd")

	siteSlug := testutil.RandomSlug("tf-test-iface-site-u")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

	cleanup.RegisterInterfaceCleanup(updatedInterfaceName, deviceName)

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

			testutil.CheckInterfaceDestroy,

			testutil.CheckDeviceDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckDeviceRoleDestroy,

			testutil.CheckManufacturerDestroy,

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(

					manufacturerName, manufacturerSlug,

					deviceTypeModel, deviceTypeSlug,

					deviceRoleName, deviceRoleSlug,

					siteName, siteSlug,

					deviceName, interfaceName,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
				),
			},

			{

				Config: testAccInterfaceResourceConfig_basic(

					manufacturerName, manufacturerSlug,

					deviceTypeModel, deviceTypeSlug,

					deviceRoleName, deviceRoleSlug,

					siteName, siteSlug,

					deviceName, updatedInterfaceName,
				),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_interface.test", "name", updatedInterfaceName),
				),
			},
		},
	})

}

func testAccInterfaceResourceConfig_basic(

	manufacturerName, manufacturerSlug,

	deviceTypeModel, deviceTypeSlug,

	deviceRoleName, deviceRoleSlug,

	siteName, siteSlug,

	deviceName, interfaceName string,

) string {

	return fmt.Sprintf(`































resource "netbox_manufacturer" "test" {































  name = %[1]q































  slug = %[2]q































}































































resource "netbox_device_type" "test" {































  manufacturer = netbox_manufacturer.test.id































  model        = %[3]q































  slug         = %[4]q































}































































resource "netbox_device_role" "test" {































  name = %[5]q































  slug = %[6]q































}































































resource "netbox_site" "test" {































  name = %[7]q































  slug = %[8]q































}































































resource "netbox_device" "test" {































  name        = %[9]q































  device_type = netbox_device_type.test.id































  role        = netbox_device_role.test.id































  site        = netbox_site.test.id































}































































resource "netbox_interface" "test" {































  device = netbox_device.test.id































  name   = %[10]q































  type   = "1000base-t"































}































`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,

		deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, interfaceName)

}

func testAccInterfaceResourceConfig_full(

	manufacturerName, manufacturerSlug,

	deviceTypeModel, deviceTypeSlug,

	deviceRoleName, deviceRoleSlug,

	siteName, siteSlug,

	deviceName, interfaceName, description string,

) string {

	return fmt.Sprintf(`































resource "netbox_manufacturer" "test" {































  name = %[1]q































  slug = %[2]q































}































































resource "netbox_device_type" "test" {































  manufacturer = netbox_manufacturer.test.id































  model        = %[3]q































  slug         = %[4]q































}































































resource "netbox_device_role" "test" {































  name = %[5]q































  slug = %[6]q































}































































resource "netbox_site" "test" {































  name = %[7]q































  slug = %[8]q































}































































resource "netbox_device" "test" {































  name        = %[9]q































  device_type = netbox_device_type.test.id































  role        = netbox_device_role.test.id































  site        = netbox_site.test.id































}































































resource "netbox_interface" "test" {































  device       = netbox_device.test.id































  name         = %[10]q































  type         = "1000base-t"































  description  = %[11]q































  enabled      = true































  mtu          = 1500































  mgmt_only    = false































}































`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,

		deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, interfaceName, description)

}
