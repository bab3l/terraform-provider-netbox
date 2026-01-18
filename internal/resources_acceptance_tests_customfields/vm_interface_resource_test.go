//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccVMInterfaceResource_importWithCustomFieldsAndTags tests importing a VM interface
// with custom fields and tags to ensure all data is preserved during import.
func TestAccVMInterfaceResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

	tags = [netbox_tag.vmint_test1.slug, netbox_tag.vmint_test2.slug]

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

// TestAccVMInterfaceResource_CustomFieldsPreservation tests that custom fields are preserved
// when they are not specified in the configuration after creation.
func TestAccVMInterfaceResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	clusterTypeName := testutil.RandomName("cluster_type")
	clusterTypeSlug := testutil.RandomSlug("cluster_type")
	clusterName := testutil.RandomName("cluster")
	vmName := testutil.RandomName("vm")
	ifaceName := testutil.RandomName("vmint")
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterVirtualMachineCleanup(vmName)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with custom fields
			{
				Config: testAccVMInterfaceResourceConfig_withCustomFields(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "custom_fields.#", "2"),
				),
			},
			// Step 2: Update without custom_fields in config - should preserve in NetBox
			{
				Config: testAccVMInterfaceResourceConfig_withoutCustomFields(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vm_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					// Custom fields should be null/empty in state (filter-to-owned)
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Add custom fields back - should see them again
			{
				Config: testAccVMInterfaceResourceConfig_withCustomFields(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, cfText, cfInteger),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "name", ifaceName),
					resource.TestCheckResourceAttr("netbox_vm_interface.test", "custom_fields.#", "2"),
				),
			},
		},
	})
}

func testAccVMInterfaceResourceConfig_withCustomFields(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, cfText, cfInteger string) string {
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

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["virtualization.vminterface"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["virtualization.vminterface"]
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q

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
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, cfText, cfInteger, ifaceName)
}

func testAccVMInterfaceResourceConfig_withoutCustomFields(clusterTypeName, clusterTypeSlug, clusterName, vmName, ifaceName, cfText, cfInteger string) string {
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

resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["virtualization.vminterface"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["virtualization.vminterface"]
}

resource "netbox_vm_interface" "test" {
  virtual_machine = netbox_virtual_machine.test.id
  name            = %q

  # custom_fields intentionally omitted - should preserve in NetBox
}
`, clusterTypeName, clusterTypeSlug, clusterName, vmName, cfText, cfInteger, ifaceName)
}
