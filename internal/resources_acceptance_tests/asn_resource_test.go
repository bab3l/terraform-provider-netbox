package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNResource_basic(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Generate a random ASN in the private range (64512-65534)

	asn := int64(acctest.RandIntRange(64512, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},

			{

				ResourceName: "netbox_asn.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})

}

func TestAccASNResource_full(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate a random ASN in the private range (64512-65534)

	asn := int64(acctest.RandIntRange(64512, 65534))

	description := testutil.RandomName("description")

	updatedDescription := "Updated ASN description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, description, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),

					resource.TestCheckResourceAttr("netbox_asn.test", "description", description),

					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),

					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),

					resource.TestCheckResourceAttrSet("netbox_asn.test", "tenant"),
				),
			},

			{

				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, updatedDescription, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccASNResource_IDPreservation(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")
	asn := int64(acctest.RandIntRange(64512, 65000))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
		},
	})

}

func TestAccASNResource_update(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	asn := int64(acctest.RandIntRange(64512, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccASNResource_external_deletion(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")
	asn := int64(acctest.RandIntRange(64512, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// Safe conversion with bounds check
					if asn < 0 || asn > 4294967295 { // Max ASN value (32-bit)
						t.Fatalf("ASN value %d out of valid range", asn)
					}
					//nolint:gosec // Safe conversion - bounds checked above
					items, _, err := client.IpamAPI.IpamAsnsList(context.Background()).Asn([]int32{int32(asn)}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find asn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAsnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete asn: %v", err)
					}
					t.Logf("Successfully externally deleted asn with ID: %d", itemID)
				},
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
				),
			},
		},
	})
}

func testAccASNResourceConfig_basic(rirName, rirSlug string, asn int64) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

resource "netbox_asn" "test" {

  asn = %d

  rir = netbox_rir.test.id

}

`, rirName, rirSlug, asn)

}

func testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug string, asn int64, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  tenant      = netbox_tenant.test.id
  description = %q
  comments    = %q
}

`, rirName, rirSlug, tenantName, tenantSlug, asn, description, comments)

}

func testAccASNResourceConfig_update(rirName, rirSlug string, asn int64, description string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

resource "netbox_asn" "test" {

  asn         = %d

  rir         = netbox_rir.test.id

  description = %q

}

`, rirName, rirSlug, asn, description)

}

// TestAccConsistency_ASN_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_ASN_LiteralNames(t *testing.T) {

	t.Parallel()
	asn := int64(65100)

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),

					resource.TestCheckResourceAttr("netbox_asn.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_asn.test", "tenant", tenantName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccASNConsistencyLiteralNamesConfig(asn int64, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}

resource "netbox_asn" "test" {

  asn = %[1]d

  # Use literal string names to mimic existing user state

  rir = "%[3]s"

  tenant = "%[4]s"

  depends_on = [netbox_rir.test, netbox_tenant.test]

}

`, asn, rirName, rirSlug, tenantName, tenantSlug)

}

func TestAccASNResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	asn := int64(acctest.RandIntRange(64512, 65534))
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
				Config: testAccASNResourceImportConfig_full(asn, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_asn.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_asn.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccASNResourceImportConfig_full(asn, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_asn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir", "custom_fields", "tags", "tenant"},
			},
		},
	})
}

func testAccASNResourceImportConfig_full(asn int64, rirName, rirSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["ipam.asn"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["ipam.asn"]
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
resource "netbox_asn" "test" {
  asn    = %d
  rir    = netbox_rir.test.slug
  tenant = netbox_tenant.test.slug

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
		asn,
	)
}
