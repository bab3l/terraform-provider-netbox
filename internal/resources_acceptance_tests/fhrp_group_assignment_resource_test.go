package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupAssignmentResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-fhrp-assign")
	interfaceName := testutil.RandomName("eth")
	groupID := testutil.RandomFHRPGroupID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")
	cleanup.RegisterFHRPGroupCleanup("vrrp2", groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "200"),
				),
			},
			{
				ResourceName:            "netbox_fhrp_group_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "interface_id", "display_name"},
			},
			{
				Config:             testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName, groupID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName string, groupID int32) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_type" "test" {
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
	group_id = %[3]d
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`, name, interfaceName, groupID)
}

func testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName string, groupID int32) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_type" "test" {
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = %[3]d
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 200
}
`, name, interfaceName, groupID)
}

func TestAccConsistency_FHRPGroupAssignment_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("test-fhrp-assign-lit")
	interfaceName := testutil.RandomName("eth")
	groupID := testutil.RandomFHRPGroupID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfg")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")
	cleanup.RegisterFHRPGroupCleanup("vrrp2", groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
				),
			},
			{
				Config:   testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName, groupID),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
				),
			},
		},
	})
}

func testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName string, groupID int32) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%[1]s-site"
  slug   = "%[1]s-site"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfg"
  slug = "%[1]s-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
}

resource "netbox_device_role" "test" {
  name = "%[1]s-role"
  slug = "%[1]s-role"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = %[3]d
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`, name, interfaceName, groupID)
}

func TestAccFHRPGroupAssignmentResource_full(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("test-fhrp-assign-full")
	interfaceName := testutil.RandomName("eth")
	groupID := testutil.RandomFHRPGroupID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")
	cleanup.RegisterFHRPGroupCleanup("vrrp2", groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
			{
				ResourceName:            "netbox_fhrp_group_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "interface_id", "display_name"},
			},
			{
				Config:   testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName, groupID),
				PlanOnly: true,
			},
		},
	})
}

func TestAccFHRPGroupAssignmentResource_update(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("test-fhrp-assign-upd")
	interfaceName := testutil.RandomName("eth")
	groupID := testutil.RandomFHRPGroupID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")
	cleanup.RegisterFHRPGroupCleanup("vrrp2", groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
				),
			},
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "200"),
				),
			},
		},
	})
}

func TestAccFHRPGroupAssignmentResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-fhrp-assign-extdel")
	interfaceName := testutil.RandomName("eth")
	groupID := testutil.RandomFHRPGroupID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-dt")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")
	cleanup.RegisterFHRPGroupCleanup("vrrp2", groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// List FHRP group assignments to find the one we created
					items, _, err := client.IpamAPI.IpamFhrpGroupAssignmentsList(context.Background()).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find FHRP group assignment for external deletion: %v", err)
					}

					// Find the assignment for this test's group ID to avoid cross-test deletion
					var assignmentID int32
					found := false
					for _, assignment := range items.Results {
						if assignment.Priority == 100 && assignment.Group.GetGroupId() == groupID {
							assignmentID = assignment.Id
							found = true
							break
						}
					}

					if !found {
						t.Fatal("FHRP group assignment not found")
					}

					_, err = client.IpamAPI.IpamFhrpGroupAssignmentsDestroy(context.Background(), assignmentID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete FHRP group assignment: %v", err)
					}
					t.Logf("Successfully externally deleted FHRP group assignment with ID: %d", assignmentID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccFHRPGroupAssignmentResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_fhrp_group_assignment",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_group_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_fhrp_group_assignment" "test" {
  # group_id missing
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_interface_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id = netbox_fhrp_group.test.id
  # interface_type missing
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_interface_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  # interface_id missing
  priority       = 100
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_priority": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model = "test-type"
  slug  = "test-type"
}

resource "netbox_device" "test" {
  name        = "test-device"
  site        = netbox_site.test.id
  device_role = netbox_device_role.test.id
  device_type = netbox_device_type.test.id
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  # priority missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
