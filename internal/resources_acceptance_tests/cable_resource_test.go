package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// NOTE: Custom field tests for cable resource are in resources_acceptance_tests_customfields package

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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
				ResourceName:      "netbox_cable.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCableResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-full")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-full")
	mfgName := testutil.RandomName("tf-test-mfg-cable-full")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable-full")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable-full")
	deviceTypeSlug := testutil.RandomSlug("device-type-full")
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig_full(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
					resource.TestCheckResourceAttr("netbox_cable.test", "description", testutil.Description1),
					resource.TestCheckResourceAttr("netbox_cable.test", "comments", testutil.Comments),
				),
			},
		},
	})
}

func TestAccCableResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-update")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-update")
	mfgName := testutil.RandomName("tf-test-mfg-cable-update")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable-update")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable-update")
	deviceTypeSlug := testutil.RandomSlug("device-type-update")
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
				Config: testAccCableResourceConfig_withDescription(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
					resource.TestCheckResourceAttr("netbox_cable.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccCableResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-imp")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-imp")
	mfgName := testutil.RandomName("tf-test-mfg-cable-imp")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable-imp")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable-imp")
	deviceTypeSlug := testutil.RandomSlug("device-type-imp")
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},
			{
				ResourceName:      "netbox_cable.test",
				ImportState:       true,
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

func testAccCableResourceConfig_full(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB string) string {
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
  status      = "connected"
  type        = "cat6"
  description = %q
  comments    = %q
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
`, siteName, siteSlug, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB, testutil.Description1, testutil.Comments)
}

func testAccCableResourceConfig_withDescription(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB string) string {
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
  status      = "connected"
  type        = "cat6"
  description = %q
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
`, siteName, siteSlug, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB, testutil.Description2)
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
				Config:   testAccCableResourceConfigLiteralNames(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName + "-a")
	cleanup.RegisterDeviceCleanup(deviceName + "-b")

	var createdCableID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
					func(s *terraform.State) error {
						r, ok := s.RootModule().Resources["netbox_cable.test"]
						if !ok {
							return fmt.Errorf("resource netbox_cable.test not found in state")
						}
						if r.Primary.ID == "" {
							return fmt.Errorf("netbox_cable.test has empty ID")
						}
						createdCableID = r.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					if createdCableID == "" {
						t.Fatalf("created cable ID was not captured from state")
					}
					cableID64, err := strconv.ParseInt(createdCableID, 10, 32)
					if err != nil {
						t.Fatalf("failed to parse created cable ID %q: %v", createdCableID, err)
					}
					cableID := int32(cableID64)
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					_, err = client.DcimAPI.DcimCablesDestroy(context.Background(), cableID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cable: %v", err)
					}
					t.Logf("Successfully externally deleted cable with ID: %d", cableID)
				},
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
				),
			},
		},
	})
}

func TestAccCableResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-cable-opt")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-opt")
	mfgName := testutil.RandomName("tf-test-mfg-cable-opt")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceRoleName := testutil.RandomName("tf-test-role-cable-opt")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeModel := testutil.RandomName("tf-test-type-cable-opt")
	deviceTypeSlug := testutil.RandomSlug("device-type-opt")
	interfaceNameA := testutil.RandomName("eth")
	interfaceNameB := testutil.RandomName("eth")
	label := "Test Label"
	length := 10.5
	lengthUnit := "m"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName + "-a")
	cleanup.RegisterDeviceCleanup(deviceName + "-b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model = %[8]q
  slug  = %[9]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name = "%[3]s-a"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name = "%[3]s-b"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  device = netbox_device.test_a.id
  name   = %[10]q
  type   = "1000base-t"
}

resource "netbox_interface" "test_b" {
  device = netbox_device.test_b.id
  name   = %[11]q
  type   = "1000base-t"
}

