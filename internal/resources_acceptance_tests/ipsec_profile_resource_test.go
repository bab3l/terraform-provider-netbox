package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECProfileResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProfileResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})

}

func TestAccIPSECProfileResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProfileResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "ah"),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})

}

func TestAccIPSECProfileResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProfileResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},

			{

				Config: testAccIPSECProfileResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})

}

func TestAccIPSECProfileResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProfileResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ipsec_profile.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIPSECProfileResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_ipsec_profile" "test" {

  name           = %q

  mode           = "esp"

  ike_policy     = netbox_ike_policy.test.id

  ipsec_policy   = netbox_ipsec_policy.test.id

}

`, testAccIPSECProfileResourcePrereqs(name), name)

}

func testAccIPSECProfileResourceConfig_full(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_ipsec_profile" "test" {

  name           = %q

  mode           = "ah"

  ike_policy     = netbox_ike_policy.test.id

  ipsec_policy   = netbox_ipsec_policy.test.id

  description    = "Test IPsec profile"

}

`, testAccIPSECProfileResourcePrereqs(name), name)

}

func testAccIPSECProfileResourcePrereqs(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_policy" "test" {

  name = %q

}

resource "netbox_ipsec_policy" "test" {

  name = %q

}

`, name+"-ike-policy", name+"-ipsec-policy")

}
