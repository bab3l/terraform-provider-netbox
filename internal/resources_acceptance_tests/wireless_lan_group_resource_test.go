package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-wlan-group")
	slug := testutil.RandomSlug("tf-test-wlan-group")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_wireless_lan_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWirelessLANGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-wlan-group-full")
	slug := testutil.RandomSlug("tf-test-wlan-group-full")
	description := "Test wireless LAN group with all fields"
	updatedDescription := "Updated wireless LAN group description"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", description),
				),
			},
			{
				Config: testAccWirelessLANGroupResourceConfig_full(name, slug, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccWirelessLANGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-wlan-group-id")
	slug := testutil.RandomSlug("tf-test-wlan-group-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),
				),
			},
		},
	})
}

func testAccWirelessLANGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccWirelessLANGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccWirelessLANGroupResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-wlan-group-upd")
	slug := testutil.RandomSlug("tf-test-wlan-group-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_withDescription(name, slug, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccWirelessLANGroupResourceConfig_withDescription(name, slug, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func testAccWirelessLANGroupResourceConfig_withDescription(name string, slug string, description string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func TestAccWirelessLANGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-wlan-group-extdel")
	slug := testutil.RandomSlug("tf-test-wlan-group-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find wireless LAN group by slug
					items, _, err := client.WirelessAPI.WirelessWirelessLanGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil {
						t.Fatalf("Failed to list wireless LAN groups: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Wireless LAN group not found with slug: %s", slug)
					}

					// Delete the wireless LAN group
					itemID := items.Results[0].Id
					_, err = client.WirelessAPI.WirelessWirelessLanGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete wireless LAN group: %v", err)
					}

					t.Logf("Successfully externally deleted wireless LAN group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWirelessLANGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	parentName := testutil.RandomName("tf-test-wlan-group-parent")
	parentSlug := testutil.RandomSlug("tf-test-wlan-group-parent")
	name := testutil.RandomName("tf-test-wlan-group-rem")
	slug := testutil.RandomSlug("tf-test-wlan-group-rem")
	const testDescription = "Test Description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterWirelessLANGroupCleanup(parentName)
	cleanup.RegisterWirelessLANGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckWirelessLANGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupResourceConfig_removeOptionalFields_withParent(parentName, parentSlug, name, slug, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", testDescription),
					resource.TestCheckResourceAttrPair("netbox_wireless_lan_group.test", "parent", "netbox_wireless_lan_group.parent", "id"),
				),
			},
			{
				Config: testAccWirelessLANGroupResourceConfig_removeOptionalFields_noOptional(parentName, parentSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_wireless_lan_group.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_wireless_lan_group.test", "parent"),
				),
			},
		},
	})
}

func testAccWirelessLANGroupResourceConfig_removeOptionalFields_withParent(parentName, parentSlug, name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_wireless_lan_group" "test" {
  name        = %q
  slug        = %q
  description = %q
  parent      = netbox_wireless_lan_group.parent.id
}
`, parentName, parentSlug, name, slug, description)
}

func testAccWirelessLANGroupResourceConfig_removeOptionalFields_noOptional(parentName, parentSlug, name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "parent" {
  name = %q
  slug = %q
}

resource "netbox_wireless_lan_group" "test" {
  name = %q
  slug = %q
}
`, parentName, parentSlug, name, slug)
}
