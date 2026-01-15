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
				ResourceName:      "netbox_rack_reservation.test",
				ImportState:       true,
				ImportStateVerify: true,
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
	cfName := testutil.RandomCustomFieldName("test_field")

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
				Config: testAccRackReservationResourceConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "3"),
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config: testAccRackReservationResourceConfig_fullUpdate(siteName, siteSlug, rackName, tenantName, tenantSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "2"),
					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "custom_fields.0.value", "updated_value"),
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

func testAccRackReservationResourceConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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

resource "netbox_custom_field" "test_field" {
  name         = %[12]q
  object_types = ["dcim.rackreservation"]
  type         = "text"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2, 3]
  user        = data.netbox_user.admin.id
  tenant      = netbox_tenant.test.id
  description = %[6]q
  comments    = %[7]q

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
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccRackReservationResourceConfig_fullUpdate(siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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

resource "netbox_custom_field" "test_field" {
  name         = %[12]q
  object_types = ["dcim.rackreservation"]
  type         = "text"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [5, 6]
  user        = data.netbox_user.admin.id
  tenant      = netbox_tenant.test.id
  description = %[6]q
  comments    = %[7]q

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
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}
`, siteName, siteSlug, rackName, tenantName, tenantSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
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

func TestAccRackReservationResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-rr")
	siteSlug := testutil.GenerateSlug(siteName)
	rackName := testutil.RandomName("rack-rr")
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
  rack        = netbox_rack.test.name
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
