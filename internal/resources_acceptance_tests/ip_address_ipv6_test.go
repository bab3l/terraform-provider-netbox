package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressResource_ipv6_simple(t *testing.T) {

	t.Parallel()

	// IPv6 groups must be 1-4 hex digits (0-ffff). Use %x to format as hex.

	ip6 := fmt.Sprintf("2001:db8:%x::1/64", acctest.RandIntRange(1, 65535))

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(ip6)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: fmt.Sprintf(`

resource "netbox_ip_address" "test_ipv6" {

  address = %q

  status  = "active"

}

`, ip6),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test_ipv6", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test_ipv6", "address", ip6),
				),
			},
		},
	})

}
