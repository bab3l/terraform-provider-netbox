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

	ifaceName := "eth0"

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

	ifaceName := "eth0"

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
