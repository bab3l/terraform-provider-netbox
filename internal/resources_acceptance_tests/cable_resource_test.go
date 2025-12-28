package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Cable resource is simpler than other resources (e.g., ASN, Aggregate) because:

// - Only a_terminations and b_terminations are required (complex nested objects)

// - Other fields are simple scalars (type, status, color, label, description, comments, etc.)

// - No complex reference validation or state drift issues

// Therefore, a single comprehensive test that validates core functionality (creating a cable

// with terminations and import/export) is sufficient to ensure the resource works correctly.

func TestAccCableResource_basic(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("test-site-cable")

	siteSlug := testutil.GenerateSlug(siteName)

	deviceName := testutil.RandomName("test-device-cable")

	mfgName := testutil.RandomName("tf-test-mfg-cable")

	mfgSlug := testutil.GenerateSlug(mfgName)

	deviceRoleName := testutil.RandomName("tf-test-role-cable")

	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)

	deviceTypeModel := testutil.RandomName("tf-test-type-cable")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	interfaceNameA := testutil.RandomName("eth")

	interfaceNameB := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName + "-a")
	cleanup.RegisterDeviceCleanup(deviceName + "-b")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),

					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},

			{

				ResourceName: "netbox_cable.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCableResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-id")
	mfgName := testutil.RandomName("tf-test-mfg-cable-id")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable-id")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	interfaceNameA := testutil.RandomName("eth")
	interfaceNameB := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName + "-a-id")
	cleanup.RegisterDeviceCleanup(deviceName + "-b-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName+"-id", mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},
		},
	})

}

func testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %[1]q

  slug = %[2]q

  status = "active"

}

resource "netbox_manufacturer" "test" {

  name = %[3]q

  slug = %[4]q

}

resource "netbox_device_role" "test" {

  name = %[5]q

  slug = %[6]q

}

resource "netbox_device_type" "test" {

  model = %[7]q

  slug  = %[8]q

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_device" "test_a" {

  name           = "%[9]s-a"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_device" "test_b" {

  name           = "%[9]s-b"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_interface" "test_a" {

  name      = %[10]q

  device    = netbox_device.test_a.id

  type      = "1000base-t"

}

resource "netbox_interface" "test_b" {

  name      = %[11]q

  device    = netbox_device.test_b.id

  type      = "1000base-t"

}

resource "netbox_cable" "test" {

  status = "connected"

  type   = "cat6"

  a_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_a.id

    }

  ]

  b_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_b.id

    }

  ]

}

`, siteName, siteSlug, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB)

}

func TestAccConsistency_Cable_LiteralNames(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-lit")

	siteSlug := testutil.GenerateSlug(siteName)

	deviceName := testutil.RandomName("test-device-cable-lit")

	mfgName := testutil.RandomName("tf-test-mfg-cable-lit")

	mfgSlug := testutil.GenerateSlug(mfgName)

	deviceRoleName := testutil.RandomName("tf-test-role-cable-lit")

	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)

	deviceTypeModel := testutil.RandomName("tf-test-type-cable-lit")

	deviceTypeSlug := testutil.RandomSlug("device-type-lit")

	interfaceNameA := testutil.RandomName("eth")

	interfaceNameB := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName + "-a")
	cleanup.RegisterDeviceCleanup(deviceName + "-b")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCableResourceConfigLiteralNames(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),

					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},

			{

				PlanOnly: true,

				Config: testAccCableResourceConfigLiteralNames(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
			},
		},
	})

}

func testAccCableResourceConfigLiteralNames(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %[1]q

  slug = %[2]q

  status = "active"

}

resource "netbox_manufacturer" "test" {

  name = %[3]q

  slug = %[4]q

}

resource "netbox_device_role" "test" {

  name = %[5]q

  slug = %[6]q

}

resource "netbox_device_type" "test" {

  model = %[7]q

  slug  = %[8]q

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_device" "test_a" {

  name           = "%[9]s-a"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_device" "test_b" {

  name           = "%[9]s-b"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_interface" "test_a" {

  name      = %[10]q

  device    = netbox_device.test_a.name

  type      = "1000base-t"

}

resource "netbox_interface" "test_b" {

  name      = %[11]q

  device    = netbox_device.test_b.name

  type      = "1000base-t"

}

resource "netbox_cable" "test" {

  status = "connected"

  type   = "cat6"

  a_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_a.id

    }

  ]

  b_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_b.id

    }

  ]

}

`, siteName, siteSlug, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB)

}

func TestAccCableResource_externalDeletion(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("test-site-cable")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable")
	mfgName := testutil.RandomName("tf-test-mfg-cable")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	interfaceNameA := testutil.RandomName("eth")
	interfaceNameB := testutil.RandomName("eth")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List cables and delete the first one (we just created it)
					items, _, err := client.DcimAPI.DcimCablesList(context.Background()).Limit(10).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to list cables for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimCablesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cable: %v", err)
					}
					t.Logf("Successfully externally deleted cable with ID: %d", itemID)
				},
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
				),
			},
		},
	})
}
