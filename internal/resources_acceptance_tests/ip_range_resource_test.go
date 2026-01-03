package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeResource_basic(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
		},
	})

}

func TestAccIPRangeResource_full(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)
	description := testutil.RandomName("ip-range-desc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_full(startAddress, endAddress, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})

}

func TestAccIPRangeResource_update(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(101, 150)
	startOctet2 := 10 + acctest.RandIntRange(1, 200)
	endOctet2 := startOctet2 + 10
	startAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet2)
	endAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet2)
	description := testutil.RandomName("ip-range-desc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress2, endAddress2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress2),
				),
			},

			{

				Config: testAccIPRangeResourceConfig_full(startAddress2, endAddress2, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})

}

func TestAccIPRangeResource_import(t *testing.T) {

	t.Parallel()

	secondOctet := acctest.RandIntRange(151, 200)
	thirdOctet := acctest.RandIntRange(151, 200)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, endOctet)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ip_range.test",

				ImportState: true,

				ImportStateVerify: true,
			},
			{
				Config:   testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				PlanOnly: true,
			},
		},
	})

}

func TestAccIPRangeResource_external_deletion(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(201, 250)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamIpRangesList(context.Background()).StartAddress([]string{startAddress}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IP range for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamIpRangesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IP range: %v", err)
					}
					t.Logf("Successfully externally deleted IP range with ID: %d", itemID)
				},
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
		},
	})
}

func TestAccIPRangeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(101, 150)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
		},
	})

}

func testAccIPRangeResourceConfig_basic(startAddress, endAddress string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %[1]q

  end_address   = %[2]q

}

`, startAddress, endAddress)

}

func testAccIPRangeResourceConfig_full(startAddress, endAddress, description string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %[1]q

  end_address   = %[2]q

  status        = "active"

  description   = %[3]q

}

`, startAddress, endAddress, description)

}

func TestAccConsistency_IPRange_LiteralNames(t *testing.T) {
	t.Parallel()
	startOctet := 50 + acctest.RandIntRange(1, 100)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("172.16.%d.10", startOctet)
	endAddress := fmt.Sprintf("172.16.%d.20", endOctet)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
			{
				Config:   testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
		},
	})
}

func testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %q
  end_address   = %q
}
`, startAddress, endAddress)
}

func TestAccIPRangeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	startAddress := "192.0.2.1"
	endAddress := "192.0.2.10"
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_ip_range.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_ip_range.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "tenant", "start_address", "end_address"},
			},
		},
	})
}

func testAccIPRangeResourceImportConfig_full(startAddress, endAddress, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.iprange"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.iprange"]
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
resource "netbox_ip_range" "test" {
  start_address = %q
  end_address   = %q
  tenant        = netbox_tenant.test.slug

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
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		startAddress, endAddress,
	)
}
