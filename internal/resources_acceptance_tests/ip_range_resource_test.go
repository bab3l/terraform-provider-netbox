package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeResource_basic(t *testing.T) {

	t.Parallel()

	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("192.0.2.%d", startOctet)
	endAddress := fmt.Sprintf("192.0.2.%d", endOctet)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
		},
	})

}

func TestAccIPRangeResource_full(t *testing.T) {

	t.Parallel()

	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("192.0.3.%d", startOctet)
	endAddress := fmt.Sprintf("192.0.3.%d", endOctet)
	description := testutil.RandomName("ip-range-desc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_full(startAddress, endAddress, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})

}

func TestAccIPRangeResource_update(t *testing.T) {

	t.Parallel()

	startOctet2 := 10 + acctest.RandIntRange(1, 200)
	endOctet2 := startOctet2 + 10
	startAddress2 := fmt.Sprintf("192.0.4.%d", startOctet2)
	endAddress2 := fmt.Sprintf("192.0.4.%d", endOctet2)
	description := testutil.RandomName("ip-range-desc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress2, endAddress2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress2),
				),
			},

			{

				Config: testAccIPRangeResourceConfig_full(startAddress2, endAddress2, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})

}

func TestAccIPRangeResource_import(t *testing.T) {

	t.Parallel()

	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("192.0.2.%d/32", startOctet)
	endAddress := fmt.Sprintf("192.0.2.%d/32", endOctet)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ip_range.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIPRangeResourceConfig_basic(startAddress, endAddress string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %[1]q

  end_address   = %[2]q

}

`, startAddress, endAddress)

}

func testAccIPRangeResourceConfig_full(startAddress, endAddress, description string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %[1]q

  end_address   = %[2]q

  status        = "active"

  description   = %[3]q

}

`, startAddress, endAddress, description)

}

func TestAccConsistency_IPRange_LiteralNames(t *testing.T) {
	t.Parallel()
	startOctet := 50 + acctest.RandIntRange(1, 100)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("172.16.%d.10", startOctet)
	endAddress := fmt.Sprintf("172.16.%d.20", endOctet)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
			{
				Config:   testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
		},
	})
}

func testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %q
  end_address   = %q
}
`, startAddress, endAddress)
}
