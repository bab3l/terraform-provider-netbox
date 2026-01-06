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
		},
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
				Config: testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "vcpus", "2"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "memory", "2048"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "disk", "50"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "description", description),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "comments", comments),
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterVirtualMachineCleanup(updatedName)
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
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
				),
			},
			{
				Config: testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, updatedName, "Updated description", "Updated comments"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "description", "Updated description"),
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
  cluster = netbox_cluster.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
}

func testAccVirtualMachineResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments string) string {
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
  cluster     = netbox_cluster.test.name
  status      = "active"
  vcpus       = 2
  memory      = 2048
  disk        = 50
  description = %q
  comments    = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, description, comments)
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
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "platform", platformName),
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
  platform = netbox_platform.test.name
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
  cluster = netbox_cluster.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName)
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
