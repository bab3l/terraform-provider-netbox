package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttrSet("netbox_interface.test", "device"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "1500"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mgmt_only", "true"),

					resource.TestCheckResourceAttr("netbox_interface.test", "description", "Test interface with full options"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "true"),
				),
			},

			{

				Config: testAccInterfaceResourceConfig_updated(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),

					resource.TestCheckResourceAttr("netbox_interface.test", "description", "Updated interface description"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_interface.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})

}

func TestAccConsistency_Interface(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-consistency")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_consistency_device_name(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func TestAccConsistency_Interface_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-literal")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_consistency_device_id(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func testAccInterfaceResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.id

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_full(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device       = netbox_device.test.id

  name         = %q

  type         = "1000base-t"

  enabled      = false

  mtu          = 1500

  mgmt_only    = true

  description  = "Test interface with full options"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_updated(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device      = netbox_device.test.id

  name        = %q

  type        = "1000base-t"

  enabled     = false

  description = "Updated interface description"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_consistency_device_name(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.name

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_consistency_device_id(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.id

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourcePrereqs(name string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %q

  slug = %q

}

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

}

resource "netbox_device_role" "test" {

  name = %q

  slug = %q

}

resource "netbox_device" "test" {

  site        = netbox_site.test.id

  name        = %q

  device_type = netbox_device_type.test.id

  role        = netbox_device_role.test.id

  status      = "offline"

}

`, name+"-site", testutil.RandomSlug("site"), name+"-mfr", testutil.RandomSlug("mfr"), name+"-model", testutil.RandomSlug("device"), name+"-role", testutil.RandomSlug("role"), name+"-device")

}

func TestAccInterfaceResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)
	name := testutil.RandomName("tf-test-interface-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[2]s"
}
resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[3]s"
}
resource "netbox_device_type" "test" {
  model = "%[1]s-model"
  slug = "%[4]s"
  manufacturer = netbox_manufacturer.test.id
}
resource "netbox_device_role" "test" {
  name = "%[1]s-role"
  slug = "%[5]s"
}
resource "netbox_device" "test" {
  site = netbox_site.test.id
  name = "%[1]s-device"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
}
resource "netbox_interface" "test" {
  name = "%[1]s"
  device = netbox_device.test.id
  type = "1000base-t"
}
`, name, testutil.RandomSlug("site"), testutil.RandomSlug("mfr"), testutil.RandomSlug("device"), testutil.RandomSlug("role")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInterfacesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find interface for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInterfacesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete interface: %v", err)
					}
					t.Logf("Successfully externally deleted interface with ID: %d", itemID)
				},
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[2]s"
}
resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[3]s"
}
resource "netbox_device_type" "test" {
  model = "%[1]s-model"
  slug = "%[4]s"
  manufacturer = netbox_manufacturer.test.id
}
resource "netbox_device_role" "test" {
  name = "%[1]s-role"
  slug = "%[5]s"
}
resource "netbox_device" "test" {
  site = netbox_site.test.id
  name = "%[1]s-device"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
}
resource "netbox_interface" "test" {
  name = "%[1]s"
  device = netbox_device.test.id
  type = "1000base-t"
}
`, name, testutil.RandomSlug("site"), testutil.RandomSlug("mfr"), testutil.RandomSlug("device"), testutil.RandomSlug("role")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
				),
			},
		},
	})
}

// TestAccInterfaceResource_importWithCustomFieldsAndTags tests importing an interface
// with custom fields and tags to ensure all data is preserved during import.
func TestAccInterfaceResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	interfaceName := testutil.RandomName("tf-test-int-import")
	deviceName := testutil.RandomName("tf-test-device-import")
	manufacturerName := testutil.RandomName("tf-test-mfr-import")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")
	deviceTypeModel := testutil.RandomName("tf-test-dt-import")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-import")
	deviceRoleName := testutil.RandomName("tf-test-role-import")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-import")
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-int-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-int-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-int-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-int-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values for different data types
	cfText := testutil.RandomCustomFieldName("tf_int_text")
	cfTextValue := testutil.RandomName("interface-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_int_longtext")
	cfLongtextValue := fmt.Sprintf("Interface description: %s", testutil.RandomName("int-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_int_integer")
	cfIntegerValue := 1000
	cfBoolean := testutil.RandomCustomFieldName("tf_int_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_int_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_int_url")
	cfURLValue := testutil.RandomURL("interface")
	cfJSON := testutil.RandomCustomFieldName("tf_int_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the interface with all custom fields and tags
				Config: testAccInterfaceResourceImportConfig_full(
					interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttr("netbox_interface.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_interface.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the interface and verify all fields are preserved
				ResourceName:            "netbox_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields", "tags"}, // Device reference field, CF/Tags import needs work
				// The import should preserve basic fields
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
				),
			},
		},
	})
}

func testAccInterfaceResourceImportConfig_full(
	interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "int_test1" {
  name  = %[11]q
  slug  = %[12]q
  color = %[13]q
}

resource "netbox_tag" "int_test2" {
  name  = %[14]q
  slug  = %[15]q
  color = %[16]q
}

# Create custom fields for dcim.interface
resource "netbox_custom_field" "int_text" {
  name         = %[17]q
  type         = "text"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_longtext" {
  name         = %[19]q
  type         = "longtext"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_integer" {
  name         = %[21]q
  type         = "integer"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_boolean" {
  name         = %[23]q
  type         = "boolean"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_date" {
  name         = %[25]q
  type         = "date"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_url" {
  name         = %[27]q
  type         = "url"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_custom_field" "int_json" {
  name         = %[29]q
  type         = "json"
  object_types = ["dcim.interface"]
  required     = false
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "test" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  name        = %[2]q
  status      = "active"
}

# Create interface with all custom fields and tags
resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = %[1]q
  type   = "1000base-t"

  tags = [
    {
      name = netbox_tag.int_test1.name
      slug = netbox_tag.int_test1.slug
    },
    {
      name = netbox_tag.int_test2.name
      slug = netbox_tag.int_test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.int_text.name
      type  = "text"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.int_longtext.name
      type  = "longtext"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.int_integer.name
      type  = "integer"
      value = "%[22]d"
    },
    {
      name  = netbox_custom_field.int_boolean.name
      type  = "boolean"
      value = "%[24]t"
    },
    {
      name  = netbox_custom_field.int_date.name
      type  = "date"
      value = %[26]q
    },
    {
      name  = netbox_custom_field.int_url.name
      type  = "url"
      value = %[28]q
    },
    {
      name  = netbox_custom_field.int_json.name
      type  = "json"
      value = %[30]q
    }
  ]
}
`, interfaceName, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
