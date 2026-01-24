package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachinePrimaryIPResource_basic(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, "", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
		},
	})
}

func TestAccVirtualMachinePrimaryIPResource_full(t *testing.T) {
	t.Parallel()

	requireIPv6Support(t)

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))
	ip6 := fmt.Sprintf("2001:db8:%d::1/64", acctest.RandIntRange(1, 65535))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4)
	cleanup.RegisterIPAddressCleanup(ip6)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, ip6, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
		},
	})
}

func TestAccVirtualMachinePrimaryIPResource_update(t *testing.T) {
	t.Parallel()

	requireIPv6Support(t)

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4a := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))
	ip4b := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))
	ip6 := fmt.Sprintf("2001:db8:%d::1/64", acctest.RandIntRange(1, 65535))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4a)
	cleanup.RegisterIPAddressCleanup(ip4b)
	cleanup.RegisterIPAddressCleanup(ip6)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4a, ip6, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4b, ip6, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
		},
	})
}

func TestAccVirtualMachinePrimaryIPResource_import(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, "", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_virtual_machine_primary_ip.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine", "primary_ip4", "primary_ip6"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_virtual_machine_primary_ip.test", "virtual_machine"),
					testutil.ReferenceFieldCheck("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					testutil.ReferenceFieldCheck("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
			{
				Config:   testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, "", false),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVirtualMachinePrimaryIPResource_externalDeletion(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4)

	testutil.RunExternalDeletionTest(t, testutil.ExternalDeletionTestConfig{
		ResourceName: "netbox_virtual_machine_primary_ip",
		Config: func() string {
			return testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, "", false)
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			client, err := testutil.GetSharedClient()
			if err != nil {
				return err
			}
			vmID64, err := strconv.ParseInt(id, 10, 32)
			if err != nil {
				return err
			}
			patch := netbox.NewPatchedWritableVirtualMachineWithConfigContextRequest()
			patch.SetPrimaryIp4Nil()
			patch.SetPrimaryIp6Nil()
			_, _, err = client.VirtualizationAPI.VirtualizationVirtualMachinesPartialUpdate(ctx, int32(vmID64)).
				PatchedWritableVirtualMachineWithConfigContextRequest(*patch).
				Execute()
			return err
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVirtualMachineDestroy,
			testutil.CheckVMInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckClusterDestroy,
			testutil.CheckClusterTypeDestroy,
		),
	})
}

func TestAccVirtualMachinePrimaryIPResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	requireIPv6Support(t)

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	vmName := testutil.RandomName("tf-test-vm")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("192.0.2.%d/24", acctest.RandIntRange(1, 254))
	ip6 := fmt.Sprintf("2001:db8:%d::1/64", acctest.RandIntRange(1, 65535))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterVMInterfaceCleanup(interfaceName, vmName)
	cleanup.RegisterIPAddressCleanup(ip4)
	cleanup.RegisterIPAddressCleanup(ip6)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, ip6, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, ip6, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
			{
				Config: testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, ip6, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_virtual_machine_primary_ip.test", "primary_ip6"),
				),
			},
		},
	})
}

func testAccVirtualMachinePrimaryIPResourceConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, ip6 string, setPrimaryIP6 bool) string {
	primaryIP6Resource := ""
	primaryIP6Attr := ""
	if ip6 != "" {
		primaryIP6Resource = fmt.Sprintf(`
resource "netbox_ip_address" "test_v6" {
  address = %q
  status  = "active"
  assigned_object_type = "virtualization.vminterface"
  assigned_object_id   = netbox_vm_interface.test.id
}
`, ip6)
		if setPrimaryIP6 {
			primaryIP6Attr = "\n  primary_ip6 = netbox_ip_address.test_v6.id"
		}
	}

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
  status  = "active"
}

resource "netbox_vm_interface" "test" {
  name            = %q
  virtual_machine = netbox_virtual_machine.test.name
}

resource "netbox_ip_address" "test_v4" {
  address              = %q
  status               = "active"
  assigned_object_type = "virtualization.vminterface"
  assigned_object_id   = netbox_vm_interface.test.id
}
%s
resource "netbox_virtual_machine_primary_ip" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  primary_ip4     = netbox_ip_address.test_v4.id%s
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, interfaceName, ip4, primaryIP6Resource, primaryIP6Attr)
}

