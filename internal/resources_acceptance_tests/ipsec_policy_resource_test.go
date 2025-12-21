package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECPolicyResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
		},
	})

}

func TestAccIPSECPolicyResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECPolicyResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "group14"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Test IPsec policy"),
				),
			},
		},
	})

}

func TestAccIPSECPolicyResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},

			{

				Config: testAccIPSECPolicyResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "group14"),

					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Test IPsec policy"),
				),
			},
		},
	})

}

func TestAccIPSECPolicyResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ipsec_policy.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIPSECPolicyResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_policy" "test" {

  name = %q

}

`, name)

}

func testAccIPSECPolicyResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_policy" "test" {

  name        = %q

  pfs_group   = "group14"

  description = "Test IPsec policy"

}

`, name)

}
