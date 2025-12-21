package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomLinkResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomLinkDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomLinkResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_text", "View in External System"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_url", "https://example.com/device/{{ object.name }}"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "object_types.#", "1"),
				),
			},

			{

				ResourceName: "netbox_custom_link.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCustomLinkResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomLinkDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomLinkResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "enabled", "true"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "weight", "50"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "group_name", "External Links"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "button_class", "blue"),

					resource.TestCheckResourceAttr("netbox_custom_link.test", "new_window", "true"),
				),
			},
		},
	})

}

func TestAccCustomLinkResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cl")

	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomLinkDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomLinkResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
				),
			},

			{

				Config: testAccCustomLinkResourceConfig_basic(updatedName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", updatedName),
				),
			},
		},
	})

}

func testAccCustomLinkResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_link" "test" {

  name         = "%s"

  object_types = ["dcim.device"]

  link_text    = "View in External System"

  link_url     = "https://example.com/device/{{ object.name }}"

}

`, name)

}

func testAccCustomLinkResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_link" "test" {

  name         = "%s"

  object_types = ["dcim.device", "dcim.site"]

  enabled      = true

  link_text    = "View Details"

  link_url     = "https://example.com/{{ object.name }}"

  weight       = 50

  group_name   = "External Links"

  button_class = "blue"

  new_window   = true

}

`, name)

}
