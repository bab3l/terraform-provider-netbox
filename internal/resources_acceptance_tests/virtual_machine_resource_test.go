package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachineResource_basic(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
			{
				ResourceName:            "netbox_virtual_machine.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster"},
			},
			{
				Config:             testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccVirtualMachineResource_import(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-import")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-import")
	clusterName := testutil.RandomName("tf-test-cluster-import")
	roleName := testutil.RandomName("tf-test-vm-role-import")
	roleSlug := testutil.RandomSlug("tf-test-vm-role-import")
	tenantName := testutil.RandomName("tf-test-vm-tenant-import")
	tenantSlug := testutil.RandomSlug("tf-test-vm-tenant-import")
	platformName := testutil.RandomName("tf-test-vm-platform-import")
	platformSlug := testutil.RandomSlug("tf-test-vm-platform-import")
	vmName := testutil.RandomName("tf-test-vm-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterPlatformCleanup(platformSlug)

	testutil.RunImportTest(t, testutil.ImportTestConfig{
		ResourceName: "netbox_virtual_machine",
		Config: func() string {
			return testAccVirtualMachineResourceConfig_importCommandIDRefs(clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName)
		},
		AdditionalChecks: testutil.ValidateReferenceIDs("netbox_virtual_machine.test", "cluster", "role", "tenant", "platform"),
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckTenantDestroy,
			testutil.CheckPlatformDestroy,
		),
	})
}

func TestAccVirtualMachineResource_full(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-full")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-full")
	clusterName := testutil.RandomName("tf-test-cluster-full")
	vmName := testutil.RandomName("tf-test-vm-full")
	description := "Test VM with all fields"
	comments := "Test comments"
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role")
	deviceTypeName := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type")
	deviceName := testutil.RandomName("tf-test-device")
	configTemplateName := testutil.RandomName("tf-test-config-template")
	serial := "VM-serial-12345"
	localContextValue := "local-value"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
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
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial, localContextValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "vcpus", "2"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "memory", "2048"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "disk", "50"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "description", description),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "local_context_data", fmt.Sprintf("{\"local_key\":\"%s\"}", localContextValue)),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "config_template"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "device"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "serial", serial),
				),
			},
		},
	})
}

func TestAccVirtualMachineResource_update(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-update")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-update")
	clusterName := testutil.RandomName("tf-test-cluster-update")
	vmName := testutil.RandomName("tf-test-vm-update")
	updatedName := testutil.RandomName("tf-test-vm-updated")
	siteName := testutil.RandomName("tf-test-site-update")
	siteSlug := testutil.RandomSlug("tf-test-site-update")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-update")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-update")
	deviceRoleName := testutil.RandomName("tf-test-device-role-update")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-update")
	deviceTypeName := testutil.RandomName("tf-test-device-type-update")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-update")
	deviceName := testutil.RandomName("tf-test-device-update")
	configTemplateName := testutil.RandomName("tf-test-config-template-update")
	serial := "VM-serial-update"
	localContextValue := "local-value"
	updatedLocalContextValue := "updated-value"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterVirtualMachineCleanup(updatedName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
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
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
			{
				Config: testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, updatedName, "Updated description", "Updated comments", siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial, localContextValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "local_context_data", fmt.Sprintf("{\"local_key\":\"%s\"}", localContextValue)),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "config_template"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "device"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "serial", serial),
				),
			},
			{
				Config: testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, updatedName, "Updated description", "Updated comments", siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial, updatedLocalContextValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "local_context_data", fmt.Sprintf("{\"local_key\":\"%s\"}", updatedLocalContextValue)),
				),
			},
		},
	})
}

func TestAccConsistency_VirtualMachine_LiteralNames(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("ct")
	clusterTypeSlug := testutil.RandomSlug("ct")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVirtualMachineConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName),
			},
		},
	})

}

func TestAccVirtualMachineResource_IDPreservation(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-id")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-id")
	clusterName := testutil.RandomName("tf-test-cluster-id")
	vmName := testutil.RandomName("tf-test-vm-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
		},
	})

}

func testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %q
	cluster = netbox_cluster.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}

func testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, configTemplateName, serial, localContextValue string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = %q
	slug = %q
}

resource "netbox_manufacturer" "test" {
	name = %q
	slug = %q
}

resource "netbox_device_role" "test" {
	name    = %q
	slug    = %q
	color   = "aa1409"
	vm_role = false
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %q
	slug         = %q
}

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_device" "test" {
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
	cluster     = netbox_cluster.test.id
	name        = %q
}

