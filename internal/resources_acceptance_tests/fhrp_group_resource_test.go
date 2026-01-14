package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const fhrpGroupProtocol = "vrrp3"

func TestAccFHRPGroupResource_basic(t *testing.T) {
	t.Parallel()

	protocol := "vrrp2"
	// Use non-overlapping range to prevent parallel test collisions
	groupID := int32(acctest.RandIntRange(106, 140)) // nolint:gosec

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				Config:   testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				PlanOnly: true,
			},
		},
	})
}

func TestAccFHRPGroupResource_full(t *testing.T) {
	t.Parallel()

	protocol := "hsrp"
	// Use non-overlapping range to prevent parallel test collisions
	groupID := int32(acctest.RandIntRange(36, 70)) // nolint:gosec
	name := testutil.RandomName("tf-test-fhrp")
	description := testutil.RandomName("description")
	authType := "plaintext"
	authKey := "secretkey123"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, name, description, authType, authKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", authType),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", authKey),
				),
			},
			{
				Config:   testAccFHRPGroupResourceConfig_full(protocol, groupID, name, description, authType, authKey),
				PlanOnly: true,
			},
		},
	})
}

func TestAccFHRPGroupResource_update(t *testing.T) {
	t.Parallel()

	protocol := fhrpGroupProtocol
	// Use non-overlapping range to prevent parallel test collisions
	groupID := int32(acctest.RandIntRange(71, 105)) // nolint:gosec
	updatedName := testutil.RandomName("tf-test-fhrp-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
				),
			},
			{
				Config:   testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				PlanOnly: true,
			},
			{
				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, updatedName, "Updated description", "md5", "newsecret456"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", "md5"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", "newsecret456"),
				),
			},
			{
				Config:   testAccFHRPGroupResourceConfig_full(protocol, groupID, updatedName, "Updated description", "md5", "newsecret456"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccFHRPGroupResource_external_deletion(t *testing.T) {
	t.Parallel()

	protocol := fhrpGroupProtocol
	// Use non-overlapping range to prevent parallel test collisions
	groupID := int32(acctest.RandIntRange(106, 140)) // nolint:gosec

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamFhrpGroupsList(context.Background()).Protocol([]string{protocol}).GroupId([]int32{groupID}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find fhrp_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamFhrpGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete fhrp_group: %v", err)
					}
					t.Logf("Successfully externally deleted fhrp_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFHRPGroupResource_import(t *testing.T) {
	t.Parallel()

	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) // nolint:gosec

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				ResourceName:      "netbox_fhrp_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				PlanOnly: true,
			},
		},
	})
}

func testAccFHRPGroupResourceConfig_basic(protocol string, groupID int32) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol = %q
  group_id = %d
}
`, protocol, groupID)
}

func testAccFHRPGroupResourceConfig_full(protocol string, groupID int32, name, description, authType, authKey string) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol    = %q
  group_id    = %d
  name        = %q
  description = %q
  auth_type   = %q
  auth_key    = %q
}
`, protocol, groupID, name, description, authType, authKey)
}

func TestAccConsistency_FHRPGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	protocol := fhrpGroupProtocol
	groupID := int32(123)
	name := testutil.RandomName("tf-test-fhrp-group-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckFHRPGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupConsistencyLiteralNamesConfig(protocol, groupID, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", "123"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", description),
				),
			},
			{
				Config:   testAccFHRPGroupConsistencyLiteralNamesConfig(protocol, groupID, name, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
				),
			},
		},
	})
}

// TestAccFHRPGroupResource_IDPreservation tests that the FHRP group resource preserves the
// ID as the immutable identifier.
func TestAccFHRPGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	protocol := fhrpGroupProtocol
	// Use non-overlapping range to prevent parallel test collisions
	groupID := int32(acctest.RandIntRange(211, 254)) // nolint:gosec

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				Config:   testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				PlanOnly: true,
			},
		},
	})
}

func testAccFHRPGroupConsistencyLiteralNamesConfig(protocol string, groupID int32, name, description string) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol    = %q
  group_id    = %d
  name        = %q
  description = %q
}
`, protocol, groupID, name, description)
}

func TestAccFHRPGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	protocol := "vrrp3"
	groupID := int32(acctest.RandIntRange(141, 175)) // nolint:gosec
	name := testutil.RandomName("tf-test-fhrp-opt")
	authType := "md5"
	authKey := "secret123"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol   = %[1]q
  group_id   = %[2]d
  name       = %[3]q
  auth_type  = %[4]q
  auth_key   = %[5]q
  description = "Description"
  comments    = "Comments"
}
`, protocol, groupID, name, authType, authKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", authType),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", authKey),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", "Description"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "comments", "Comments"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol = %[1]q
  group_id = %[2]d
}
`, protocol, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckNoResourceAttr("netbox_fhrp_group.test", "name"),
					resource.TestCheckNoResourceAttr("netbox_fhrp_group.test", "auth_type"),
					resource.TestCheckNoResourceAttr("netbox_fhrp_group.test", "auth_key"),
					resource.TestCheckNoResourceAttr("netbox_fhrp_group.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_fhrp_group.test", "comments"),
				),
			},
		},
	})
}
