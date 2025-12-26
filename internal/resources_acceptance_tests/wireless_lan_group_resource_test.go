package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANGroupResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-wlan-group")
	slug := testutil.RandomSlug("tf-test-wlan-group")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
