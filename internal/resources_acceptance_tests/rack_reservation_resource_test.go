package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackReservationResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	rackName := testutil.RandomName("tf-test-rack")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "2"),
				),
			},
			// ImportState test
			{
				ResourceName: "netbox_rack_reservation.test",
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_rack_reservation.test", "rack"),
				),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:             testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccRackReservationResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	rackName := testutil.RandomName("tf-test-rack")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("updated-description")
	comments := "Initial reservation comments"
	updatedComments := "Updated reservation comments"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "3"),
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_fullUpdate(siteName, siteSlug, rackName, tenantName, tenantSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "2"),
				),
			},
		},
	})

}

func TestAccRackReservationResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	rackName := testutil.RandomName("tf-test-rack")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccRackReservationResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-tags")
	rackName := testutil.RandomName("tf-test-rack-tags")
	description := testutil.RandomName("description")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_tags(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_tags(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_tags(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_tags(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRackReservationResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-site-tag-order")
	rackName := testutil.RandomName("tf-test-rack-tag-order")
	description := testutil.RandomName("description")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_tagsOrder(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_tagsOrder(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_rack_reservation.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = %[4]q
}
`, siteName, siteSlug, rackName, description)
}

func testAccRackReservationResourceConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_tag" "tag1" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_tag" "tag2" {
  name = %[10]q
  slug = %[11]q
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2, 3]
  user        = data.netbox_user.admin.id
  tenant      = netbox_tenant.test.id
  description = %[6]q
  comments    = %[7]q

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccRackReservationResourceConfig_fullUpdate(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_tag" "tag1" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_tag" "tag2" {
  name = %[10]q
  slug = %[11]q
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [5, 6]
  user        = data.netbox_user.admin.id
  tenant      = netbox_tenant.test.id
  description = %[6]q
  comments    = %[7]q

  tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccRackReservationResourceConfig_tags(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_rack" "test" {
	name     = %[3]q
	site     = netbox_site.test.id
	status   = "active"
	u_height = 42
}

data "netbox_user" "admin" {
	username = "admin"
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[7]s"
	slug = %[7]q
}

resource "netbox_rack_reservation" "test" {
	rack        = netbox_rack.test.id
	units       = [1, 2]
	user        = data.netbox_user.admin.id
	description = %[4]q
	%[8]s
}
`, siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRackReservationResourceConfig_tagsOrder(siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_rack" "test" {
	name     = %[3]q
	site     = netbox_site.test.id
	status   = "active"
	u_height = 42
}

data "netbox_user" "admin" {
	username = "admin"
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_rack_reservation" "test" {
	rack        = netbox_rack.test.id
	units       = [1, 2]
	user        = data.netbox_user.admin.id
	description = %[4]q
	%[7]s
}
`, siteName, siteSlug, rackName, description, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_RackReservation_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	rackName := testutil.RandomName("rack")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description),
			},
		},
	})
}

func testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
	rack        = netbox_rack.test.id
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = %[4]q
}
`, siteName, siteSlug, rackName, description)
}

func TestAccRackReservationResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-site-ed")
	rackName := testutil.RandomName("tf-test-rack-extdel")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					results, _, err := client.DcimAPI.DcimRackReservationsList(context.Background()).Description([]string{description}).Execute()
					if err != nil || results == nil || len(results.Results) == 0 {
						t.Fatalf("Failed to find rack reservation for external deletion: %v", err)
					}
					reservationID := results.Results[0].Id
					_, err = client.DcimAPI.DcimRackReservationsDestroy(context.Background(), reservationID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack reservation: %v", err)
					}
					t.Logf("Successfully externally deleted rack reservation with ID: %d", reservationID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccRackReservationResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccRackReservationResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-rem")
	rackName := testutil.RandomName("tf-test-rack-rem")
	tenantName := testutil.RandomName("tf-test-tenant-rem")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-rem")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationResourceConfig_withTenant(siteName, siteSlug, rackName, tenantName, tenantSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "tenant"),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_withoutTenantRef(siteName, siteSlug, rackName, tenantName, tenantSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
					resource.TestCheckNoResourceAttr("netbox_rack_reservation.test", "tenant"),
				),
			},
		},
	})
}

func testAccRackReservationResourceConfig_withTenant(siteName, siteSlug, rackName, tenantName, tenantSlug, description string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = %[6]q
  tenant      = netbox_tenant.test.id
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description)
}
func testAccRackReservationResourceConfig_withoutTenantRef(siteName, siteSlug, rackName, tenantName, tenantSlug, description string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = %[6]q
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description)
}

func TestAccRackReservationResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_rack_reservation",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_rack": {
				Config: func() string {
					return `
data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  units = [1, 2]
  user  = data.netbox_user.admin.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_units": {
				Config: func() string {
					return `
data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack = "1"
  user = data.netbox_user.admin.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_user": {
				Config: func() string {
					return `
resource "netbox_rack_reservation" "test" {
  rack  = "1"
  units = [1, 2]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_rack_reference": {
				Config: func() string {
					return `
data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack  = "99999"
  units = [1, 2]
  user  = data.netbox_user.admin.id
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
