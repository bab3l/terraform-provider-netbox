package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	isPoolFalse = "false"
	isPoolTrue  = "true"
)

func TestAccPrefixResource_basic(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
		},
	})
}

func TestAccPrefixResource_IDPreservation(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
				),
			},
		},
	})
}

func TestAccPrefixResource_update(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_forUpdate(prefix, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "false"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccPrefixResourceConfig_forUpdate(prefix, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "reserved"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "true"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccPrefixResource_full(t *testing.T) {

	t.Parallel()
	prefix := testutil.RandomIPv4Prefix()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	vrfName := testutil.RandomName("vrf")
	vlanName := testutil.RandomName("vlan")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("updated-description")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_full(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, description, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", description),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "true"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "mark_utilized", "true"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vlan"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "role"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "tags.#", "2"),
				),
			},
			{
				Config: testAccPrefixResourceConfig_fullUpdate(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "false"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "mark_utilized", "false"),
				),
			},
		},
	})
}

func TestAccPrefixResource_withVRF(t *testing.T) {

	t.Parallel()
	prefix := testutil.RandomIPv4Prefix()
	vrfName := testutil.RandomName("tf-test-vrf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPrefixDestroy,
			testutil.CheckVRFDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_withVRF(prefix, vrfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vrf"),
				),
			},
		},
	})
}

func TestAccPrefixResource_ipv6(t *testing.T) {

	t.Parallel()
	prefix := testutil.RandomIPv6Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
		},
	})
}

func TestAccPrefixResource_externalDeletion(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamPrefixesList(context.Background()).Prefix([]string{prefix}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find prefix for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamPrefixesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete prefix: %v", err)
					}
					t.Logf("Successfully externally deleted prefix with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPrefixResourceConfig_basic(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %q
}
`, prefix)
}

func testAccPrefixResourceConfig_forUpdate(prefix, description string) string {
	// Toggle status and is_pool based on description
	status := testutil.PrefixStatusActive
	isPool := isPoolFalse

	if description == testutil.Description2 {
		status = testutil.PrefixStatusReserved
		isPool = isPoolTrue
	}

	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix      = %[1]q
  status      = %[2]q
  is_pool     = %[3]s
  description = %[4]q
}
`, prefix, status, isPool, description)
}

func testAccPrefixResourceConfig_full(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, description, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[2]q
  slug   = %[3]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_vrf" "test" {
  name = %[6]q
}

resource "netbox_vlan" "test" {
  name = %[7]q
  vid  = 100
  site = netbox_site.test.id
}

resource "netbox_role" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_tag" "tag1" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_tag" "tag2" {
  name = %[13]q
  slug = %[14]q
}

resource "netbox_prefix" "test" {
  prefix        = %[1]q
  site          = netbox_site.test.id
  tenant        = netbox_tenant.test.id
  vrf           = netbox_vrf.test.id
  vlan          = netbox_vlan.test.id
  role          = netbox_role.test.id
  description   = %[10]q
  status        = "active"
  is_pool       = true
  mark_utilized = true

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
`, prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, description, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccPrefixResourceConfig_fullUpdate(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, description, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[2]q
  slug   = %[3]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_vrf" "test" {
  name = %[6]q
}

resource "netbox_vlan" "test" {
  name = %[7]q
  vid  = 100
  site = netbox_site.test.id
}

resource "netbox_role" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_tag" "tag1" {
  name = %[11]q
  slug = %[12]q
}

resource "netbox_tag" "tag2" {
  name = %[13]q
  slug = %[14]q
}

resource "netbox_prefix" "test" {
  prefix        = %[1]q
  site          = netbox_site.test.id
  tenant        = netbox_tenant.test.id
  vrf           = netbox_vrf.test.id
  vlan          = netbox_vlan.test.id
  role          = netbox_role.test.id
  description   = %[10]q
  status        = "active"
  is_pool       = false
  mark_utilized = false

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
`, prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, roleName, roleSlug, description, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccPrefixResourceConfig_withVRF(prefix, vrfName string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %q
}

resource "netbox_prefix" "test" {
  prefix = %q
  vrf    = netbox_vrf.test.name
}
`, vrfName, prefix)
}

func TestAccPrefixResource_import(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
			{
				ResourceName:      "netbox_prefix.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_Prefix(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	vlanName := testutil.RandomName("vlan")
	vlanVid := 100

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_prefix.test", "tenant", tenantName),
					resource.TestCheckResourceAttr("netbox_prefix.test", "vlan", vlanName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
			},
		},
	})
}

func testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName string, vlanVid int) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_vlan" "test" {
  name = "%[6]s"
  vid  = %[7]d
  site = netbox_site.test.id
}

resource "netbox_prefix" "test" {
  prefix = "%[1]s"
  site = netbox_site.test.name
  tenant = netbox_tenant.test.name
  vlan = netbox_vlan.test.name
}
`, prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid)
}

func TestAccConsistency_Prefix_LiteralNames(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	vlanName := testutil.RandomName("vlan")
	vlanVid := 200

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixConsistencyLiteralNamesConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_prefix.test", "tenant", tenantName),
					resource.TestCheckResourceAttr("netbox_prefix.test", "vlan", vlanName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPrefixConsistencyLiteralNamesConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
			},
		},
	})
}

func testAccPrefixConsistencyLiteralNamesConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName string, vlanVid int) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_vlan" "test" {
  name = "%[6]s"
  vid  = %[7]d
  site = netbox_site.test.id
}

