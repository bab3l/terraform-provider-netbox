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

// Acceptance Tests.

func TestAccASNRangeResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-asn-range")
	slug := testutil.RandomSlug("tf-test-asn-range")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-asn-range-full")
	slug := testutil.RandomSlug("tf-test-asn-range-full")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "65000"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "65100"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Test ASN range with full options"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-asn-range-upd")
	slug := testutil.RandomSlug("tf-test-asn-range-upd")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
				),
			},
			{
				Config: testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64700"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-ext-del")
	slug := testutil.RandomSlug("tf-test-asn-range-ext-del")
	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamAsnRangesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find ASN range for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAsnRangesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete ASN range: %v", err)
					}
					t.Logf("Successfully externally deleted ASN range with ID: %d", itemID)
				},
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asn-range-id")
	slug := testutil.RandomSlug("tf-test-asn-range-id")
	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
				),
			},
		},
	})

}

func testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = "64512"
  end   = "64612"
}
`, rirName, rirSlug, name, slug)
}

func testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "65000"
  end         = "65100"
  tenant      = netbox_tenant.test.id
  description = "Test ASN range with full options"
}
`, rirName, rirSlug, tenantName, tenantSlug, name, slug)
}

func testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "64512"
  end         = "64700"
  description = "Updated description"
}
`, rirName, rirSlug, name, slug)
}

func TestAccASNRangeResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-asn-range")
	slug := testutil.RandomSlug("tf-test-asn-range")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(name)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
			{
				ResourceName:            "netbox_asn_range.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
			{
				Config:   testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_ASNRange(t *testing.T) {
	t.Parallel()

	rangeName := testutil.RandomName("asn-range")
	rangeSlug := testutil.RandomSlug("asn-range")
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(rangeName)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_asn_range" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  rir = netbox_rir.test.id
  tenant = netbox_tenant.test.name
  start = 65000
  end = 65100
}
`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)
}

// TestAccConsistency_ASNRange_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ASNRange_LiteralNames(t *testing.T) {
	t.Parallel()

	rangeName := testutil.RandomName("asn-range")
	rangeSlug := testutil.RandomSlug("asn-range")
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(rangeName)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "rir", rirSlug),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNRangeConsistencyLiteralNamesConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_asn_range" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  # Use literal string names to mimic existing user state
  rir = "%[4]s"
  tenant = "%[5]s"
  start = 65000
  end = 65100
  depends_on = [netbox_rir.test, netbox_tenant.test]
}
`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)
}

func TestAccASNRangeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	rangeName := testutil.RandomName("asn_range")
	rangeSlug := testutil.RandomSlug("asn_range")
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
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
				Config: testAccASNRangeResourceImportConfig_full(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_asn_range.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_asn_range.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccASNRangeResourceImportConfig_full(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_asn_range.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir", "custom_fields", "tags", "tenant"},
			}, {
				Config:   testAccASNRangeResourceImportConfig_full(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			}},
	})
}

func testAccASNRangeResourceImportConfig_full(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.asnrange"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.asnrange"]
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
resource "netbox_asn_range" "test" {
  name    = %q
  slug    = %q
  rir     = netbox_rir.test.slug
  tenant  = netbox_tenant.test.slug
  start   = 64512
  end     = 64600

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
		rirName, rirSlug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		rangeName, rangeSlug,
	)
}
