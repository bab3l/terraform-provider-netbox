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

func TestConsoleServerPortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	if r == nil {

		t.Fatal("Expected non-nil ConsoleServerPort resource")

	}

}

func TestConsoleServerPortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"device", "name"}

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

	optionalAttrs := []string{"label", "type", "speed", "description", "mark_connected", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestConsoleServerPortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_console_server_port"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestConsoleServerPortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewConsoleServerPortResource().(*resources.ConsoleServerPortResource)

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

func TestAccConsoleServerPortResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	roleName := testutil.RandomName("tf-test-role")

	roleSlug := testutil.RandomSlug("tf-test-role")

	deviceName := testutil.RandomName("tf-test-device")

	consoleServerPortName := testutil.RandomName("tf-test-csp")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	cleanup.RegisterDeviceRoleCleanup(roleSlug)

	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", consoleServerPortName),
				),
			},

			{

				ResourceName: "netbox_console_server_port.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})

}

func TestAccConsoleServerPortResource_full(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site-full")

	siteSlug := testutil.RandomSlug("tf-test-site-full")

	mfgName := testutil.RandomName("tf-test-mfg-full")

	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")

	dtModel := testutil.RandomName("tf-test-dt-full")

	dtSlug := testutil.RandomSlug("tf-test-dt-full")

	roleName := testutil.RandomName("tf-test-role-full")

	roleSlug := testutil.RandomSlug("tf-test-role-full")

	deviceName := testutil.RandomName("tf-test-device-full")

	consoleServerPortName := testutil.RandomName("tf-test-csp-full")

	description := "Test console server port with all fields"

	updatedDescription := "Updated console server port description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	cleanup.RegisterDeviceRoleCleanup(roleSlug)

	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", consoleServerPortName),

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", description),

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "type", "rj-45"),
				),
			},

			{

				Config: testAccConsoleServerPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccConsoleServerPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName string) string {

	return fmt.Sprintf(`







resource "netbox_site" "test" {



  name   = %q



  slug   = %q



  status = "active"



}







resource "netbox_manufacturer" "test" {



  name = %q



  slug = %q



}







resource "netbox_device_type" "test" {



  manufacturer = netbox_manufacturer.test.id



  model        = %q



  slug         = %q



}







resource "netbox_device_role" "test" {



  name  = %q



  slug  = %q



  color = "aa1409"



}







resource "netbox_device" "test" {



  name        = %q







  device_type = netbox_device_type.test.id



  role        = netbox_device_role.test.id



  site        = netbox_site.test.id



}







resource "netbox_console_server_port" "test" {



  device = netbox_device.test.id



  name   = %q



}







`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName)

}

func testAccConsoleServerPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, description string) string {

	return fmt.Sprintf(`







resource "netbox_site" "test" {



  name   = %q



  slug   = %q



  status = "active"



}







resource "netbox_manufacturer" "test" {



  name = %q



  slug = %q



}







resource "netbox_device_type" "test" {



  manufacturer = netbox_manufacturer.test.id



  model        = %q



  slug         = %q



}







resource "netbox_device_role" "test" {



  name  = %q



  slug  = %q



  color = "aa1409"



}







resource "netbox_device" "test" {



  name        = %q







  device_type = netbox_device_type.test.id



  role        = netbox_device_role.test.id



  site        = netbox_site.test.id



}







resource "netbox_console_server_port" "test" {



  device      = netbox_device.test.id



  name        = %q



  type        = "rj-45"



  description = %q



}







`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, description)

}

func TestAccConsistency_ConsoleServerPort(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	deviceRoleName := testutil.RandomName("device-role")

	deviceRoleSlug := testutil.RandomSlug("device-role")

	deviceName := testutil.RandomName("device")

	portName := testutil.RandomName("console-server-port")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "device", deviceName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccConsoleServerPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, portName),
			},
		},
	})

}

func testAccConsoleServerPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, portName string) string {

	return fmt.Sprintf(`







resource "netbox_site" "test" {



  name = "%[1]s"



  slug = "%[2]s"



}







resource "netbox_manufacturer" "test" {



  name = "%[3]s"



  slug = "%[4]s"



}







resource "netbox_device_type" "test" {



  model = "%[5]s"



  slug = "%[6]s"



  manufacturer = netbox_manufacturer.test.id



}







resource "netbox_device_role" "test" {



  name = "%[7]s"



  slug = "%[8]s"



}







resource "netbox_device" "test" {



  name = "%[9]s"







  device_type = netbox_device_type.test.id



  role = netbox_device_role.test.id



  site = netbox_site.test.id



}







resource "netbox_console_server_port" "test" {



  device = netbox_device.test.name



  name = "%[10]s"



  type = "rj-45"



}







`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, portName)

}
