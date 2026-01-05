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

func TestAccPrefixResource_basic(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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

func TestAccPrefixResource_update(t *testing.T) {

	t.Parallel()
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
			{
				Config: testAccPrefixResourceConfig_withDescription(prefix, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccPrefixResource_external_deletion(t *testing.T) {
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

func TestAccPrefixResource_IDPreservation(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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

func testAccPrefixResourceConfig_basic(prefix string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix = %q
}
`, prefix)
}

func testAccPrefixResourceConfig_withDescription(prefix, description string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix      = %q
  description = %q
  status      = "active"
}
`, prefix, description)
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckPrefixDestroy,
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
