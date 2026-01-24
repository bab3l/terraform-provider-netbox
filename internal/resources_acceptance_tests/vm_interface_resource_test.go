package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	vmInterfaceNameEth0       = "eth0"
	vmInterfaceNameEth0Parent = "eth0-parent"
	vmInterfaceNameEth0Bridge = "eth0-bridge"
)

// NOTE: Custom field tests for VM interface resource are in resources_acceptance_tests_customfields package

func TestAccVMInterfaceResource_basic(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	ifaceName := testutil.InterfaceName

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
				),
			},
		},
	})
}

func TestAccVMInterfaceResource_full(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-full")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-full")
	clusterName := testutil.RandomName("tf-test-cluster-full")
	vmName := testutil.RandomName("tf-test-vm-full")
	siteName := testutil.RandomName("tf-test-site-vm-int-full")
	siteSlug := testutil.RandomSlug("tf-test-site-vm-int-full")
	untaggedVLANName := testutil.RandomName("tf-test-vlan-untagged-full")
	untaggedVLANVID := testutil.RandomVID()
	taggedVLANOneName := testutil.RandomName("tf-test-vlan-tagged1-full")
	taggedVLANOneVID := testutil.RandomVID()
	taggedVLANTwoName := testutil.RandomName("tf-test-vlan-tagged2-full")
	taggedVLANTwoVID := testutil.RandomVID()
	vrfName := testutil.RandomName("tf-test-vrf-full")
	const ifaceName = vmInterfaceNameEth0
	parentIfaceName := vmInterfaceNameEth0Parent
	bridgeIfaceName := vmInterfaceNameEth0Bridge
	description := "Test VM interface with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(parentIfaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(bridgeIfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANCleanup(untaggedVLANVID)
	cleanup.RegisterVLANCleanup(taggedVLANOneVID)
	cleanup.RegisterVLANCleanup(taggedVLANTwoVID)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_full(
					clusterTypeName,
					clusterTypeSlug,
					clusterName,
					vmName,
					ifaceName,
					description,
					siteName,
					siteSlug,
					untaggedVLANName,
					untaggedVLANVID,
					taggedVLANOneName,
					taggedVLANOneVID,
					taggedVLANTwoName,
					taggedVLANTwoVID,
					vrfName,
					parentIfaceName,
					bridgeIfaceName,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "description", description),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mode", "tagged"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "parent"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "bridge"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "untagged_vlan"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "vrf"),
					testutil.ReferenceListNumericCheck("netbox_vm_interface.test", "tagged_vlans"),
				),
			},
		},
	})
}

func TestAccConsistency_VMInterface_LiteralNames(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("ct")
	clusterTypeSlug := testutil.RandomSlug("ct")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	ifaceName := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
			},
		},
	})
}

func TestAccVMInterfaceResource_update(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-update")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-update")
	clusterName := testutil.RandomName("tf-test-cluster-update")
	vmName := testutil.RandomName("tf-test-vm-update")
	siteName := testutil.RandomName("tf-test-site-vm-int-update")
	siteSlug := testutil.RandomSlug("tf-test-site-vm-int-update")
	untaggedVLANName := testutil.RandomName("tf-test-vlan-untagged-update")
	untaggedVLANVID := testutil.RandomVID()
	taggedVLANOneName := testutil.RandomName("tf-test-vlan-tagged1-update")
	taggedVLANOneVID := testutil.RandomVID()
	taggedVLANTwoName := testutil.RandomName("tf-test-vlan-tagged2-update")
	taggedVLANTwoVID := testutil.RandomVID()
	vrfName := testutil.RandomName("tf-test-vrf-update")
	const ifaceName = vmInterfaceNameEth0
	updatedIfaceName := "eth1"
	parentIfaceName := vmInterfaceNameEth0Parent
	bridgeIfaceName := vmInterfaceNameEth0Bridge

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(updatedIfaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(parentIfaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(bridgeIfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANCleanup(untaggedVLANVID)
	cleanup.RegisterVLANCleanup(taggedVLANOneVID)
	cleanup.RegisterVLANCleanup(taggedVLANTwoVID)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
				),
			},
			{
				Config: testAccVMInterfaceResourceConfig_full(
					clusterTypeName,
					clusterTypeSlug,
					clusterName,
					vmName,
					updatedIfaceName,
					"Updated description",
					siteName,
					siteSlug,
					untaggedVLANName,
					untaggedVLANVID,
					taggedVLANOneName,
					taggedVLANOneVID,
					taggedVLANTwoName,
					taggedVLANTwoVID,
					vrfName,
					parentIfaceName,
					bridgeIfaceName,
				),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", updatedIfaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "description", "Updated description"),
					testutil.ReferenceListNumericCheck("netbox_vm_interface.test", "tagged_vlans"),
				),
			},
		},
	})
}

func TestAccVMInterfaceResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-tags")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-tags")
	clusterName := testutil.RandomName("tf-test-cluster-tags")
	vmName := testutil.RandomName("tf-test-vm-tags")
	ifaceName := vmInterfaceNameEth0

	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_vm_interface",
		ConfigWithoutTags: func() string {
			return testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName)
		},
		ConfigWithTags: func() string {
			return testAccVMInterfaceResourceConfig_withTags(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				ifaceName,
				tag1Name,
				tag1Slug,
				tag2Name,
				tag2Slug,
				tag3Name,
				tag3Slug,
				"netbox_tag.tag1.slug, netbox_tag.tag2.slug",
			)
		},
		ConfigWithDifferentTags: func() string {
			return testAccVMInterfaceResourceConfig_withTags(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				ifaceName,
				tag1Name,
				tag1Slug,
				tag2Name,
				tag2Slug,
				tag3Name,
				tag3Slug,
				"netbox_tag.tag2.slug, netbox_tag.tag3.slug",
			)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
	})
}

func TestAccVMInterfaceResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-tag-order")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-tag-order")
	clusterName := testutil.RandomName("tf-test-cluster-tag-order")
	vmName := testutil.RandomName("tf-test-vm-tag-order")
	ifaceName := vmInterfaceNameEth0

	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_vm_interface",
		ConfigWithTagsOrderA: func() string {
			return testAccVMInterfaceResourceConfig_withTags(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				ifaceName,
				tag1Name,
				tag1Slug,
				tag2Name,
				tag2Slug,
				tag3Name,
				tag3Slug,
				"netbox_tag.tag1.slug, netbox_tag.tag2.slug",
			)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccVMInterfaceResourceConfig_withTags(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				ifaceName,
				tag1Name,
				tag1Slug,
				tag2Name,
				tag2Slug,
				tag3Name,
				tag3Slug,
				"netbox_tag.tag2.slug, netbox_tag.tag1.slug",
			)
		},
		ExpectedTagCount: 2,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
	})
}

