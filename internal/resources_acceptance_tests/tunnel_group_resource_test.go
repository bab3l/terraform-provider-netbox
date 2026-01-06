package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelGroupResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tunnel-group")
	slug := testutil.RandomSlug("tf-test-tunnel-grp")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tg-id")
	slug := testutil.GenerateSlug(name)

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-full")
	slug := testutil.RandomSlug("tf-test-tg-full")
	description := "Test tunnel group with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-upd")
	slug := testutil.RandomSlug("tf-test-tg-upd")
	updatedDescription := testutil.Description2

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_full(name, slug, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-imp")
	slug := testutil.RandomSlug("tf-test-tg-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
			},
			{
				ResourceName:      "netbox_tunnel_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTunnelGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-group-extdel")
	slug := testutil.RandomSlug("tf-test-tg-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnTunnelGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tunnel group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnTunnelGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tunnel group: %v", err)
					}
					t.Logf("Successfully externally deleted tunnel group with ID: %d", itemID)
				},
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
				),
			},
		},
	})
}

func TestAccConsistency_TunnelGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tg")
	slug := testutil.RandomSlug("tg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug),
			},
		},
	})
}

func testAccTunnelGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccTunnelGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}