func requireIPv6Support(t *testing.T) {
	t.Helper()
	if !supportsIPv6(t) {
		t.Skip("IPv6 IP address creation not supported by NetBox instance")
	}
}

func supportsIPv6(t *testing.T) bool {
	t.Helper()
	client, err := testutil.GetSharedClient()
	if err != nil {
		t.Fatalf("Failed to get shared client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	statusInfo, _, statusErr := client.StatusAPI.StatusRetrieve(ctx).Execute()
	if statusErr == nil && statusInfo != nil {
		if versionValue, ok := statusInfo["netbox_version"]; ok {
			if strings.HasPrefix(strings.TrimSpace(fmt.Sprint(versionValue)), "4.1.11") {
				return false
			}
		} else if versionValue, ok := statusInfo["netbox-version"]; ok {
			if strings.HasPrefix(strings.TrimSpace(fmt.Sprint(versionValue)), "4.1.11") {
				return false
			}
		}
	}

	clusterTypeName := testutil.RandomName("tf-test-ct-ipv6")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-ipv6")
	clusterName := testutil.RandomName("tf-test-cluster-ipv6")
	vmName := testutil.RandomName("tf-test-vm-ipv6")
	interfaceName := testutil.RandomName("eth")

	clusterTypeReq := netbox.NewClusterTypeRequest(clusterTypeName, clusterTypeSlug)
	clusterType, _, err := client.VirtualizationAPI.VirtualizationClusterTypesCreate(ctx).ClusterTypeRequest(*clusterTypeReq).Execute()
	if err != nil || clusterType == nil {
		return false
	}
	defer func() {
		_, _ = client.VirtualizationAPI.VirtualizationClusterTypesDestroy(ctx, clusterType.GetId()).Execute()
	}()

	clusterReq := netbox.NewWritableClusterRequest(clusterName, *netbox.NewBriefClusterTypeRequest(clusterTypeName, clusterTypeSlug))
	cluster, _, err := client.VirtualizationAPI.VirtualizationClustersCreate(ctx).WritableClusterRequest(*clusterReq).Execute()
	if err != nil || cluster == nil {
		return false
	}
	defer func() {
		_, _ = client.VirtualizationAPI.VirtualizationClustersDestroy(ctx, cluster.GetId()).Execute()
	}()

	vmReq := netbox.NewWritableVirtualMachineWithConfigContextRequest(vmName)
	status := netbox.MODULESTATUSVALUE_ACTIVE
	vmReq.Status = &status
	vmReq.Cluster = *netbox.NewNullableBriefClusterRequest(netbox.NewBriefClusterRequest(clusterName))
	vm, _, err := client.VirtualizationAPI.VirtualizationVirtualMachinesCreate(ctx).WritableVirtualMachineWithConfigContextRequest(*vmReq).Execute()
	if err != nil || vm == nil {
		return false
	}
	defer func() {
		_, _ = client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(ctx, vm.GetId()).Execute()
	}()

	ifaceReq := netbox.NewWritableVMInterfaceRequest(*netbox.NewBriefVirtualMachineRequest(vmName), interfaceName)
	iface, _, err := client.VirtualizationAPI.VirtualizationInterfacesCreate(ctx).WritableVMInterfaceRequest(*ifaceReq).Execute()
	if err != nil || iface == nil {
		return false
	}
	defer func() {
		_, _ = client.VirtualizationAPI.VirtualizationInterfacesDestroy(ctx, iface.GetId()).Execute()
	}()

	address := fmt.Sprintf("2001:db8:%d::1/64", acctest.RandIntRange(1, 65535))
	ipReq := netbox.NewWritableIPAddressRequest(address)
	ipStatus := netbox.PATCHEDWRITABLEIPADDRESSREQUESTSTATUS_ACTIVE
	ipReq.Status = &ipStatus
	ipReq.AssignedObjectType = *netbox.NewNullableString(netbox.PtrString("virtualization.vminterface"))
	ipReq.AssignedObjectId = *netbox.NewNullableInt64(netbox.PtrInt64(int64(iface.GetId())))

	ip, _, err := client.IpamAPI.IpamIpAddressesCreate(ctx).WritableIPAddressRequest(*ipReq).Execute()
	if err != nil || ip == nil {
		return false
	}

	_, _ = client.IpamAPI.IpamIpAddressesDestroy(ctx, ip.GetId()).Execute()
	return true
}