func TestAccVMInterfaceResource_externalDeletion(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-ext-del")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-ext-del")
	clusterName := testutil.RandomName("tf-test-cluster-ext-del")
	vmName := testutil.RandomName("tf-test-vm-ext-del")
	ifaceName := testutil.RandomName("eth-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationInterfacesList(context.Background()).NameIc([]string{ifaceName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find vm_interface for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationInterfacesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete vm_interface: %v", err)
					}
					t.Logf("Successfully externally deleted vm_interface with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVMInterfaceResource_import(t *testing.T) {
	t.Parallel()

	clusterTypeName := "test-cluster-type-" + testutil.GenerateSlug("ct")
	clusterTypeSlug := testutil.GenerateSlug("ct")
	clusterName := "test-cluster-" + testutil.GenerateSlug("cluster")
	vmName := "test-vm-" + testutil.GenerateSlug("vm")
	ifaceName := "test-iface-" + testutil.GenerateSlug("iface")
	siteName := "test-site-" + testutil.GenerateSlug("site")
	siteSlug := testutil.GenerateSlug("site")
	untaggedVLANName := "test-vlan-untagged-" + testutil.GenerateSlug("vlan")
	untaggedVLANVID := testutil.RandomVID()
	taggedVLANOneName := "test-vlan-tagged1-" + testutil.GenerateSlug("vlan1")
	taggedVLANOneVID := testutil.RandomVID()
	taggedVLANTwoName := "test-vlan-tagged2-" + testutil.GenerateSlug("vlan2")
	taggedVLANTwoVID := testutil.RandomVID()
	vrfName := "test-vrf-" + testutil.GenerateSlug("vrf")
	parentIfaceName := "test-iface-parent-" + testutil.GenerateSlug("iface")
	bridgeIfaceName := "test-iface-bridge-" + testutil.GenerateSlug("iface")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(parentIfaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(bridgeIfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANCleanup(untaggedVLANVID)
	cleanup.RegisterVLANCleanup(taggedVLANOneVID)
	cleanup.RegisterVLANCleanup(taggedVLANTwoVID)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_full(
					clusterTypeName,
					clusterTypeSlug,
					clusterName,
					vmName,
					ifaceName,
					"Import description",
					siteName,
					siteSlug,
					untaggedVLANName,
					untaggedVLANVID,
					taggedVLANOneName,
					taggedVLANOneVID,
					taggedVLANTwoName,
					taggedVLANTwoVID,
					vrfName,
					parentIfaceName,
					bridgeIfaceName,
				),
			},
			{
				ResourceName:            "netbox_vm_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_vm_interface.test", "virtual_machine"),
					testutil.ReferenceFieldCheck("netbox_vm_interface.test", "untagged_vlan"),
					testutil.ReferenceFieldCheck("netbox_vm_interface.test", "vrf"),
					testutil.ReferenceFieldCheck("netbox_vm_interface.test", "parent"),
					testutil.ReferenceFieldCheck("netbox_vm_interface.test", "bridge"),
					testutil.ReferenceListNumericCheck("netbox_vm_interface.test", "tagged_vlans"),
				),
			},
			{
				Config: testAccVMInterfaceResourceConfig_full(
					clusterTypeName,
					clusterTypeSlug,
					clusterName,
					vmName,
					ifaceName,
					"Import description",
					siteName,
					siteSlug,
					untaggedVLANName,
					untaggedVLANVID,
					taggedVLANOneName,
					taggedVLANOneVID,
					taggedVLANTwoName,
					taggedVLANTwoVID,
					vrfName,
					parentIfaceName,
					bridgeIfaceName,
				),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVMInterfaceResource_IDPreservation(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-id")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-id")
	clusterName := testutil.RandomName("tf-test-cluster-id")
	vmName := testutil.RandomName("tf-test-vm-id")
	ifaceName := testutil.RandomName("eth-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
				),
			},
		},
	})
}

func testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName string) string {
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

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName)
}

func testAccVMInterfaceResourceConfig_withTags(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagList string) string {
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

resource "netbox_tag" "tag1" {
	name = %q
	slug = %q
}

resource "netbox_tag" "tag2" {
	name = %q
	slug = %q
}

resource "netbox_tag" "tag3" {
	name = %q
	slug = %q
}

resource "netbox_vm_interface" "test" {
	virtual_machine = netbox_virtual_machine.test.name
	name            = %q
	tags            = [%s]
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, ifaceName, tagList)
}

func testAccVMInterfaceResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description, siteName, siteSlug, untaggedVLANName string, untaggedVLANVID int32, taggedVLANOneName string, taggedVLANOneVID int32, taggedVLANTwoName string, taggedVLANTwoVID int32, vrfName, parentIfaceName, bridgeIfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
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
	site    = netbox_site.test.name
}

