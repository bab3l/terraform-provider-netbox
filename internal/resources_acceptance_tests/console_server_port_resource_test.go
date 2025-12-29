package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortResource_basic(t *testing.T) {

	t.Parallel()

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

	t.Parallel()

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

	description := testutil.RandomName("description")

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

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

func TestAccConsistency_ConsoleServerPort_LiteralNames(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	deviceName := testutil.RandomName("device")

	resourceName := testutil.RandomName("console_server_port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccConsoleServerPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", resourceName),

					resource.TestCheckResourceAttr("netbox_console_server_port.test", "device", deviceName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccConsoleServerPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})

}

func TestAccConsoleServerPortResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-update")
	siteSlug := testutil.RandomSlug("tf-test-site-update")
	mfgName := testutil.RandomName("tf-test-mfg-update")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-update")
	dtModel := testutil.RandomName("tf-test-dt-update")
	dtSlug := testutil.RandomSlug("tf-test-dt-update")
	roleName := testutil.RandomName("tf-test-role-update")
	roleSlug := testutil.RandomSlug("tf-test-role-update")
	deviceName := testutil.RandomName("tf-test-device-update")
	consoleServerPortName := testutil.RandomName("tf-test-csp-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", consoleServerPortName),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccConsoleServerPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccConsoleServerPortResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-ext-del")
	siteSlug := testutil.RandomSlug("tf-test-site-ext-del")
	mfgName := testutil.RandomName("tf-test-mfg-ext-del")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ext-del")
	dtModel := testutil.RandomName("tf-test-dt-ext-del")
	dtSlug := testutil.RandomSlug("tf-test-dt-ext-del")
	roleName := testutil.RandomName("tf-test-role-ext-del")
	roleSlug := testutil.RandomSlug("tf-test-role-ext-del")
	deviceName := testutil.RandomName("tf-test-device-ext-del")
	consoleServerPortName := testutil.RandomName("tf-test-csp-ext-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_server_port.test", "name", consoleServerPortName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimConsoleServerPortsList(context.Background()).NameIc([]string{consoleServerPortName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find console_server_port for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimConsoleServerPortsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete console_server_port: %v", err)
					}
					t.Logf("Successfully externally deleted console_server_port with ID: %d", itemID)
				},
				Config: testAccConsoleServerPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, consoleServerPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_server_port.test", "id"),
				),
			},
		},
	})
}

func TestAccConsoleServerPortResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-id")
	siteSlug := testutil.RandomSlug("tf-test-site-id")
	mfgName := testutil.RandomName("tf-test-mfg-id")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-id")
	dtModel := testutil.RandomName("tf-test-dt-id")
	dtSlug := testutil.RandomSlug("tf-test-dt-id")
	roleName := testutil.RandomName("tf-test-role-id")
	roleSlug := testutil.RandomSlug("tf-test-role-id")
	deviceName := testutil.RandomName("tf-test-device-id")
	consoleServerPortName := testutil.RandomName("tf-test-csp-id")

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

func testAccConsoleServerPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.id

  subdevice_role = "parent"  # Enable device bays

}

resource "netbox_site" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_role" "test" {

  name = %q

  slug = %q

  color = "ff0000"

}

resource "netbox_device" "test" {

  name        = %q

  device_type = netbox_device_type.test.id

  site        = netbox_site.test.id

  role        = netbox_device_role.test.id

  status      = "active"

}

resource "netbox_console_server_port" "test" {

  # Use literal string name to mimic existing user state

  device = %q

  name = %q

  type = "rj-45"

  depends_on = [netbox_device.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, resourceName)

}