resource "netbox_config_template" "test" {
	name          = %q
	template_code = "hostname {{ device.name }}"
}

resource "netbox_virtual_machine" "test" {
  name        = %q
	cluster     = netbox_cluster.test.id
  status      = "active"
  vcpus       = 2
  memory      = 2048
  disk        = 50
  description = %q
  comments    = %q
	serial      = %q
	device      = netbox_device.test.id
	config_template = netbox_config_template.test.id
	local_context_data = jsonencode({
		local_key = %q
	})
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, clusterTypeName, clusterTypeSlug, clusterName, deviceName, configTemplateName, vmName, description, comments, serial, localContextValue)
}

func TestAccConsistency_VirtualMachine_PlatformNamePersistence(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-platform")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-platform")
	clusterName := testutil.RandomName("tf-test-cluster-platform")
	platformName := testutil.RandomName("tf-test-platform")
	platformSlug := testutil.RandomSlug("tf-test-platform")
	vmName := testutil.RandomName("tf-test-vm-platform")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckPlatformDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_platformNamePersistence(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					testutil.ReferenceFieldCheck("netbox_virtual_machine.test", "platform"),
				),
			},
			{
				// Verify no drift when re-applied
				PlanOnly: true,
				Config:   testAccVirtualMachineResourceConfig_platformNamePersistence(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName),
			},
		},
	})
}

func TestAccVirtualMachineResource_ImportCommandRequiresApply(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-import")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-import")
	clusterName := testutil.RandomName("tf-test-cluster-import")
	roleName := testutil.RandomName("tf-test-vm-role-import")
	roleSlug := testutil.RandomSlug("tf-test-vm-role-import")
	tenantName := testutil.RandomName("tf-test-vm-tenant-import")
	tenantSlug := testutil.RandomSlug("tf-test-vm-tenant-import")
	platformName := testutil.RandomName("tf-test-vm-platform-import")
	platformSlug := testutil.RandomSlug("tf-test-vm-platform-import")
	vmName := testutil.RandomName("tf-test-vm-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterPlatformCleanup(platformSlug)

	// Variable to capture created resource IDs from first test case
	var vmID, clusterID, roleID, tenantID, platformID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckTenantDestroy,
			testutil.CheckPlatformDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create resources using ID references and capture the IDs
			{
				Config: testAccVirtualMachineResourceConfig_importCommandIDRefs(clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					// Capture the IDs for later use
					testutil.ExtractResourceAttr("netbox_virtual_machine.test", "id", &vmID),
					testutil.ExtractResourceAttr("netbox_cluster.test", "id", &clusterID),
					testutil.ExtractResourceAttr("netbox_device_role.test", "id", &roleID),
					testutil.ExtractResourceAttr("netbox_tenant.test", "id", &tenantID),
					testutil.ExtractResourceAttr("netbox_platform.test", "id", &platformID),
				),
			},
			// Step 2: Simulate fresh import by using ImportCommandWithID
			// This uses the SAME config (with ID references), but simulates a user
			// running `terraform import` on the command line
			{
				Config:            testAccVirtualMachineResourceConfig_importCommandIDRefs(clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName),
				ResourceName:      "netbox_virtual_machine.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportCommandWithID,
				ImportStateVerify: false,
			},
			// Step 3: Verify imported state contains numeric IDs (not names)
			// After import, reference fields MUST contain IDs to match configs using resource.id references.
			// This validates that UpdateReferenceAttribute correctly returns numeric IDs during import.
			{
				Config: testAccVirtualMachineResourceConfig_importCommandIDRefs(clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName),
				Check: resource.ComposeTestCheckFunc(
					// Use the standard helper to validate all reference fields are numeric IDs
					testutil.ReferenceFieldCheck("netbox_virtual_machine.test", "cluster"),
					testutil.ReferenceFieldCheck("netbox_virtual_machine.test", "role"),
					testutil.ReferenceFieldCheck("netbox_virtual_machine.test", "tenant"),
					testutil.ReferenceFieldCheck("netbox_virtual_machine.test", "platform"),
				),
			},
		},
	})
}

func TestAccVirtualMachineResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-tags")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-tags")
	clusterName := testutil.RandomName("tf-test-cluster-tags")
	vmName := testutil.RandomName("tf-test-vm-tags")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_virtual_machine",
		ConfigWithoutTags: func() string {
			return testAccVirtualMachineResourceConfig_tagLifecycle(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccVirtualMachineResourceConfig_tagLifecycle(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, caseTag1Uscore2)
		},
		ConfigWithDifferentTags: func() string {
			return testAccVirtualMachineResourceConfig_tagLifecycle(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, caseTag3)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
		CheckDestroy:              testutil.CheckVirtualMachineDestroy,
	})
}

func TestAccVirtualMachineResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-tagorder")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-tagorder")
	clusterName := testutil.RandomName("tf-test-cluster-tagorder")
	vmName := testutil.RandomName("tf-test-vm-tagorder")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_virtual_machine",
		ConfigWithTagsOrderA: func() string {
			return testAccVirtualMachineResourceConfig_tagOrder(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, caseTag1Uscore2)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccVirtualMachineResourceConfig_tagOrder(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, caseTag2Uscore1)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckVirtualMachineDestroy,
	})
}

func testAccVirtualMachineResourceConfig_importCommandIDRefs(clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name = %q
  slug = %q
}

resource "netbox_virtual_machine" "test" {
  name     = %q
  cluster  = netbox_cluster.test.id
  role     = netbox_device_role.test.id
  tenant   = netbox_tenant.test.id
  platform = netbox_platform.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, roleName, roleSlug, tenantName, tenantSlug, platformName, platformSlug, vmName)
}

func testAccVirtualMachineResourceConfig_platformNamePersistence(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_platform" "test" {
  name = %q
  slug = %q
}

resource "netbox_virtual_machine" "test" {
  name     = %q
  cluster  = netbox_cluster.test.id
	platform = netbox_platform.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName)
}

func testAccVirtualMachineConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %[4]q
	cluster = netbox_cluster.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}

func testAccVirtualMachineResourceConfig_tagLifecycle(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagCase string) string {
	var tagsList string
	switch tagCase {
	case caseTag1Uscore2:
		tagsList = tagsDoubleSlug
	case caseTag3:
		tagsList = tagsSingleSlug
	default:
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_cluster" "test" {
	name = %[3]q
	type = netbox_cluster_type.test.id
}

resource "netbox_tag" "tag1" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_tag" "tag2" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_tag" "tag3" {
	name = %[9]q
	slug = %[10]q
}

resource "netbox_virtual_machine" "test" {
	name    = %[4]q
	cluster = netbox_cluster.test.id
	%[11]s
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagsList)
}

func testAccVirtualMachineResourceConfig_tagOrder(clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tagOrder string) string {
	var tagsList string
	switch tagOrder {
	case caseTag1Uscore2:
		tagsList = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsList = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_cluster" "test" {
	name = %[3]q
	type = netbox_cluster_type.test.id
}

resource "netbox_tag" "tag1" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_tag" "tag2" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_virtual_machine" "test" {
	name    = %[4]q
	cluster = netbox_cluster.test.id
	%[9]s
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tagsList)
}

func TestAccVirtualMachineResource_externalDeletion(t *testing.T) {
	t.Parallel()

	vmName := testutil.RandomName("test-vm-del")
	clusterName := testutil.RandomName("test-cluster")
	clusterTypeName := testutil.RandomName("test-cluster-type")
	clusterTypeSlug := testutil.GenerateSlug(clusterTypeName)
	siteName := testutil.RandomName("test-site")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationVirtualMachinesList(context.Background()).Name([]string{vmName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find virtual_machine for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete virtual_machine: %v", err)
					}
					t.Logf("Successfully externally deleted virtual_machine with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVirtualMachineResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	// Random names for all resources
	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer")
	platformName := testutil.RandomName("tf-test-platform")
	platformSlug := testutil.RandomSlug("tf-test-platform")
	vmName := testutil.RandomName("tf-test-vm")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckTenantDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create with all optional fields
			{
				Config: testAccVirtualMachineResourceConfig_withAllFields(
					clusterTypeName, clusterTypeSlug, clusterName,
					siteName, siteSlug,
					tenantName, tenantSlug,
					roleName, roleSlug,
					manufacturerName, manufacturerSlug,
					platformName, platformSlug,
					vmName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "cluster"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "platform"),
					// New fields to test
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "status", "staged"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "vcpus", "4"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "memory", "8192"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "disk", "100"),
				),
			},
			// Step 2: Remove optional fields (should set them to null)
			{
				Config: testAccVirtualMachineResourceConfig_withoutOptionalFields(
					clusterTypeName, clusterTypeSlug, clusterName,
					siteName, siteSlug,
					tenantName, tenantSlug,
					roleName, roleSlug,
					manufacturerName, manufacturerSlug,
					platformName, platformSlug,
					vmName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "site"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "role"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "platform"),
					// New fields should be cleared (status reverts to default 'active')
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "status", "active"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "vcpus"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "memory"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine.test", "disk"),
				),
			},
			// Step 3: Re-add optional fields (verify they can be set again)
			{
				Config: testAccVirtualMachineResourceConfig_withAllFields(
					clusterTypeName, clusterTypeSlug, clusterName,
					siteName, siteSlug,
					tenantName, tenantSlug,
					roleName, roleSlug,
					manufacturerName, manufacturerSlug,
					platformName, platformSlug,
					vmName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "cluster"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "platform"),
					// New fields should be re-added
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "status", "staged"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "vcpus", "4"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "memory", "8192"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "disk", "100"),
				),
			},
		},
	})
}

