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

func TestVirtualDiskResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	if r == nil {

		t.Fatal("Expected non-nil VirtualDisk resource")
	}
}

func TestVirtualDiskResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"virtual_machine", "name", "size"}

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

	optionalAttrs := []string{"description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestVirtualDiskResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_virtual_disk"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestVirtualDiskResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualDiskResource()

	// Type assert to access Configure method

	configurable, ok := r.(fwresource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")
	}

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)
	}
}

// Acceptance Tests

func TestAccVirtualDiskResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	diskName := testutil.RandomName("tf-test-disk")

	vmName := testutil.RandomName("tf-test-vm")

	clusterName := testutil.RandomName("tf-test-cluster")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(diskName)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskResourceConfig_basic(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "size", "100"),

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "virtual_machine"),
				),
			},
		},
	})
}

func TestAccVirtualDiskResource_full(t *testing.T) {

	// Generate unique names

	diskName := testutil.RandomName("tf-test-disk-full")

	vmName := testutil.RandomName("tf-test-vm-full")

	clusterName := testutil.RandomName("tf-test-cluster-full")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(diskName)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskResourceConfig_full(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "size", "500"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "description", "Test virtual disk with full options"),

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "virtual_machine"),
				),
			},
		},
	})
}

func TestAccVirtualDiskResource_update(t *testing.T) {

	// Generate unique names

	diskName := testutil.RandomName("tf-test-disk-upd")

	vmName := testutil.RandomName("tf-test-vm-upd")

	clusterName := testutil.RandomName("tf-test-cluster-upd")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(diskName)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskResourceConfig_basic(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "size", "100"),
				),
			},

			{

				Config: testAccVirtualDiskResourceConfig_updated(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "size", "200"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccVirtualDiskResourceConfig_basic(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {

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

  # Ignore disk changes since Netbox auto-computes this from virtual_disks

  lifecycle {

    ignore_changes = [disk]
  }
}

resource "netbox_virtual_disk" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q

  size            = 100
}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName)
}

func testAccVirtualDiskResourceConfig_full(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {

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

  # Ignore disk changes since Netbox auto-computes this from virtual_disks

  lifecycle {

    ignore_changes = [disk]
  }
}

resource "netbox_virtual_disk" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q

  size            = 500
  description     = "Test virtual disk with full options"
}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName)
}

func testAccVirtualDiskResourceConfig_updated(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug string) string {

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

  # Ignore disk changes since Netbox auto-computes this from virtual_disks

  lifecycle {

    ignore_changes = [disk]
  }
}

resource "netbox_virtual_disk" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q

  size            = 200
  description     = "Updated description"
}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName)
}

func TestAccVirtualDiskResource_import(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	diskName := testutil.RandomName("tf-test-disk")

	vmName := testutil.RandomName("tf-test-vm")

	clusterName := testutil.RandomName("tf-test-cluster")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(diskName)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskResourceConfig_basic(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "size", "100"),

					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "virtual_machine"),
				),
			},

			{

				ResourceName: "netbox_virtual_disk.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_VirtualDisk(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	clusterName := testutil.RandomName("cluster")

	vmName := testutil.RandomName("vm")

	diskName := testutil.RandomName("disk")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskConsistencyConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "virtual_machine", vmName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccVirtualDiskConsistencyConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName),
			},
		},
	})
}

func testAccVirtualDiskConsistencyConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_virtual_machine" "test" {
  name = "%[4]s"
  cluster = netbox_cluster.test.id

  lifecycle {

    ignore_changes = [disk]
  }
}

resource "netbox_virtual_disk" "test" {
  virtual_machine = netbox_virtual_machine.test.name
  name = "%[5]s"

  size = 100
}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName)
}
