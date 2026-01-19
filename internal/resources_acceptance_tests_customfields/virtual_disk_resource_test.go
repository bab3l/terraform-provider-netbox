//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDiskDestroy,
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
				ResourceName:      "netbox_virtual_disk.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccVirtualDiskResourceImportConfig_full(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
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

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
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

// TestAccVirtualDiskResource_CustomFieldsPreservation tests that custom fields are preserved
// when they are not specified in the configuration after creation.
func TestAccVirtualDiskResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	diskName := testutil.RandomName("virtual_disk")
	vmName := testutil.RandomName("vm")
	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster_type")
	clusterTypeSlug := testutil.RandomSlug("cluster_type")
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVirtualDiskCleanup(diskName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVirtualDiskDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with custom fields
			{
				Config: testAccVirtualDiskResourceConfig_withCustomFields(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "custom_fields.#", "2"),
				),
			},
			// Step 2: Update without custom_fields in config - should preserve in NetBox
			{
				Config: testAccVirtualDiskResourceConfig_withoutCustomFields(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_disk.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),
					// Custom fields should be null/empty in state (filter-to-owned)
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Add custom fields back - should see them again
			{
				Config: testAccVirtualDiskResourceConfig_withCustomFields(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "name", diskName),
					resource.TestCheckResourceAttr("netbox_virtual_disk.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

func testAccVirtualDiskResourceConfig_withCustomFields(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, cfText, cfInteger string) string {
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

  lifecycle {
    ignore_changes = [disk]
  }
}

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["virtualization.virtualdisk"]
}

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
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    }
  ]
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, cfText, cfInteger, diskName)
}

func testAccVirtualDiskResourceConfig_withoutCustomFields(diskName, vmName, clusterName, clusterTypeName, clusterTypeSlug, cfText, cfInteger string) string {
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

  lifecycle {
    ignore_changes = [disk]
  }
}

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["virtualization.virtualdisk"]
}

resource "netbox_virtual_disk" "test" {
  name         = %q
  virtual_machine = netbox_virtual_machine.test.name
  size         = 100

  # custom_fields intentionally omitted - should preserve in NetBox
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, cfText, cfInteger, diskName)
}
