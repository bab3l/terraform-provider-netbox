package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeResource_basic(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)
	description := testutil.RandomName("ip-range-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(101, 150)
	startOctet2 := 10 + acctest.RandIntRange(1, 200)
	endOctet2 := startOctet2 + 10
	startAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet2)
	endAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet2)
	description := testutil.RandomName("ip-range-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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

	secondOctet := acctest.RandIntRange(151, 200)
	thirdOctet := acctest.RandIntRange(151, 200)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d/32", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ip_range.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				PlanOnly: true,
			},
		},
	})
}

func TestAccIPRangeResource_external_deletion(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(201, 250)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamIpRangesList(context.Background()).StartAddress([]string{startAddress}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IP range for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamIpRangesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IP range: %v", err)
					}
					t.Logf("Successfully externally deleted IP range with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIPRangeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(101, 150)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
			{
				Config:   testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
		},
	})
}
