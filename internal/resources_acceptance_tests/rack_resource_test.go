package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
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
				ResourceName:            "netbox_rack.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site"},
			},
			{
				Config:   testAccRackResourceConfig_import(siteName, siteSlug, rackName),
				PlanOnly: true,
			},
		},
	})
}

// NOTE: Custom field tests for rack resource are in resources_acceptance_tests_customfields package.
func TestAccConsistency_Rack(t *testing.T) {
	t.Parallel()

	rackName := testutil.RandomName("rack")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRackRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
				Config:   testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})
}

func TestAccConsistency_Rack_LiteralNames(t *testing.T) {
	t.Parallel()

	rackName := testutil.RandomName("tf-test-rack-lit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

// TestAccRackResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccRackResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rack")
	siteSlug := testutil.RandomSlug("tf-test-site-rack")
	locationName := testutil.RandomName("tf-test-loc-rack")
	locationSlug := testutil.RandomSlug("tf-test-loc-rack")
	tenantName := testutil.RandomName("tf-test-tenant-rack")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-rack")
	roleName := testutil.RandomName("tf-test-role-rack")
	roleSlug := testutil.RandomSlug("tf-test-role-rack")
	rackTypeName := testutil.RandomName("tf-test-racktype")
	rackTypeSlug := testutil.RandomSlug("tf-test-racktype")
	rackName := testutil.RandomName("tf-test-rack")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterLocationCleanup(locationSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRackRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRackDestroy,
			testutil.CheckLocationDestroy,
			testutil.CheckTenantDestroy,
			testutil.CheckRackRoleDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create rack with location, tenant, role, and rack_type
			{
				Config: testAccRackResourceConfig_withAllFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "location"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "rack_type"),
				),
			},
			// Step 2: Remove location, tenant, role, and rack_type - should set them to null
			{
				Config: testAccRackResourceConfig_withoutOptionalFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckNoResourceAttr("netbox_rack.test", "location"),
					resource.TestCheckNoResourceAttr("netbox_rack.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_rack.test", "role"),
					resource.TestCheckNoResourceAttr("netbox_rack.test", "rack_type"),
				),
			},
			// Step 3: Re-add location, tenant, role, and rack_type - verify they can be set again
			{
				Config: testAccRackResourceConfig_withAllFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "location"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_rack.test", "rack_type"),
				),
			},
		},
	})
}

func testAccRackResourceConfig_withAllFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_rack_type" "test" {
  model         = %q
  slug          = %q
  manufacturer  = netbox_manufacturer.test.id
  form_factor   = "4-post-cabinet"
}

resource "netbox_rack" "test" {
  name      = %q
  site      = netbox_site.test.id
  location  = netbox_location.test.id
  tenant    = netbox_tenant.test.id
  role      = netbox_rack_role.test.id
  rack_type = netbox_rack_type.test.id
}
`, siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName)
}

func testAccRackResourceConfig_withoutOptionalFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_rack_type" "test" {
  model         = %q
  slug          = %q
  manufacturer  = netbox_manufacturer.test.id
  form_factor   = "4-post-cabinet"
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.id
}
`, siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackTypeName, rackTypeSlug, rackName)
}