func testAccVirtualMachineResourceConfig_withAllFields(
	clusterTypeName, clusterTypeSlug, clusterName,
	siteName, siteSlug,
	tenantName, tenantSlug,
	roleName, roleSlug,
	manufacturerName, manufacturerSlug,
	platformName, platformSlug,
	vmName string,
) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_cluster" "test" {
  name = "%s"
  type = netbox_cluster_type.test.id
}

resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_tenant" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name   = "%s"
  slug   = "%s"
  vm_role = true
}

resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_platform" "test" {
  name         = "%s"
  slug         = "%s"
	manufacturer = netbox_manufacturer.test.id
}

resource "netbox_virtual_machine" "test" {
  name     = "%s"
	cluster  = netbox_cluster.test.id
	site     = netbox_site.test.id
	tenant   = netbox_tenant.test.id
	role     = netbox_device_role.test.id
	platform = netbox_platform.test.id
  status   = "staged"
  vcpus    = 4
  memory   = 8192
  disk     = 100
}
`, clusterTypeName, clusterTypeSlug, clusterName,
		siteName, siteSlug,
		tenantName, tenantSlug,
		roleName, roleSlug,
		manufacturerName, manufacturerSlug,
		platformName, platformSlug,
		vmName)
}

func testAccVirtualMachineResourceConfig_withoutOptionalFields(
	clusterTypeName, clusterTypeSlug, clusterName,
	siteName, siteSlug,
	tenantName, tenantSlug,
	roleName, roleSlug,
	manufacturerName, manufacturerSlug,
	platformName, platformSlug,
	vmName string,
) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_cluster" "test" {
  name = "%s"
  type = netbox_cluster_type.test.id
}

resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_tenant" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name   = "%s"
  slug   = "%s"
  vm_role = true
}

resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_platform" "test" {
  name         = "%s"
  slug         = "%s"
	manufacturer = netbox_manufacturer.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "%s"
	cluster = netbox_cluster.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName,
		siteName, siteSlug,
		tenantName, tenantSlug,
		roleName, roleSlug,
		manufacturerName, manufacturerSlug,
		platformName, platformSlug,
		vmName)
}

func TestAccVirtualMachineResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-vm-desc")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-vm-desc")
	clusterName := testutil.RandomName("tf-test-cluster-vm-desc")
	vmName := testutil.RandomName("tf-test-vm-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_virtual_machine",
		BaseConfig: func() string {
			return testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName)
		},
		ConfigWithFields: func() string {
			return testAccVirtualMachineResourceConfig_withDescriptionAndComments(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		CheckDestroy: testutil.CheckVirtualMachineDestroy,
	})
}

func testAccVirtualMachineResourceConfig_withDescriptionAndComments(clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name        = %q
  cluster     = netbox_cluster.test.id
  description = %q
  comments    = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments)
}

// TestAccVirtualMachineResource_validationErrors tests validation error scenarios.
func TestAccVirtualMachineResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_virtual_machine",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_cluster_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  cluster = netbox_cluster.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_cluster_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "Test VM"
  cluster = netbox_cluster.test.id
  status  = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_cluster_reference": {
				Config: func() string {
					return `
resource "netbox_virtual_machine" "test" {
  name    = "Test VM"
  cluster = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_cluster_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_cluster" "test" {
  name = "Test Cluster"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "Test VM"
  cluster = netbox_cluster.test.id
  tenant  = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
