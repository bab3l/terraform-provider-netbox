package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceResource_basic(t *testing.T) {

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
	serviceName := testutil.RandomName("tf-test-svc")

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
				Config: testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "protocol", "tcp"),
				),
			},
			{
				ResourceName:            "netbox_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccServiceResource_full(t *testing.T) {

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
	serviceName := testutil.RandomName("tf-test-svc-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated service description"

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
				Config: testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service.test", "id"),
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service.test", "description", description),
				),
			},
			{
				Config: testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_Service(t *testing.T) {

	t.Parallel()
	serviceName := testutil.RandomName("service")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "virtual_machine", vmName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
			},
		},
	})
}

func TestAccConsistency_Service_LiteralNames(t *testing.T) {

	t.Parallel()
	serviceName := testutil.RandomName("service")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("netbox_service.test", "virtual_machine", vmName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
			},
		},
	})
}

func testAccServiceResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName string) string {
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

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = %q
  protocol = "tcp"
  ports    = [22]
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName)
}

func testAccServiceResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description string) string {
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

resource "netbox_service" "test" {
  device      = netbox_device.test.id
  name        = %q
  protocol    = "tcp"
  ports       = [22, 443]
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, serviceName, description)
}

func testAccServiceConsistencyConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[2]s"
  cluster = netbox_cluster.test.id
}

resource "netbox_service" "test" {
  name = "%[1]s"
  virtual_machine = netbox_virtual_machine.test.name
  ports = [80]
  protocol = "tcp"
}
`, serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug)
}

func testAccServiceConsistencyLiteralNamesConfig(serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[2]s"
  cluster = netbox_cluster.test.id
}

resource "netbox_service" "test" {
  name = "%[1]s"
  # Use literal string name to mimic existing user state
  virtual_machine = "%[2]s"
  ports = [80]
  protocol = "tcp"

  depends_on = [netbox_virtual_machine.test]
}
`, serviceName, vmName, clusterName, clusterTypeName, clusterTypeSlug)
}