resource "netbox_vlan" "untagged" {
	name = %q
	vid  = %d
	site = netbox_site.test.id
}

resource "netbox_vlan" "tagged_1" {
	name = %q
	vid  = %d
	site = netbox_site.test.id
}

resource "netbox_vlan" "tagged_2" {
	name = %q
	vid  = %d
	site = netbox_site.test.id
}

resource "netbox_vrf" "test" {
	name = %q
}

resource "netbox_vm_interface" "parent" {
	virtual_machine = netbox_virtual_machine.test.name
	name            = %q
}

resource "netbox_vm_interface" "bridge" {
	virtual_machine = netbox_virtual_machine.test.name
	name            = %q
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %q
  enabled         = true
  mtu             = 1500
	mac_address     = "00:11:22:33:44:55"
  description     = %q
	mode            = "tagged"
	parent          = netbox_vm_interface.parent.id
	bridge          = netbox_vm_interface.bridge.id
	untagged_vlan   = netbox_vlan.untagged.id
	tagged_vlans    = [netbox_vlan.tagged_1.id, netbox_vlan.tagged_2.id]
	vrf             = netbox_vrf.test.id
}
`, clusterTypeName, clusterTypeSlug, siteName, siteSlug, clusterName, vmName, untaggedVLANName, untaggedVLANVID, taggedVLANOneName, taggedVLANOneVID, taggedVLANTwoName, taggedVLANTwoVID, vrfName, parentIfaceName, bridgeIfaceName, ifaceName, description)
}

func TestAccConsistency_VMInterface(t *testing.T) {
	t.Parallel()

	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	interfaceName := vmInterfaceNameEth0
	macAddress := "AA:BB:CC:DD:EE:FF" // Uppercase to test case sensitivity
	vlanName := testutil.RandomName("vlan")
	vlanVid := int32(100)
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mac_address", macAddress),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			{
				// Verify no drift
				PlanOnly: true,
				Config:   testAccVMInterfaceConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

func testAccVMInterfaceConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName string, vlanVid int32, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_site" "test" {
  name = "%[9]s"
  slug = "%[10]s"
}

resource "netbox_cluster" "test" {
  name = "%[2]s"
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[1]s"
  cluster = netbox_cluster.test.id
  site = netbox_site.test.id
}

resource "netbox_vlan" "test" {
  name = "%[7]s"
  vid  = %[8]d
  site = netbox_site.test.id
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name = "%[5]s"
  mac_address = "%[6]s"
  mode = "access"
  untagged_vlan = netbox_vlan.test.name
}
`, vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName, vlanVid, siteName, siteSlug)
}

func testAccVMInterfaceConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName string) string {
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

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName)
}

// TestAccVMInterface_WithVRF verifies that VM interfaces with VRF maintain consistency.
func TestAccVMInterface_WithVRF(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-vrf")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-vrf")
	clusterName := testutil.RandomName("tf-test-cluster-vrf")
	vmName := testutil.RandomName("tf-test-vm-vrf")
	ifaceName := testutil.RandomName("eth-vrf")
	vrfName := testutil.RandomName("tf-test-vrf")
	vrfRD := testutil.RandomSlug("vrf-rd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create VM interface with VRF by name
			{
				Config: testAccVMInterfaceConfig_withVRF(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vrfName, vrfRD),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
					// VRF should be stored as name, not ID
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "vrf", vrfName),
				),
			},
			// Step 2: Refresh state and verify no drift
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "vrf", vrfName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_withVRF(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vrfName, vrfRD),
			},
		},
	})
}

func testAccVMInterfaceConfig_withVRF(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vrfName, vrfRD string) string {
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

resource "netbox_vrf" "test" {
  name = %[6]q
  rd   = %[7]q
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
  vrf             = netbox_vrf.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vrfName, vrfRD)
}

