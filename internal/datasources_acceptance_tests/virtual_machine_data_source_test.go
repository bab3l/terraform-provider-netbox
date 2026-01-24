package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachineDataSource_byID(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceRoleName := testutil.RandomName("device-role")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceName := testutil.RandomName("device")
	configTemplateName := testutil.RandomName("config-template")
	serial := "VM-serial-data-source"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineDataSourceConfigByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_id", "name", vmName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_id", "status", "active"),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_machine.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_id", "device", deviceName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_id", "serial", serial),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_id", "config_template", configTemplateName),
				),
			},
		},
	})
}

func TestAccVirtualMachineDataSource_byName(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceRoleName := testutil.RandomName("device-role")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceName := testutil.RandomName("device")
	configTemplateName := testutil.RandomName("config-template")
	serial := "VM-serial-data-source"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineDataSourceConfigByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "name", vmName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "status", "active"),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_machine.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "device", deviceName),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "serial", serial),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "config_template", configTemplateName),
				),
			},
		},
	})
}

func TestAccVirtualMachineDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type-id")
	clusterTypeSlug := testutil.RandomSlug("cluster-type-id")
	clusterName := testutil.RandomName("cluster-id")
	vmName := testutil.RandomName("vm-id")
	siteName := testutil.RandomName("site-id")
	siteSlug := testutil.RandomSlug("site-id")
	manufacturerName := testutil.RandomName("manufacturer-id")
	manufacturerSlug := testutil.RandomSlug("manufacturer-id")
	deviceRoleName := testutil.RandomName("device-role-id")
	deviceRoleSlug := testutil.RandomSlug("device-role-id")
	deviceTypeName := testutil.RandomName("device-type-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	deviceName := testutil.RandomName("device-id")
	configTemplateName := testutil.RandomName("config-template-id")
	serial := "VM-serial-data-source-id"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckClusterTypeDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckVirtualMachineDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineDataSourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_virtual_machine.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_machine.by_name", "name", vmName),
				),
			},
		},
	})
}

func testAccVirtualMachineDataSourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_manufacturer" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_device_role" "test" {
	name    = "%s"
	slug    = "%s"
	color   = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "%s"
	slug         = "%s"
}

resource "netbox_cluster_type" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_cluster" "test" {
	name = "%s"
	type = netbox_cluster_type.test.id
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	cluster     = netbox_cluster.test.id
	name        = "%s"
}

resource "netbox_config_template" "test" {
	name          = "%s"
	template_code = "hostname {{ device.name }}"
}

resource "netbox_virtual_machine" "test" {
	name    = "%s"
	cluster = netbox_cluster.test.id
	status  = "active"
	serial  = "%s"
	device  = netbox_device.test.id
	config_template = netbox_config_template.test.id
}

data "netbox_virtual_machine" "by_name" {
	name = netbox_virtual_machine.test.name
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, clusterTypeName, clusterTypeSlug, clusterName, deviceName, configTemplateName, vmName, serial)
}

func testAccVirtualMachineDataSourceConfigByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_manufacturer" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_device_role" "test" {
	name    = "%s"
	slug    = "%s"
	color   = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "%s"
	slug         = "%s"
}

resource "netbox_cluster_type" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_cluster" "test" {
	name = "%s"
	type = netbox_cluster_type.test.id
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	cluster     = netbox_cluster.test.id
	name        = "%s"
}

resource "netbox_config_template" "test" {
	name          = "%s"
	template_code = "hostname {{ device.name }}"
}

resource "netbox_virtual_machine" "test" {
	name    = "%s"
	cluster = netbox_cluster.test.id
	status  = "active"
	serial  = "%s"
	device  = netbox_device.test.id
	config_template = netbox_config_template.test.id
}

data "netbox_virtual_machine" "by_id" {
	id = netbox_virtual_machine.test.id
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, clusterTypeName, clusterTypeSlug, clusterName, deviceName, configTemplateName, vmName, serial)
}

func testAccVirtualMachineDataSourceConfigByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_manufacturer" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_device_role" "test" {
	name    = "%s"
	slug    = "%s"
	color   = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "%s"
	slug         = "%s"
}

resource "netbox_cluster_type" "test" {
	name = "%s"
	slug = "%s"
}

resource "netbox_cluster" "test" {
	name = "%s"
	type = netbox_cluster_type.test.id
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	cluster     = netbox_cluster.test.id
	name        = "%s"
}

resource "netbox_config_template" "test" {
	name          = "%s"
	template_code = "hostname {{ device.name }}"
}

resource "netbox_virtual_machine" "test" {
	name    = "%s"
	cluster = netbox_cluster.test.id
	status  = "active"
	serial  = "%s"
	device  = netbox_device.test.id
	config_template = netbox_config_template.test.id
}

data "netbox_virtual_machine" "by_name" {
	name = netbox_virtual_machine.test.name
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, clusterTypeName, clusterTypeSlug, clusterName, deviceName, configTemplateName, vmName, serial)
}
