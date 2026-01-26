package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for virtual device context resource are in resources_acceptance_tests_customfields package

func TestAccVirtualDeviceContextResource_basic(t *testing.T) {
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
	vdcName := testutil.RandomName("tf-test-vdc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "status", "active"),
				),
			},
			{
				ResourceName:            "netbox_virtual_device_context.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_virtual_device_context.test", "device"),
					testutil.ReferenceFieldCheck("netbox_virtual_device_context.test", "tenant"),
				),
			},
			{
				Config:             testAccVirtualDeviceContextResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccVirtualDeviceContextResource_full(t *testing.T) {
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
	vdcName := testutil.RandomName("tf-test-vdc-full")
	tenantName := testutil.RandomName("tf-test-tenant-full")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-full")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("updated-description")
	comments := testutil.RandomName("comments")
	updatedComments := testutil.RandomName("updated-comments")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	ipAddress := testutil.RandomIPv4Address()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_fullWithAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "identifier", "100"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", description),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "comments", comments),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "primary_ip4"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tags.#", "2"),
				),
			},
			{
				Config: testAccVirtualDeviceContextResourceConfig_fullWithAllFieldsUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "identifier", "200"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "comments", updatedComments),
				),
			},
		},
	})

}

func TestAccVirtualDeviceContextResource_update(t *testing.T) {
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
	vdcName := testutil.RandomName("tf-test-vdc")
	description1 := "Initial description"
	description2 := "Updated description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", description1),
				),
			},
			{
				Config: testAccVirtualDeviceContextResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", description2),
				),
			},
		},
	})
}

func TestAccConsistency_VirtualDeviceContext(t *testing.T) {
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
	vdcName := testutil.RandomName("vdc")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_virtual_device_context.test", "device"),
					testutil.ReferenceFieldCheck("netbox_virtual_device_context.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVirtualDeviceContextConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_VirtualDeviceContext_LiteralNames(t *testing.T) {
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
	vdcName := testutil.RandomName("vdc")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextConsistencyLiteralNamesConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "device", deviceName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVirtualDeviceContextConsistencyLiteralNamesConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccVirtualDeviceContextResource_IDPreservation(t *testing.T) {
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
	vdcName := testutil.RandomName("tf-test-vdc-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
				),
			},
		},
	})

}

func testAccVirtualDeviceContextResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model = %q
  slug = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name = %q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_device_context" "test" {
  name = %q
  device = netbox_device.test.id
  status = "active"
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName)
}

func testAccVirtualDeviceContextResourceConfig_fullWithAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model = %[5]q
  slug = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name = %[9]q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_tag" "tag1" {
  name = %[15]q
  slug = %[16]q
}

resource "netbox_tag" "tag2" {
  name = %[17]q
  slug = %[18]q
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_ip_address" "test" {
	address = %[19]q
  status  = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.test.id
}

resource "netbox_virtual_device_context" "test" {
  name = %[10]q
  device = netbox_device.test.id
  status = "active"
  identifier = 100
  tenant = netbox_tenant.test.id
  primary_ip4 = netbox_ip_address.test.id
  description = %[13]q
  comments = %[14]q

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress)
}

func testAccVirtualDeviceContextResourceConfig_fullWithAllFieldsUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model = %[5]q
  slug = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name = %[9]q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_tag" "tag1" {
  name = %[15]q
  slug = %[16]q
}

resource "netbox_tag" "tag2" {
  name = %[17]q
  slug = %[18]q
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_ip_address" "test" {
	address = %[19]q
  status  = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.test.id
}

resource "netbox_virtual_device_context" "test" {
  name = %[10]q
  device = netbox_device.test.id
  status = "active"
  identifier = 200
  tenant = netbox_tenant.test.id
  primary_ip4 = netbox_ip_address.test.id
  description = %[13]q
  comments = %[14]q

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress)
}

func testAccVirtualDeviceContextResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model = %q
  slug = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name = %q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_device_context" "test" {
  name = %q
  device = netbox_device.test.id
  status = "active"
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, description)
}

func testAccVirtualDeviceContextConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug string) string {
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

resource "netbox_tenant" "test" {
  name = "%[11]s"
  slug = "%[12]s"
}

resource "netbox_virtual_device_context" "test" {
  name = "%[10]s"
	device = netbox_device.test.id
	tenant = netbox_tenant.test.id
  status = "active"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug)
}

func testAccVirtualDeviceContextConsistencyLiteralNamesConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug string) string {
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
  status = "active"
}

resource "netbox_tenant" "test" {
  name = "%[11]s"
  slug = "%[12]s"
}

resource "netbox_virtual_device_context" "test" {
  name = "%[10]s"
  device = "%[9]s"
  tenant = "%[11]s"
  status = "active"
  depends_on = [netbox_device.test, netbox_tenant.test]
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug)
}

func TestAccVirtualDeviceContextResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site")
	siteSlug := testutil.GenerateSlug(siteName)
	mfgName := testutil.RandomName("test-mfg")
	mfgSlug := testutil.GenerateSlug(mfgName)
	dtModel := testutil.RandomName("test-dt")
	dtSlug := testutil.GenerateSlug(dtModel)
	roleName := testutil.RandomName("test-role")
	roleSlug := testutil.GenerateSlug(roleName)
	deviceName := testutil.RandomName("test-device")
	vdcName := testutil.RandomName("test-vdc-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimVirtualDeviceContextsList(context.Background()).Name([]string{vdcName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find virtual_device_context for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimVirtualDeviceContextsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete virtual_device_context: %v", err)
					}
					t.Logf("Successfully externally deleted virtual_device_context with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccVirtualDeviceContextResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-opt")
	siteSlug := testutil.RandomSlug("tf-test-site-opt")
	mfgName := testutil.RandomName("tf-test-mfg-opt")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-opt")
	dtModel := testutil.RandomName("tf-test-dt-opt")
	dtSlug := testutil.RandomSlug("tf-test-dt-opt")
	roleName := testutil.RandomName("tf-test-role-opt")
	roleSlug := testutil.RandomSlug("tf-test-role-opt")
	deviceName := testutil.RandomName("tf-test-device-opt")
	vdcName := testutil.RandomName("tf-test-vdc-opt")
	tenantName := testutil.RandomName("tf-test-tenant-opt")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-opt")
	ipAddress4 := testutil.RandomIPv4Address()
	ipAddress6 := testutil.RandomIPv6Address()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_ip_address" "ipv4" {
  address            = %[13]q
  status             = "active"
  assigned_object_id = netbox_interface.test.id
  assigned_object_type = "dcim.interface"
}

resource "netbox_interface" "test2" {
  device = netbox_device.test.id
  name   = "eth1"
  type   = "1000base-t"
}

resource "netbox_ip_address" "ipv6" {
  address            = %[14]q
  status             = "active"
  assigned_object_id = netbox_interface.test2.id
  assigned_object_type = "dcim.interface"
}

resource "netbox_virtual_device_context" "test" {
  name        = %[10]q
  device      = netbox_device.test.id
  status      = "active"
  identifier  = 42
  tenant      = netbox_tenant.test.id
  primary_ip4 = netbox_ip_address.ipv4.id
  primary_ip6 = netbox_ip_address.ipv6.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, ipAddress4, ipAddress6),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "identifier", "42"),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "primary_ip6"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_ip_address" "ipv4" {
  address            = %[13]q
  status             = "active"
  assigned_object_id = netbox_interface.test.id
  assigned_object_type = "dcim.interface"
}

resource "netbox_interface" "test2" {
  device = netbox_device.test.id
  name   = "eth1"
  type   = "1000base-t"
}

resource "netbox_ip_address" "ipv6" {
  address            = %[14]q
  status             = "active"
  assigned_object_id = netbox_interface.test2.id
  assigned_object_type = "dcim.interface"
}

resource "netbox_virtual_device_context" "test" {
  name   = %[10]q
  device = netbox_device.test.id
  status = "active"
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, ipAddress4, ipAddress6),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_virtual_device_context.test", "identifier"),
					resource.TestCheckNoResourceAttr("netbox_virtual_device_context.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_virtual_device_context.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_virtual_device_context.test", "primary_ip6"),
				),
			},
		},
	})
}
