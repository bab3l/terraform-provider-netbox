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

func TestAccRackResource_basic(t *testing.T) {

	t.Parallel()

	// Generate unique names to avoid conflicts between test runs

	siteName := testutil.RandomName("tf-test-rack-site")

	siteSlug := testutil.RandomSlug("tf-test-rack-site")

	rackName := testutil.RandomName("tf-test-rack")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttrPair("netbox_rack.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})

}

func TestAccRackResource_full(t *testing.T) {

	t.Parallel()

	// Generate unique names

	siteName := testutil.RandomName("tf-test-rack-site-full")

	siteSlug := testutil.RandomSlug("tf-test-rack-s-full")

	rackName := testutil.RandomName("tf-test-rack-full")

	description := testutil.RandomName("description")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackResourceConfig_full(siteName, siteSlug, rackName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttr("netbox_rack.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_rack.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rack.test", "u_height", "42"),

					resource.TestCheckResourceAttr("netbox_rack.test", "width", "19"),
				),
			},
		},
	})

}

func TestAccRackResource_update(t *testing.T) {

	t.Parallel()

	// Generate unique names

	siteName := testutil.RandomName("tf-test-rack-site-upd")

	siteSlug := testutil.RandomSlug("tf-test-rack-s-upd")

	rackName := testutil.RandomName("tf-test-rack-upd")

	updatedName := testutil.RandomName("tf-test-rack-upd-name")

	// Register cleanup (use original name for initial cleanup, register updated name too)

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterRackCleanup(updatedName)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},

			{

				Config: testAccRackResourceConfig_basic(siteName, siteSlug, updatedName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", updatedName),
				),
			},
		},
	})

}

func TestAccRackResource_withLocation(t *testing.T) {

	t.Parallel()

	// Generate unique names

	siteName := testutil.RandomName("tf-test-rack-site-loc")

	siteSlug := testutil.RandomSlug("tf-test-rack-s-loc")

	locationName := testutil.RandomName("tf-test-rack-location")

	locationSlug := testutil.RandomSlug("tf-test-rack-loc")

	rackName := testutil.RandomName("tf-test-rack-with-loc")

	// Register cleanup (rack first, then location, then site due to dependency)

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterLocationCleanup(locationSlug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttrPair("netbox_rack.test", "location", "netbox_location.test", "id"),
				),
			},
		},
	})

}

func TestAccRackResource_import(t *testing.T) {

	t.Parallel()

	// Generate unique names to avoid conflicts between test runs

	siteName := testutil.RandomName("tf-test-rack-site")

	siteSlug := testutil.RandomSlug("tf-test-rack-site")

	rackName := testutil.RandomName("tf-test-rack")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackResourceConfig_import(siteName, siteSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttrPair("netbox_rack.test", "site", "netbox_site.test", "id"),
				),
			},

			{

				ResourceName: "netbox_rack.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"site"},
			},

			{
				Config:   testAccRackResourceConfig_import(siteName, siteSlug, rackName),
				PlanOnly: true,
			},
		},
	})

}

func TestAccRackResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	siteName := testutil.RandomName("tf-test-rack-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-site")
	rackName := testutil.RandomName("tf-test-rack")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate test data for all custom field types (once, used in all steps)
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
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceImportConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				ResourceName:            "netbox_rack.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site", "tenant", "custom_fields", "tags"},
			},
			{
				Config: testAccRackResourceImportConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue),
				PlanOnly: true,
			},
		},
	})
}

func testAccRackResourceImportConfig_full(
	siteName, siteSlug, rackName, tenantName, tenantSlug string,
	tag1, tag1Slug, tag2, tag2Slug string,
	cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON string,
	textValue, longtextValue string, intValue int, boolValue bool, dateValue, urlValue, jsonValue string,
) string {

	return fmt.Sprintf(`
# Dependencies
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

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

# Custom Fields for dcim.rack object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["dcim.rack"]
}

# Rack with comprehensive custom fields and tags
resource "netbox_rack" "test" {
  name   = %q
  site   = netbox_site.test.id
  status = "active"
  tenant = netbox_tenant.test.id

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
`, siteName, siteSlug, tenantName, tenantSlug,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		rackName, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

func TestAccConsistency_Rack(t *testing.T) {

	t.Parallel()

	rackName := testutil.RandomName("rack")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttr("netbox_rack.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_rack.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_rack.test", "role", roleName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})

}

func TestAccConsistency_Rack_LiteralNames(t *testing.T) {
	t.Parallel()
	rackName := testutil.RandomName("tf-test-rack-lit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "site", siteName),
				),
			},
			{
				Config:   testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
				),
			},
		},
	})
}

func TestAccRackResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-id")
	siteSlug := testutil.RandomSlug("tf-test-site-id")
	rackName := testutil.RandomName("tf-test-rack-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "site"),
				),
			},
		},
	})
}

func testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.name
}
`, siteName, siteSlug, rackName)
}

// testAccRackResourceConfig_basic returns a basic test configuration.

func testAccRackResourceConfig_basic(siteName, siteSlug, rackName string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}

resource "netbox_rack" "test" {

  name = %q

  site = netbox_site.test.id

}

`, siteName, siteSlug, rackName)

}

// testAccRackResourceConfig_full returns a test configuration with all fields.

func testAccRackResourceConfig_full(siteName, siteSlug, rackName, description string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}

resource "netbox_rack" "test" {

  name        = %q

  site        = netbox_site.test.id

  status      = "active"

  u_height    = 42

  width       = 19

  description = %q

}

`, siteName, siteSlug, rackName, description)

}

// testAccRackResourceConfig_withLocation returns a test configuration with location.

func testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %q

  slug   = %q

  status = "active"

}

resource "netbox_location" "test" {

  name = %q

  slug = %q

  site = netbox_site.test.id

}

resource "netbox_rack" "test" {

  name     = %q

  site     = netbox_site.test.id

  location = netbox_location.test.id

}

`, siteName, siteSlug, locationName, locationSlug, rackName)

}

func testAccRackResourceConfig_import(siteName, siteSlug, rackName string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %[1]q

  slug = %[2]q

}

resource "netbox_rack" "test" {

  name = %[3]q

  site = netbox_site.test.id

}

`, siteName, siteSlug, rackName)

}

func testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}

resource "netbox_rack_role" "test" {

  name = "%[6]s"

  slug = "%[7]s"

}

resource "netbox_rack" "test" {

  name = "%[1]s"

  site = netbox_site.test.name

  tenant = netbox_tenant.test.name

  role = netbox_rack_role.test.name

}

`, rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug)

}

func TestAccRackResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-rack-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-rack-site-ed")
	rackName := testutil.RandomName("tf-test-rack-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					racks, _, err := client.DcimAPI.DcimRacksList(context.Background()).Name([]string{rackName}).Execute()
					if err != nil || racks == nil || len(racks.Results) == 0 {
						t.Fatalf("Failed to find rack for external deletion: %v", err)
					}
					rackID := racks.Results[0].Id
					_, err = client.DcimAPI.DcimRacksDestroy(context.Background(), rackID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack: %v", err)
					}
					t.Logf("Successfully externally deleted rack with ID: %d", rackID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
