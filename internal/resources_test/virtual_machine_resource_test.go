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

func TestVirtualMachineResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	if r == nil {

		t.Fatal("Expected non-nil Virtual Machine resource")

	}

}

func TestVirtualMachineResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name"}

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

	optionalAttrs := []string{"status", "cluster", "vcpus", "memory", "disk", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVirtualMachineResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualMachineResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_virtual_machine"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVirtualMachineResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualMachineResource().(*resources.VirtualMachineResource)

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

func TestAccVirtualMachineResource_basic(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")

	clusterName := testutil.RandomName("tf-test-cluster")

	vmName := testutil.RandomName("tf-test-vm")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

func testAccVirtualMachineResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {

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































































  type = netbox_cluster_type.test.slug































































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

func TestAccVirtualMachineResource_import(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")

	clusterName := testutil.RandomName("tf-test-cluster")

	vmName := testutil.RandomName("tf-test-vm")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

				ResourceName: "netbox_virtual_machine.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_VirtualMachine(t *testing.T) {

	t.Parallel()

	vmName := testutil.RandomName("vm")

	clusterName := testutil.RandomName("cluster")

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	platformName := testutil.RandomName("platform")

	platformSlug := testutil.RandomSlug("platform")

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualMachineConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, siteName, siteSlug, tenantName, tenantSlug, platformName, platformSlug, manufacturerName, manufacturerSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "cluster", clusterName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "platform", platformName),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "role", roleName),
				),
			},

			{

				// Verify no drift

				PlanOnly: true,

				Config: testAccVirtualMachineConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, siteName, siteSlug, tenantName, tenantSlug, platformName, platformSlug, manufacturerName, manufacturerSlug, roleName, roleSlug),
			},
		},
	})

}

func testAccVirtualMachineConsistencyConfig(vmName, clusterName, clusterTypeName, clusterTypeSlug, siteName, siteSlug, tenantName, tenantSlug, platformName, platformSlug, manufacturerName, manufacturerSlug, roleName, roleSlug string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = "%[3]s"

  slug = "%[4]s"

}



resource "netbox_cluster" "test" {

  name = "%[2]s"

  type = netbox_cluster_type.test.id

}



resource "netbox_site" "test" {

  name = "%[5]s"

  slug = "%[6]s"

}



resource "netbox_tenant" "test" {

  name = "%[7]s"

  slug = "%[8]s"

}



resource "netbox_manufacturer" "test" {

  name = "%[11]s"

  slug = "%[12]s"

}



resource "netbox_platform" "test" {

  name = "%[9]s"

  slug = "%[10]s"

  manufacturer = netbox_manufacturer.test.id

}



resource "netbox_device_role" "test" {

  name = "%[13]s"

  slug = "%[14]s"

}



resource "netbox_virtual_machine" "test" {

  name = "%[1]s"

  cluster = netbox_cluster.test.name

  site = netbox_site.test.name

  tenant = netbox_tenant.test.name

  platform = netbox_platform.test.name

  role = netbox_device_role.test.name

}

`, vmName, clusterName, clusterTypeName, clusterTypeSlug, siteName, siteSlug, tenantName, tenantSlug, platformName, platformSlug, manufacturerName, manufacturerSlug, roleName, roleSlug)

}
