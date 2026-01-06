//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
