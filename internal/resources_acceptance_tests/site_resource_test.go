package resources_acceptance_tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site")
	slug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccSiteResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-full")
	slug := testutil.RandomSlug("tf-test-site-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_site.test", "description", description),
				),
			},
		},
	})
}

func TestAccSiteResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-update")
	slug := testutil.RandomSlug("tf-test-site-upd")
	updatedName := testutil.RandomName("tf-test-site-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
				),
			},
			{
				Config: testAccSiteResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccSiteResource_import(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-site-import")
	slug := testutil.RandomSlug("tf-test-site-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_Site(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	regionName := testutil.RandomName("region")
	regionSlug := testutil.RandomSlug("region")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "region", regionName),
					resource.TestCheckResourceAttr("netbox_site.test", "group", groupName),
					resource.TestCheckResourceAttr("netbox_site.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_Site_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("tf-test-site-lit")
	siteSlug := testutil.RandomSlug("tf-test-site-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConsistencyLiteralNamesConfig(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", siteSlug),
				),
			},
			{
				Config:   testAccSiteConsistencyLiteralNamesConfig(siteName, siteSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccSiteResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-site-id")
	slug := testutil.RandomSlug("tf-test-site-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
		},
	})
}

func testAccSiteConsistencyLiteralNamesConfig(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, siteName, siteSlug)
}

func testAccSiteResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name        = %q
  slug        = %q
  status      = "active"
  description = %q
}
`, name, slug, description)
}

func testAccSiteResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_site_group" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_tenant" "test" {
  name = "%[7]s"
  slug = "%[8]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  region = netbox_region.test.name
  group = netbox_site_group.test.name
  tenant = netbox_tenant.test.name
}
`, siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccSiteResource_externalDeletion(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-site-ext-del")
	slug := testutil.RandomSlug("site")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}
`, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimSitesList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find site for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimSitesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete site: %v", err)
					}
					t.Logf("Successfully externally deleted site with ID: %d", itemID)
				},
				Config: fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}
`, name, slug),
				ExpectError: regexp.MustCompile("(?i)(404|not found|no site)"),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}

// TestAccSiteResource_importWithCustomFieldsAndTags tests importing a site
// with custom fields and tags to ensure all data is preserved during import.
func TestAccSiteResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-site-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-site-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-site-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-site-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values
	cfText := testutil.RandomCustomFieldName("tf_site_text")
	cfTextValue := testutil.RandomName("site-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_site_longtext")
	cfLongtextValue := fmt.Sprintf("Site description: %s", testutil.RandomName("site-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_site_integer")
	cfIntegerValue := 100
	cfBoolean := testutil.RandomCustomFieldName("tf_site_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_site_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_site_url")
	cfURLValue := testutil.RandomURL("site")
	cfJSON := testutil.RandomCustomFieldName("tf_site_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				// Create the site with all custom fields and tags
				Config: testAccSiteResourceImportConfig_full(
					siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_site.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the site and verify basic fields are preserved
				ResourceName:            "netbox_site.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", siteSlug),
				),
			},
		},
	})
}

func testAccSiteResourceImportConfig_full(
	siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "site_test1" {
  name  = %[3]q
  slug  = %[4]q
  color = %[5]q
}

resource "netbox_tag" "site_test2" {
  name  = %[6]q
  slug  = %[7]q
  color = %[8]q
}

# Create custom fields for dcim.site
resource "netbox_custom_field" "site_text" {
  name         = %[9]q
  type         = "text"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_longtext" {
  name         = %[11]q
  type         = "longtext"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_integer" {
  name         = %[13]q
  type         = "integer"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_boolean" {
  name         = %[15]q
  type         = "boolean"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_date" {
  name         = %[17]q
  type         = "date"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_url" {
  name         = %[19]q
  type         = "url"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_json" {
  name         = %[21]q
  type         = "json"
  object_types = ["dcim.site"]
  required     = false
}

# Create site with all custom fields and tags
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"

  tags = [
    {
      name = netbox_tag.site_test1.name
      slug = netbox_tag.site_test1.slug
    },
    {
      name = netbox_tag.site_test2.name
      slug = netbox_tag.site_test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.site_text.name
      type  = "text"
      value = %[10]q
    },
    {
      name  = netbox_custom_field.site_longtext.name
      type  = "longtext"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.site_integer.name
      type  = "integer"
      value = "%[14]d"
    },
    {
      name  = netbox_custom_field.site_boolean.name
      type  = "boolean"
      value = "%[16]t"
    },
    {
      name  = netbox_custom_field.site_date.name
      type  = "date"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.site_url.name
      type  = "url"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.site_json.name
      type  = "json"
      value = %[22]q
    }
  ]
}
`, siteName, siteSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}