// TestAccVMInterface_VirtualMachineNameNotID verifies that when virtual_machine is specified by name,
// the state stores the name (not the numeric ID) and remains consistent after refresh.
func TestAccVMInterface_VirtualMachineNameNotID(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-vmname")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-vmname")
	clusterName := testutil.RandomName("tf-test-cluster-vmname")
	vmName := testutil.RandomName("tf-test-vm-vmname")
	ifaceName := testutil.RandomName("eth-vmname")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VM name
			{
				Config: testAccVMInterfaceConfig_vmByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// virtual_machine should be stored as NAME, not numeric ID
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
				),
			},
			// Step 2: Refresh state and verify no drift
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					// After refresh, virtual_machine should still be the name
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_vmByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
			},
		},
	})
}

func testAccVMInterfaceConfig_vmByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName string) string {
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

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName)
}

// TestAccVMInterface_UntaggedVLANNameNotID verifies that when untagged_vlan is specified by name,
// the state stores the name (not the numeric ID) and remains consistent after refresh.
func TestAccVMInterface_UntaggedVLANNameNotID(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-uvlan")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-uvlan")
	clusterName := testutil.RandomName("tf-test-cluster-uvlan")
	vmName := testutil.RandomName("tf-test-vm-uvlan")
	ifaceName := testutil.RandomName("eth-uvlan")
	vlanName := testutil.RandomName("tf-test-vlan-uvlan")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-uvlan")
	siteSlug := testutil.RandomSlug("tf-test-site-uvlan")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VLAN name
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// untagged_vlan should be stored as NAME, not numeric ID
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mode", "access"),
				),
			},
			// Step 2: Refresh state and verify no drift
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					// After refresh, untagged_vlan should still be the name
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

func testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName string, vlanVid int32, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %[4]q
  cluster = netbox_cluster.test.name
  site    = netbox_site.test.name
}

resource "netbox_vlan" "test" {
  name = %[6]q
  vid  = %[7]d
  site = netbox_site.test.name
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
  mode            = "access"
  untagged_vlan   = netbox_vlan.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug)
}

// TestAccVMInterface_UntaggedVLANByID_StoresID verifies that when untagged_vlan is specified
// by ID (via resource reference), the ID is preserved consistently.
func TestAccVMInterface_UntaggedVLANByID_StoresID(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-uvid")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-uvid")
	clusterName := testutil.RandomName("tf-test-cluster-uvid")
	vmName := testutil.RandomName("tf-test-vm-uvid")
	ifaceName := testutil.RandomName("eth-uvid")
	vlanName := testutil.RandomName("tf-test-vlan-uvid")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-uvid")
	siteSlug := testutil.RandomSlug("tf-test-site-uvid")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VLAN ID (via resource reference)
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// When specified by ID, untagged_vlan should be stored as ID
					resource.TestCheckResourceAttrPair("netbox_vm_interface.test", "untagged_vlan", "netbox_vlan.test", "id"),
				),
			},
			// Step 2: Refresh state and verify ID is still stored
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("netbox_vm_interface.test", "untagged_vlan", "netbox_vlan.test", "id"),
				),
			},
			// Step 3: Plan only - verify no changes detected
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_untaggedVLANByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

func testAccVMInterfaceConfig_untaggedVLANByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName string, vlanVid int32, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %[4]q
  cluster = netbox_cluster.test.name
  site    = netbox_site.test.name
}

resource "netbox_vlan" "test" {
  name = %[6]q
  vid  = %[7]d
  site = netbox_site.test.name
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
  mode            = "access"
  untagged_vlan   = netbox_vlan.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug)
}

