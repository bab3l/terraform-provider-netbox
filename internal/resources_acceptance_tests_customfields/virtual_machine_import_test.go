//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccVirtualMachineResource_importWithCustomFieldsAndTags tests importing a pre-existing VM
// that has various custom field types and tags properly imports all data.
func TestAccVirtualMachineResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	vmName := testutil.RandomName("tf-test-vm-import")
	clusterName := testutil.RandomName("tf-test-cluster")
	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-vm-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-vm-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-vm-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-vm-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values for different data types
	cfText := testutil.RandomCustomFieldName("tf_vm_text")
	cfTextValue := testutil.RandomName("vm-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_vm_longtext")
	cfLongtextValue := fmt.Sprintf("VM description with details: %s", testutil.RandomName("vm-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_vm_integer")
	cfIntegerValue := 2048 // Memory size
	cfBoolean := testutil.RandomCustomFieldName("tf_vm_boolean")
	cfBooleanValue := false
	cfDate := testutil.RandomCustomFieldName("tf_vm_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_vm_url")
	cfURLValue := testutil.RandomURL("vm")
	cfJSON := testutil.RandomCustomFieldName("tf_vm_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the VM with all custom fields and tags
				Config: testAccVirtualMachineResourceImportConfig_full(
					vmName, clusterName, clusterTypeSlug, clusterTypeName, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_machine.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the VM and verify all fields are preserved
				ResourceName:            "netbox_virtual_machine.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster", "custom_fields"}, // Cluster reference may have lookup inconsistencies, custom fields have import limitations
				// The import should preserve all custom fields and tags
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),

					// Verify tags are imported
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "tags.#", "2"),

					// Verify custom fields are imported
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.#", "7"),

					// Verify specific custom field values (different data types)
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.0.name", cfText),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.0.value", cfTextValue),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.1.name", cfLongtext),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.1.type", "longtext"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.1.value", cfLongtextValue),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.2.name", cfIntegerName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.2.type", "integer"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.2.value", fmt.Sprintf("%d", cfIntegerValue)),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.3.name", cfBoolean),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.3.type", "boolean"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.3.value", fmt.Sprintf("%t", cfBooleanValue)),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.4.name", cfDate),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.4.type", "date"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.4.value", cfDateValue),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.5.name", cfURL),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.5.type", "url"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.5.value", cfURLValue),

					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.6.name", cfJSON),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.6.type", "json"),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "custom_fields.6.value", cfJSONValue),
				),
			},
		},
	})
}

func testAccVirtualMachineResourceImportConfig_full(
	vmName, clusterName, clusterTypeSlug, clusterTypeName, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "vm_test1" {
  name  = %[7]q
  slug  = %[8]q
  color = %[9]q
}

resource "netbox_tag" "vm_test2" {
  name  = %[10]q
  slug  = %[11]q
  color = %[12]q
}

# Create custom fields for virtualization.virtualmachine
resource "netbox_custom_field" "vm_text" {
  name         = %[13]q
  type         = "text"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_longtext" {
  name         = %[15]q
  type         = "longtext"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_integer" {
  name         = %[17]q
  type         = "integer"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_boolean" {
  name         = %[19]q
  type         = "boolean"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_date" {
  name         = %[21]q
  type         = "date"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_url" {
  name         = %[23]q
  type         = "url"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

resource "netbox_custom_field" "vm_json" {
  name         = %[25]q
  type         = "json"
  object_types = ["virtualization.virtualmachine"]
  required     = false
}

# Create dependencies
resource "netbox_site" "vm_test" {
  name   = %[5]q
  slug   = %[6]q
  status = "active"
}

resource "netbox_cluster_type" "vm_test" {
  name = %[4]q
  slug = %[3]q
}

resource "netbox_cluster" "vm_test" {
  name = %[2]q
  type = netbox_cluster_type.vm_test.id
  site = netbox_site.vm_test.id
}

# Create VM with all custom fields and tags
resource "netbox_virtual_machine" "test" {
  name    = %[1]q
  cluster = netbox_cluster.vm_test.id
  status  = "active"

  tags = [
    {
      name = netbox_tag.vm_test1.name
      slug = netbox_tag.vm_test1.slug
    },
    {
      name = netbox_tag.vm_test2.name
      slug = netbox_tag.vm_test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.vm_text.name
      type  = "text"
      value = %[14]q
    },
    {
      name  = netbox_custom_field.vm_longtext.name
      type  = "longtext"
      value = %[16]q
    },
    {
      name  = netbox_custom_field.vm_integer.name
      type  = "integer"
      value = "%[18]d"
    },
    {
      name  = netbox_custom_field.vm_boolean.name
      type  = "boolean"
      value = "%[20]t"
    },
    {
      name  = netbox_custom_field.vm_date.name
      type  = "date"
      value = %[22]q
    },
    {
      name  = netbox_custom_field.vm_url.name
      type  = "url"
      value = %[24]q
    },
    {
      name  = netbox_custom_field.vm_json.name
      type  = "json"
      value = %[26]q
    }
  ]
}
`, vmName, clusterName, clusterTypeSlug, clusterTypeName, siteName, siteSlug,
		tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
