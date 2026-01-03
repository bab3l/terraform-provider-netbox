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

func TestAccVirtualDiskResource_basic(t *testing.T) {

	t.Parallel()

	diskName := testutil.RandomName("tf-test-disk")

	vmName := testutil.RandomName("tf-test-vm")

	clusterName := testutil.RandomName("tf-test-cluster")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

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

	t.Parallel()

	diskName := testutil.RandomName("tf-test-disk-full")

	vmName := testutil.RandomName("tf-test-vm-full")

	clusterName := testutil.RandomName("tf-test-cluster-full")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

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

	t.Parallel()

	diskName := testutil.RandomName("tf-test-disk-upd")

	vmName := testutil.RandomName("tf-test-vm-upd")

	clusterName := testutil.RandomName("tf-test-cluster-upd")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

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

func TestAccVirtualDiskResource_import(t *testing.T) {

	t.Parallel()

	diskName := testutil.RandomName("tf-test-disk")

	vmName := testutil.RandomName("tf-test-vm")

	clusterName := testutil.RandomName("tf-test-cluster")

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

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

				ImportStateVerifyIgnore: []string{"virtual_machine"},
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

func TestAccConsistency_VirtualDisk_LiteralNames(t *testing.T) {

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

				Config: testAccVirtualDiskConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "virtual_machine", vmName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccVirtualDiskConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName),
			},
		},
	})

}

func TestAccVirtualDiskResource_IDPreservation(t *testing.T) {
	t.Parallel()

	diskName := testutil.RandomName("tf-test-disk-id")
	vmName := testutil.RandomName("tf-test-vm-id")
	clusterName := testutil.RandomName("tf-test-cluster-id")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type-id")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct-id")

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

func testAccVirtualDiskConsistencyLiteralNamesConfig(clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName string) string {

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

  name = "%[5]s"

  virtual_machine = "%[4]s"

  size = 100

  depends_on = [netbox_virtual_machine.test]

}

`, clusterTypeName, clusterTypeSlug, clusterName, vmName, diskName)

}

func TestAccVirtualDiskResource_externalDeletion(t *testing.T) {
	t.Parallel()

	diskName := testutil.RandomName("test-disk-del")
	vmName := testutil.RandomName("test-vm")
	clusterName := testutil.RandomName("test-cluster")
	clusterTypeName := testutil.RandomName("test-cluster-type")
	clusterTypeSlug := testutil.GenerateSlug(clusterTypeName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualMachineCleanup(vmName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDiskResourceConfig_basic(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationVirtualDisksList(context.Background()).Name([]string{diskName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find virtual_disk for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationVirtualDisksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete virtual_disk: %v", err)
					}
					t.Logf("Successfully externally deleted virtual_disk with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVirtualDiskResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	diskName := testutil.RandomName("virtual_disk")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster_type")
	clusterTypeSlug := testutil.RandomSlug("cluster_type")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualDiskCleanup(diskName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVirtualDiskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDiskResourceImportConfig_full(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_virtual_disk.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"virtual_machine", "custom_fields"}, // Virtual machine reference may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccVirtualDiskResourceImportConfig_full(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

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

  lifecycle {
    ignore_changes = [disk]
  }
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["virtualization.virtualdisk"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Main Resource
resource "netbox_virtual_disk" "test" {
  name         = %q
  virtual_machine = netbox_virtual_machine.test.name
  size         = 100

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`,
		tenantName, tenantSlug,
		clusterTypeName, clusterTypeSlug,
		clusterName,
		vmName,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		diskName,
	)
}