// TestAccVMInterface_UnknownValueResolution tests that when a reference value
// starts as unknown (computed from another resource), it resolves correctly.
// This simulates the real-world scenario where netbox_vlan.test.name is unknown
// during planning and only becomes known during apply.
func TestAccVMInterface_UnknownValueResolution(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-unk")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-unk")
	clusterName := testutil.RandomName("tf-test-cluster-unk")
	vmName := testutil.RandomName("tf-test-vm-unk")
	ifaceName := testutil.RandomName("eth-unk")
	vlanName := testutil.RandomName("tf-test-vlan-unk")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-unk")
	siteSlug := testutil.RandomSlug("tf-test-site-unk")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create all resources - the key here is that untagged_vlan
			// uses netbox_vlan.test.name which is unknown during planning
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// After create, untagged_vlan should be the NAME, not numeric ID
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			// Step 2: Run another apply with the same config - this checks for drift
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					// After second apply, untagged_vlan should STILL be the name
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			// Step 3: Plan only - verify no changes detected (critical test for drift)
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

// TestAccVMInterface_SwitchFromIDToName tests that when the config is changed
// from using an ID to using a name, the state updates correctly without drift.
func TestAccVMInterface_SwitchFromIDToName(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-switch")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-switch")
	clusterName := testutil.RandomName("tf-test-cluster-switch")
	vmName := testutil.RandomName("tf-test-vm-switch")
	ifaceName := testutil.RandomName("eth-switch")
	vlanName := testutil.RandomName("tf-test-vlan-switch")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-switch")
	siteSlug := testutil.RandomSlug("tf-test-site-switch")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VLAN ID
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// untagged_vlan should be stored as ID
					resource.TestCheckResourceAttrPair("netbox_vm_interface.test", "untagged_vlan", "netbox_vlan.test", "id"),
				),
			},
			// Step 2: Switch to using VLAN name in config
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					// After switching to name, untagged_vlan should now be stored as NAME
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			// Step 3: Verify no drift after switch
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

// TestAccVMInterface_ImportThenUseByName tests that after importing a resource,
// if the config uses names, the state correctly updates to use names.
func TestAccVMInterface_ImportThenUseByName(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-import")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-import")
	clusterName := testutil.RandomName("tf-test-cluster-import")
	vmName := testutil.RandomName("tf-test-vm-import")
	ifaceName := testutil.RandomName("eth-import")
	vlanName := testutil.RandomName("tf-test-vlan-import")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with ID (this sets up the scenario)
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByID(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
				),
			},
			// Step 2: Import the resource
			{
				ResourceName:            "netbox_vm_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine", "untagged_vlan"},
			},
			// Step 3: Apply with name config - this should update state to use name
			{
				Config: testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					// After apply with name, state should have the name
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "untagged_vlan", vlanName),
				),
			},
			// Step 4: Verify no drift
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_untaggedVLANByName(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

// TestAccVMInterface_ModeNotInConfig verifies that when mode is not specified in config,
// but the interface has a mode in Netbox (e.g., from VLAN assignment), there is no drift.
// This tests the bug where Terraform would show: mode = "access" -> null.
func TestAccVMInterface_ModeNotInConfig(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-ct-mode")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-mode")
	clusterName := testutil.RandomName("tf-test-cluster-mode")
	vmName := testutil.RandomName("tf-test-vm-mode")
	ifaceName := testutil.RandomName("eth-mode")
	vlanName := testutil.RandomName("tf-test-vlan-mode")
	vlanVid := testutil.RandomVID()
	siteName := testutil.RandomName("tf-test-site-mode")
	siteSlug := testutil.RandomSlug("tf-test-site-mode")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with mode specified
			{
				Config: testAccVMInterfaceConfig_withMode(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mode", "access"),
				),
			},
			// Step 2: Remove mode from config (but interface still has mode in Netbox)
			// This should NOT show drift or cause errors
			{
				Config: testAccVMInterfaceConfig_withoutMode(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					// mode should be null in state when not in config
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "mode"),
				),
			},
			// Step 3: Verify no drift
			{
				PlanOnly: true,
				Config:   testAccVMInterfaceConfig_withoutMode(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, siteName, siteSlug),
			},
		},
	})
}

