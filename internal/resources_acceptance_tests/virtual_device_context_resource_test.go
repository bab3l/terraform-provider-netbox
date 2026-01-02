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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVirtualDeviceContextDestroy,

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

				ResourceName: "netbox_virtual_device_context.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceConfig_fullWithAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "identifier", "100"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", description),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "comments", comments),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "primary_ip4"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config: testAccVirtualDeviceContextResourceConfig_fullWithAllFieldsUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "identifier", "200"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.0.value", "updated_value"),
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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVirtualDeviceContextDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDeviceContextConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "device", deviceName),

					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccVirtualDeviceContextConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

				Config: testAccVirtualDeviceContextConsistencyLiteralNamesConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, vdcName, tenantName, tenantSlug),
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVirtualDeviceContextDestroy,
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

func testAccVirtualDeviceContextResourceConfig_fullWithAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := testutil.RandomCustomFieldName("test_field")
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

resource "netbox_custom_field" "test_field" {
  name         = %[19]q
  object_types = ["dcim.virtualdevicecontext"]
  type         = "text"
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_ip_address" "test" {
  address = "192.0.2.1/32"
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

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccVirtualDeviceContextResourceConfig_fullWithAllFieldsUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := testutil.RandomCustomFieldName("test_field")
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

resource "netbox_custom_field" "test_field" {
  name         = %[19]q
  object_types = ["dcim.virtualdevicecontext"]
  type         = "text"
}

resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

resource "netbox_ip_address" "test" {
  address = "192.0.2.1/32"
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

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, vdcName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
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

  device = netbox_device.test.name

  tenant = netbox_tenant.test.name

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

func TestAccVirtualDeviceContextResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	vdcName := testutil.RandomName("vdc")
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	// Note: VDC cleanup function requires ID (int32), but we only have name - skipping VDC cleanup
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVirtualDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextResourceImportConfig_full(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "name", vdcName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_virtual_device_context.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_virtual_device_context.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccVirtualDeviceContextResourceImportConfig_full(vdcName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "9e9e9e"
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["dcim.virtualdevicecontext"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["dcim.virtualdevicecontext"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Virtual Device Context with comprehensive custom fields and tags
resource "netbox_virtual_device_context" "test" {
  name   = %q
  device = netbox_device.test.id
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this virtual device context resource for testing purposes."
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
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
`, tenantName, tenantSlug, siteName, siteSlug, mfgName, mfgSlug, roleName, roleSlug, dtModel, dtSlug, deviceName, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, vdcName)
}
