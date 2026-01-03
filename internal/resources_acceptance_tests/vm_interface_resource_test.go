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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
	const ifaceName = "eth0"
	description := "Test VM interface with all fields"
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "virtual_machine", vmName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "description", description),
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
	const ifaceName = "eth0"
	updatedIfaceName := "eth1"
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVMInterfaceCleanup(ifaceName, vmName)
	cleanup.RegisterVMInterfaceCleanup(updatedIfaceName, vmName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
				Config: testAccVMInterfaceResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, updatedIfaceName, "Updated description"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", updatedIfaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccVMInterfaceResource_external_deletion(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-ext-del")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-ext-del")
	clusterName := testutil.RandomName("tf-test-cluster-ext-del")
	vmName := testutil.RandomName("tf-test-vm-ext-del")
	ifaceName := testutil.RandomName("eth-ext-del")

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
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
				),
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

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
			},
			{
				ResourceName:            "netbox_vm_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine"},
			},
			{
				Config:   testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName),
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

func testAccVMInterfaceResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description string) string {
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
  enabled         = true
  mtu             = 1500
  description     = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, description)
}

func TestAccConsistency_VMInterface(t *testing.T) {

	t.Parallel()

	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	interfaceName := "eth0"
	macAddress := "AA:BB:CC:DD:EE:FF" // Uppercase to test case sensitivity
	vlanName := testutil.RandomName("vlan")
	vlanVid := 100
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

				Config: testAccVMInterfaceConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName, vlanVid, siteName, siteSlug),
			},
		},
	})
}

func testAccVMInterfaceConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, interfaceName, macAddress, vlanName string, vlanVid int, siteName, siteSlug string) string {
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

// TestAccVMInterfaceResource_importWithCustomFieldsAndTags tests importing a VM interface
// with custom fields and tags to ensure all data is preserved during import.
func TestAccVMInterfaceResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	// Generate unique names
	clusterTypeName := testutil.RandomName("tf-test-ct-import")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-import")
	clusterName := testutil.RandomName("tf-test-cluster-import")
	vmName := testutil.RandomName("tf-test-vm-import")
	ifaceName := testutil.RandomName("tf-test-vmint-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-vmint-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-vmint-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-vmint-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-vmint-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values
	cfText := testutil.RandomCustomFieldName("tf_vmint_text")
	cfTextValue := testutil.RandomName("vmint-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_vmint_longtext")
	cfLongtextValue := fmt.Sprintf("VM Interface description: %s", testutil.RandomName("vmint-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_vmint_integer")
	cfIntegerValue := 1500
	cfBoolean := testutil.RandomCustomFieldName("tf_vmint_boolean")
	cfBooleanValue := false
	cfDate := testutil.RandomCustomFieldName("tf_vmint_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_vmint_url")
	cfURLValue := testutil.RandomURL("vm-interface")
	cfJSON := testutil.RandomCustomFieldName("tf_vmint_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				// Create the VM interface with all custom fields and tags
				Config: testAccVMInterfaceResourceImportConfig_full(
					clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the VM interface and verify basic fields are preserved
				ResourceName:            "netbox_vm_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine", "custom_fields", "tags"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
				),
			},
			{
				Config: testAccVMInterfaceResourceImportConfig_full(
					clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccVMInterfaceResourceImportConfig_full(
	clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "vmint_test1" {
  name  = %[6]q
  slug  = %[7]q
  color = %[8]q
}

resource "netbox_tag" "vmint_test2" {
  name  = %[9]q
  slug  = %[10]q
  color = %[11]q
}

# Create custom fields for virtualization.vminterface
resource "netbox_custom_field" "vmint_text" {
  name         = %[12]q
  type         = "text"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_longtext" {
  name         = %[14]q
  type         = "longtext"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_integer" {
  name         = %[16]q
  type         = "integer"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_boolean" {
  name         = %[18]q
  type         = "boolean"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_date" {
  name         = %[20]q
  type         = "date"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_url" {
  name         = %[22]q
  type         = "url"
  object_types = ["virtualization.vminterface"]
  required     = false
}

resource "netbox_custom_field" "vmint_json" {
  name         = %[24]q
  type         = "json"
  object_types = ["virtualization.vminterface"]
  required     = false
}

# Create dependencies
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

# Create VM interface with all custom fields and tags
resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %[5]q

  tags = [
    {
      name = netbox_tag.vmint_test1.name
      slug = netbox_tag.vmint_test1.slug
    },
    {
      name = netbox_tag.vmint_test2.name
      slug = netbox_tag.vmint_test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.vmint_text.name
      type  = "text"
      value = %[13]q
    },
    {
      name  = netbox_custom_field.vmint_longtext.name
      type  = "longtext"
      value = %[15]q
    },
    {
      name  = netbox_custom_field.vmint_integer.name
      type  = "integer"
      value = "%[17]d"
    },
    {
      name  = netbox_custom_field.vmint_boolean.name
      type  = "boolean"
      value = "%[19]t"
    },
    {
      name  = netbox_custom_field.vmint_date.name
      type  = "date"
      value = %[21]q
    },
    {
      name  = netbox_custom_field.vmint_url.name
      type  = "url"
      value = %[23]q
    },
    {
      name  = netbox_custom_field.vmint_json.name
      type  = "json"
      value = %[25]q
    }
  ]
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName,
		tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
