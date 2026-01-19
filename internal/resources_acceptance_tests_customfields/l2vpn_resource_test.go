//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2vpnResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	name := testutil.RandomName("test-l2vpn")
	slug := testutil.RandomSlug("test-l2vpn")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate test data for all custom field types
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42 // Fixed value for reproducibility
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2vpnResourceImportConfig_full(
					name, slug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				Config: testAccL2vpnResourceImportConfig_full(
					name, slug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				ResourceName:      "netbox_l2vpn.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config: testAccL2vpnResourceImportConfig_full(
					name, slug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccL2vpnResourceImportConfig_full(
	name, slug, tenantName, tenantSlug string,
	textValue, longtextValue string, intValue int, boolValue bool, dateValue, urlValue, jsonValue string,
	tag1, tag1Slug, tag2, tag2Slug string,
	cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON string,
) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
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

# Custom Fields for vpn.l2vpn object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["vpn.l2vpn"]
}

# L2VPN with comprehensive custom fields and tags
resource "netbox_l2vpn" "test" {
  name   = %q
  slug   = %q
  type   = "vxlan"
  tenant = netbox_tenant.test.id

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = %q
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "%d"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "%t"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = %q
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = %q
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = %q
    },
  ]

  depends_on = [
    netbox_custom_field.test_text,
    netbox_custom_field.test_longtext,
    netbox_custom_field.test_integer,
    netbox_custom_field.test_boolean,
    netbox_custom_field.test_date,
    netbox_custom_field.test_url,
    netbox_custom_field.test_json,
  ]
}
`, tenantName, tenantSlug,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		name, slug, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

func TestAccL2VPNResource_CustomFieldsPreservation(t *testing.T) {
	l2vpnName := testutil.RandomName("tf-test-l2vpn")
	cfName := testutil.RandomCustomFieldName("tf_l2vpn_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create L2VPN with custom field defined and populated
			{
				Config: testAccL2VPNResourcePreservationConfig_step1(l2vpnName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", l2vpnName),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_l2vpn.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update same L2VPN without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccL2VPNResourcePreservationConfig_step2(l2vpnName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", l2vpnName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Import to verify custom fields still exist in NetBox
			{
				ResourceName:            "netbox_l2vpn.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccL2VPNResourcePreservationConfig_step1(l2vpnName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", l2vpnName),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_l2vpn.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccL2VPNResourcePreservationConfig_step1(
	l2vpnName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "l2vpn_pres" {
  name = %[1]q
  object_types = ["vpn.l2vpn"]
  type = "text"
}

resource "netbox_l2vpn" "test" {
  name = %[2]q
  slug = "tf-l2vpn-pres-%[2]s"
  type = "vpls"
  custom_fields = [
    {
      name = netbox_custom_field.l2vpn_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.l2vpn_pres]
}
`, cfName, l2vpnName)
}

func testAccL2VPNResourcePreservationConfig_step2(
	l2vpnName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "l2vpn_pres" {
  name = %[1]q
  object_types = ["vpn.l2vpn"]
  type = "text"
}

resource "netbox_l2vpn" "test" {
  name = %[2]q
  slug = "tf-l2vpn-pres-%[2]s"
  type = "vpls"
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.l2vpn_pres]
}
`, cfName, l2vpnName)
}
