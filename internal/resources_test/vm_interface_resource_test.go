package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVMInterfaceResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVMInterfaceResource()

	if r == nil {

		t.Fatal("Expected non-nil VM Interface resource")

	}

}

func TestVMInterfaceResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVMInterfaceResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"virtual_machine", "name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"enabled", "mtu", "mac_address", "description", "mode"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVMInterfaceResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVMInterfaceResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_vm_interface"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVMInterfaceResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVMInterfaceResource().(*resources.VMInterfaceResource)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccVMInterfaceResource_basic(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")

	clusterName := testutil.RandomName("tf-test-cluster")

	vmName := testutil.RandomName("tf-test-vm")

	ifaceName := interfaceName

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

func TestAccVMInterfaceResource_update(t *testing.T) {

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

func testAccVMInterfaceResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName string) string {

	return fmt.Sprintf(`































































resource "netbox_cluster_type" "test" {































































  name = %q































































  slug = %q































































}































































































































resource "netbox_cluster" "test" {































































  name = %q































































  type = netbox_cluster_type.test.slug































































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































































  type = netbox_cluster_type.test.slug































































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

func TestAccVMInterfaceResource_import(t *testing.T) {

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

				ResourceName: "netbox_vm_interface.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_VMInterface(t *testing.T) {

	t.Parallel()

	vmName := testutil.RandomName("vm")

	clusterName := testutil.RandomName("cluster")

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	const interfaceName = "eth0"

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