resource "netbox_cable" "test" {
  status      = "planned"
  type        = "cat6"
  description = "Description"
  comments    = "Comments"
  label       = %[12]q
  length      = %[13]f
  length_unit = %[14]q
  color       = "ff0000"
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
`, siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, label, length, lengthUnit),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
					resource.TestCheckResourceAttr("netbox_cable.test", "label", label),
					resource.TestCheckResourceAttr("netbox_cable.test", "length", fmt.Sprintf("%g", length)),
					resource.TestCheckResourceAttr("netbox_cable.test", "length_unit", lengthUnit),
					resource.TestCheckResourceAttr("netbox_cable.test", "color", "ff0000"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model = %[8]q
  slug  = %[9]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name = "%[3]s-a"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name = "%[3]s-b"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  device = netbox_device.test_a.id
  name   = %[10]q
  type   = "1000base-t"
}

resource "netbox_interface" "test_b" {
  device = netbox_device.test_b.id
  name   = %[11]q
  type   = "1000base-t"
}

resource "netbox_cable" "test" {
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
`, siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"), // Default
					resource.TestCheckNoResourceAttr("netbox_cable.test", "type"),
					resource.TestCheckNoResourceAttr("netbox_cable.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_cable.test", "length"),
					resource.TestCheckNoResourceAttr("netbox_cable.test", "length_unit"),
					resource.TestCheckNoResourceAttr("netbox_cable.test", "color"),
				),
			},
		},
	})
}

func TestAccCableResource_removeDescriptionCommentsLabel(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-cable-optional")
	siteSlug := testutil.RandomSlug("tf-test-site-cable")
	deviceName := testutil.RandomName("tf-test-device-cable")
	mfgName := testutil.RandomName("tf-test-manufacturer-cable")
	mfgSlug := testutil.RandomSlug("tf-test-mfr-cable")
	deviceRoleName := testutil.RandomName("tf-test-role-cable")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-cable")
	deviceTypeModel := testutil.RandomName("tf-test-devtype-cable")
	deviceTypeSlug := testutil.RandomSlug("tf-test-devtype-cable")
	interfaceNameA := testutil.RandomName("tf-test-iface-a")
	interfaceNameB := testutil.RandomName("tf-test-iface-b")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_cable",
		BaseConfig: func() string {
			return testAccCableResourceConfig(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB)
		},
		ConfigWithFields: func() string {
			return testAccCableResourceConfig_withDescriptionCommentsLabel(
				siteName,
				siteSlug,
				deviceName,
				mfgName,
				mfgSlug,
				deviceRoleName,
				deviceRoleSlug,
				deviceTypeModel,
				deviceTypeSlug,
				interfaceNameA,
				interfaceNameB,
				"Test description",
				"Test comments",
				"Test label",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
			"label":       "Test label",
		},
		RequiredFields: map[string]string{
			"status": "connected",
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckDeviceTypeDestroy,
		),
	})
}

func testAccCableResourceConfig_withDescriptionCommentsLabel(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, description, comments, label string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_type" "test" {
  model        = %[8]q
  slug         = %[9]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "aa1409"
}

