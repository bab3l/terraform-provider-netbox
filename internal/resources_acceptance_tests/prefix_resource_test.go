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
	description := testutil.RandomName("description")

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
				Config: testAccPrefixResourceConfig_full(prefix, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_prefix.test", "description", description),
					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "false"),
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
				Config: testAccPrefixResourceConfig_full(prefix, "Updated description"),
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
				Config: testAccPrefixResourceConfig_basic(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),
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

func testAccPrefixResourceConfig_full(prefix, description string) string {
	return fmt.Sprintf(`
resource "netbox_prefix" "test" {
  prefix      = %q
  description = %q
  status      = "active"
  is_pool     = false
}
`, prefix, description)
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