func testAccVMInterfaceConfig_withMode(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName string, vlanVid int32, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %[4]q
  cluster = netbox_cluster.test.name
  site    = netbox_site.test.name
}

resource "netbox_vlan" "test" {
  name = %[6]q
  vid  = %[7]d
  site = netbox_site.test.name
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
  mode            = "access"
  untagged_vlan   = netbox_vlan.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, vlanName, vlanVid, siteName, siteSlug)
}

func testAccVMInterfaceConfig_withoutMode(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_site" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_cluster" "test" {
  name = %[3]q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %[4]q
  cluster = netbox_cluster.test.name
  site    = netbox_site.test.name
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name            = %[5]q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, siteName, siteSlug)
}

// TestAccVMInterfaceResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccVMInterfaceResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-remove")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-remove")
	clusterName := testutil.RandomName("tf-test-cluster-remove")
	siteName := testutil.RandomName("tf-test-site-vm-iface")
	siteSlug := testutil.RandomSlug("tf-test-site-vm-iface")
	vmName := testutil.RandomName("tf-test-vm-remove")
	ifaceName := vmInterfaceNameEth0
	vlanName := testutil.RandomName("tf-test-vlan-vm-iface")
	vlanVID := testutil.RandomVID()
	taggedVLANOneName := testutil.RandomName("tf-test-vlan-tagged1-vm-iface")
	taggedVLANOneVID := testutil.RandomVID()
	taggedVLANTwoName := testutil.RandomName("tf-test-vlan-tagged2-vm-iface")
	taggedVLANTwoVID := testutil.RandomVID()
	vrfName := testutil.RandomName("tf-test-vrf-vm-iface")
	parentIfaceName := vmInterfaceNameEth0Parent
	bridgeIfaceName := vmInterfaceNameEth0Bridge

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(parentIfaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(bridgeIfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANCleanup(vlanVID)
	cleanup.RegisterVLANCleanup(taggedVLANOneVID)
	cleanup.RegisterVLANCleanup(taggedVLANTwoVID)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create VM interface with all optional fields
			{
				Config: testAccVMInterfaceResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, vmName, ifaceName, vlanName, vlanVID, taggedVLANOneName, taggedVLANOneVID, taggedVLANTwoName, taggedVLANTwoVID, vrfName, parentIfaceName, bridgeIfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mode", "tagged"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "parent"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "bridge"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "untagged_vlan"),
					testutil.ReferenceListNumericCheck("netbox_vm_interface.test", "tagged_vlans"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "vrf"),
				),
			},
			// Step 2: Remove optional fields - should clear them
			{
				Config: testAccVMInterfaceResourceConfig_withoutOptionalFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, vmName, ifaceName, vlanName, vlanVID, vrfName, parentIfaceName, bridgeIfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "enabled", "true"), // Should revert to default
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "mtu"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "mac_address"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "mode"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "parent"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "bridge"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "untagged_vlan"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "tagged_vlans"),
					resource.TestCheckNoResourceAttr("netbox_vm_interface.test", "vrf"),
				),
			},
			// Step 3: Re-add optional fields - verify they can be set again
			{
				Config: testAccVMInterfaceResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, vmName, ifaceName, vlanName, vlanVID, taggedVLANOneName, taggedVLANOneVID, taggedVLANTwoName, taggedVLANTwoVID, vrfName, parentIfaceName, bridgeIfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mac_address", "00:11:22:33:44:55"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mode", "tagged"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "parent"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "bridge"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "untagged_vlan"),
					testutil.ReferenceListNumericCheck("netbox_vm_interface.test", "tagged_vlans"),
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "vrf"),
				),
			},
		},
	})
}

func testAccVMInterfaceResourceConfig_withAllFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, vmName, ifaceName, vlanName string, vlanVID int32, taggedVLANOneName string, taggedVLANOneVID int32, taggedVLANTwoName string, taggedVLANTwoVID int32, vrfName, parentIfaceName, bridgeIfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %q
  cluster = netbox_cluster.test.id
	site    = netbox_site.test.name
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
  site = netbox_site.test.id
}