resource "netbox_prefix" "test" {
  prefix = "%[1]s"
  site = "%[2]s"
  tenant = "%[4]s"
  vlan = "%[6]s"
  depends_on = [netbox_site.test, netbox_tenant.test, netbox_vlan.test]
}
`, prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid)
}

func TestAccPrefixResource_importWithTags(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceImportConfig_full(prefix, tenantName, tenantSlug, tag1Name, tag1Slug, tag2Name, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_prefix.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_prefix.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant"}, // Tenant may have lookup inconsistencies
			},
		},
	})
}

func testAccPrefixResourceImportConfig_full(prefix, tenantName, tenantSlug, tag1Name, tag1Slug, tag2Name, tag2Slug string) string {
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

# Prefix with tags (no custom fields support)
resource "netbox_prefix" "test" {
  prefix = %q
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
}
`, tenantName, tenantSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, prefix)
}

// TestAccPrefixResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccPrefixResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	siteName := testutil.RandomName("test-site-prefix")
	siteSlug := testutil.RandomSlug("test-site-prefix")
	tenantName := testutil.RandomName("test-tenant-prefix")
	tenantSlug := testutil.RandomSlug("test-tenant-prefix")
	vrfName := testutil.RandomName("tf-test-vrf-prefix")
	vlanName := testutil.RandomName("tf-test-vlan-prefix")
	vlanVID := testutil.RandomVID()
	roleName := testutil.RandomName("test-role-prefix")
	roleSlug := testutil.RandomSlug("test-role-prefix")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterVRFCleanup(vrfName)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create prefix with all nullable fields
			{
				Config: testAccPrefixResourceConfig_withAllFields(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, vlanVID, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vlan"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "role"),
				),
			},
			// Step 2: Remove all nullable fields - should set them to null, and verify defaults
			{
				Config: testAccPrefixResourceConfig_withoutFields(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, vlanVID, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckNoResourceAttr("netbox_prefix.test", "site"),
					resource.TestCheckNoResourceAttr("netbox_prefix.test", "vrf"),
					resource.TestCheckNoResourceAttr("netbox_prefix.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_prefix.test", "vlan"),
					resource.TestCheckNoResourceAttr("netbox_prefix.test", "role"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "false"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "mark_utilized", "false"),
				),
			},
			// Step 3: Re-add all fields - verify they can be set again
			{
				Config: testAccPrefixResourceConfig_withAllFields(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, vlanVID, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vlan"),
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "role"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "reserved"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "true"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "mark_utilized", "true"),
				),
			},
		},
	})
}

func testAccPrefixResourceConfig_withAllFields(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName string, vlanVID int32, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_vrf" "test" {
  name = %q
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
  site = netbox_site.test.id
}

resource "netbox_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_prefix" "test" {
  prefix        = %q
  site          = netbox_site.test.id
  vrf           = netbox_vrf.test.id
  tenant        = netbox_tenant.test.id
  vlan          = netbox_vlan.test.id
  role          = netbox_role.test.id
  status        = "reserved"
  is_pool       = true
  mark_utilized = true
}
`, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, vlanVID, roleName, roleSlug, prefix)
}

func testAccPrefixResourceConfig_withoutFields(prefix, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName string, vlanVID int32, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_vrf" "test" {
  name = %q
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
  site = netbox_site.test.id
}

resource "netbox_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_prefix" "test" {
  prefix = %q
}
`, siteName, siteSlug, tenantName, tenantSlug, vrfName, vlanName, vlanVID, roleName, roleSlug, prefix)
}

func TestAccPrefixResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_prefix",
		BaseConfig: func() string {
			return testAccPrefixResourceConfig_basic(prefix)
		},
		ConfigWithFields: func() string {
			return testAccPrefixResourceConfig_withDescriptionAndComments(
				prefix,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"prefix": prefix,
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
	})
}

func testAccPrefixResourceConfig_withDescriptionAndComments(prefix, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix      = %[1]q
  status      = "active"
  description = %[2]q
  comments    = %[3]q
}
`, prefix, description, comments)
}

func TestAccPrefixResource_validationErrors(t *testing.T) {
	t.Parallel()

	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_prefix",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_prefix": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  status = "active"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_cidr_format": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "not-a-valid-cidr"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidFormat,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/8"
  status = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_site_reference": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/8"
  site   = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_vrf_reference": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/8"
  vrf    = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/8"
  tenant = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_vlan_reference": {
				Config: func() string {
					return `
resource "netbox_prefix" "test" {
  prefix = "10.0.0.0/8"
  vlan   = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccPrefixResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccPrefixResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_prefix",
		ConfigWithoutTags: func() string {
			return testAccPrefixResourceConfig_tagLifecycle(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccPrefixResourceConfig_tagLifecycle(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccPrefixResourceConfig_tagLifecycle(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckPrefixDestroy,
	})
}

func testAccPrefixResourceConfig_tagLifecycle(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = %[5]q
  slug = %[6]q
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %[1]q
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
`, prefix)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %[1]q
  tags = [
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag3.name
      slug = netbox_tag.tag3.slug
    }
  ]
}
`, prefix)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %[1]q
  tags   = []
}
`, prefix)
	}
}

// TestAccPrefixResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccPrefixResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_prefix",
		ConfigWithTagsOrderA: func() string {
			return testAccPrefixResourceConfig_tagOrder(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccPrefixResourceConfig_tagOrder(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckPrefixDestroy,
	})
}

func testAccPrefixResourceConfig_tagOrder(prefix, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = %[3]q
  slug = %[4]q
}
`, tag1Name, tag1Slug, tag2Name, tag2Slug)

	if tag1First {
		return baseConfig + fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %[1]q
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
`, prefix)
	}

	return baseConfig + fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %[1]q
  tags = [
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    },
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    }
  ]
}
`, prefix)
}