resource "netbox_device" "test_a" {
  name        = "%[3]s-a"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_device" "test_b" {
  name        = "%[3]s-b"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
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
  status      = "connected"
  type        = "cat6"
  description = %[12]q
  comments    = %[13]q
  label       = %[14]q
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
`, siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, description, comments, label)
}

func TestAccCableResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_cable",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_a_terminations": {
				Config: func() string {
					return `
resource "netbox_cable" "test" {
  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = 1
    }
  ]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_b_terminations": {
				Config: func() string {
					return `
resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = 1
    }
  ]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func TestAccCableResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-tag")
	siteSlug := testutil.RandomSlug("site-tag")
	deviceName := testutil.RandomName("device-tag")
	mfgName := testutil.RandomName("manufacturer-tag")
	mfgSlug := testutil.RandomSlug("manufacturer-tag")
	deviceRoleName := testutil.RandomName("device-role-tag")
	deviceRoleSlug := testutil.RandomSlug("device-role-tag")
	deviceTypeModel := testutil.RandomName("device-type-tag")
	deviceTypeSlug := testutil.RandomSlug("device-type-tag")
	interfaceNameA := testutil.RandomName("interface-a-tag")
	interfaceNameB := testutil.RandomName("interface-b-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_cable",
		ConfigWithoutTags: func() string {
			return testAccCableResourceConfig_tagLifecycle(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccCableResourceConfig_tagLifecycle(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccCableResourceConfig_tagLifecycle(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckCableDestroy,
	})
}

func TestAccCableResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-tagord")
	siteSlug := testutil.RandomSlug("site-tagord")
	deviceName := testutil.RandomName("device-tagord")
	mfgName := testutil.RandomName("manufacturer-tagord")
	mfgSlug := testutil.RandomSlug("manufacturer-tagord")
	deviceRoleName := testutil.RandomName("device-role-tagord")
	deviceRoleSlug := testutil.RandomSlug("device-role-tagord")
	deviceTypeModel := testutil.RandomName("device-type-tagord")
	deviceTypeSlug := testutil.RandomSlug("device-type-tagord")
	interfaceNameA := testutil.RandomName("interface-a-tagord")
	interfaceNameB := testutil.RandomName("interface-b-tagord")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_cable",
		ConfigWithTagsOrderA: func() string {
			return testAccCableResourceConfig_tagOrder(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccCableResourceConfig_tagOrder(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, false)
		},
		ExpectedTagCount: 3,
		CheckDestroy:     testutil.CheckCableDestroy,
	})
}

// Configuration functions for tag tests.

func testAccCableResourceConfig_tagLifecycle(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[8]q
  slug         = %[9]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name         = %[3]q
  device_type  = netbox_device_type.test.id
  role         = netbox_device_role.test.id
  site         = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = %[10]q
  type   = "1000base-t"
  device = netbox_device.test.id
}

resource "netbox_interface" "test_b" {
  name   = %[11]q
  type   = "1000base-t"
  device = netbox_device.test.id
}

resource "netbox_tag" "tag1" {
  name = %[12]q
  slug = %[13]q
}

resource "netbox_tag" "tag2" {
  name = %[14]q
  slug = %[15]q
}

resource "netbox_tag" "tag3" {
  name = %[16]q
  slug = %[17]q
}
`, siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + `
resource "netbox_cable" "test" {
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
  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`
	case "tag2_tag3":
		return baseConfig + `
resource "netbox_cable" "test" {
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
  tags = [
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    }
  ]
}
`
	default: // "none"
		return baseConfig + `
resource "netbox_cable" "test" {
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
`
	}
}

func testAccCableResourceConfig_tagOrder(siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug string, tag123Order bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %[8]q
  slug         = %[9]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name         = %[3]q
  device_type  = netbox_device_type.test.id
  role         = netbox_device_role.test.id
  site         = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = %[10]q
  type   = "1000base-t"
  device = netbox_device.test.id
}

resource "netbox_interface" "test_b" {
  name   = %[11]q
  type   = "1000base-t"
  device = netbox_device.test.id
}

resource "netbox_tag" "tag1" {
  name = %[12]q
  slug = %[13]q
}

resource "netbox_tag" "tag2" {
  name = %[14]q
  slug = %[15]q
}

resource "netbox_tag" "tag3" {
  name = %[16]q
  slug = %[17]q
}
`, siteName, siteSlug, deviceName, mfgName, mfgSlug, deviceRoleName, deviceRoleSlug, deviceTypeModel, deviceTypeSlug, interfaceNameA, interfaceNameB, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	if tag123Order {
		return baseConfig + `
resource "netbox_cable" "test" {
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
  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    }
  ]
}
`
	}

	return baseConfig + `
resource "netbox_cable" "test" {
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
  tags = [
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    }
  ]
}
`
}