resource "netbox_vlan" "tagged_1" {
	name = %q
	vid  = %d
	site = netbox_site.test.id
}

resource "netbox_vlan" "tagged_2" {
	name = %q
	vid  = %d
	site = netbox_site.test.id
}

resource "netbox_vrf" "test" {
  name = %q
}

resource "netbox_vm_interface" "parent" {
	virtual_machine = netbox_virtual_machine.test.id
	name            = %q
}

resource "netbox_vm_interface" "bridge" {
	virtual_machine = netbox_virtual_machine.test.id
	name            = %q
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q
  enabled         = false
  mtu             = 1500
  mac_address     = "00:11:22:33:44:55"
	mode            = "tagged"
	parent          = netbox_vm_interface.parent.id
	bridge          = netbox_vm_interface.bridge.id
  untagged_vlan   = netbox_vlan.test.id
	tagged_vlans    = [netbox_vlan.tagged_1.id, netbox_vlan.tagged_2.id]
  vrf             = netbox_vrf.test.id
}
`, clusterTypeName, clusterTypeSlug, siteName, siteSlug, clusterName, vmName, vlanName, vlanVID, taggedVLANOneName, taggedVLANOneVID, taggedVLANTwoName, taggedVLANTwoVID, vrfName, parentIfaceName, bridgeIfaceName, ifaceName)
}

func testAccVMInterfaceResourceConfig_withoutOptionalFields(clusterTypeName, clusterTypeSlug, clusterName, siteName, siteSlug, vmName, ifaceName, vlanName string, vlanVID int32, vrfName, parentIfaceName, bridgeIfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
  site = netbox_site.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = %q
  cluster = netbox_cluster.test.id
	site    = netbox_site.test.name
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
  site = netbox_site.test.id
}

resource "netbox_vrf" "test" {
  name = %q
}

resource "netbox_vm_interface" "parent" {
	virtual_machine = netbox_virtual_machine.test.id
	name            = %q
}

resource "netbox_vm_interface" "bridge" {
	virtual_machine = netbox_virtual_machine.test.id
	name            = %q
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q
}
`, clusterTypeName, clusterTypeSlug, siteName, siteSlug, clusterName, vmName, vlanName, vlanVID, vrfName, parentIfaceName, bridgeIfaceName, ifaceName)
}

func TestAccVMInterfaceResource_removeDescription(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-vm-int-desc")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-vm-int-desc")
	clusterName := testutil.RandomName("tf-test-cluster-vm-int-desc")
	vmName := testutil.RandomName("tf-test-vm-int-desc")
	ifaceName := testutil.RandomName("tf-test-int-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_vm_interface",
		BaseConfig: func() string {
			return testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName)
		},
		ConfigWithFields: func() string {
			return testAccVMInterfaceResourceConfig_withDescription(
				clusterTypeName,
				clusterTypeSlug,
				clusterName,
				vmName,
				ifaceName,
				"Test description",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
		},
		CheckDestroy: testutil.CheckVMInterfaceDestroy,
	})
}

func testAccVMInterfaceResourceConfig_withDescription(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description string) string {
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

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q
  description     = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description)
}

func TestAccVMInterfaceResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_vm_interface",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_virtual_machine": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_vm_interface" "test" {
  # virtual_machine missing
  name = "eth0"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_cluster_type" "test" {
  name = "test-cluster-type"
  slug = "test-cluster-type"
}

resource "netbox_cluster" "test" {
  name = "test-cluster"
  cluster_type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name    = "test-vm"
  cluster = netbox_cluster.test.id
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"tagged_vlans_without_tagged_mode": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_vm_interface" "test" {
  virtual_machine = "test-vm"
  name            = "eth0"
  mode            = "access"
  tagged_vlans    = ["10"]
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidValue,
			},
		},
	})
}
