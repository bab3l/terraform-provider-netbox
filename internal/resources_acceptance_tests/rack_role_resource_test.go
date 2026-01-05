package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackRoleResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	rackRoleName := testutil.RandomName("tf-test-rack-role")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
				),
			},
		},
	})
}

func TestAccRackRoleResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-full")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-f")
	description := testutil.RandomName("description")
	color := testutil.ColorOrange

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", color),
				),
			},
		},
	})
}

func TestAccRackRoleResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-upd")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-u")
	updatedDescription := testutil.Description2

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
				),
			},
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, updatedDescription, "00bcd4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", "00bcd4"),
				),
			},
		},
	})
}

func TestAccRackRoleResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names
	rackRoleName := testutil.RandomName("tf-test-rack-role-imp")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-i")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
			},
			{
				ResourceName:      "netbox_rack_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_RackRole(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("rack-role")
	rackRoleSlug := testutil.RandomSlug("rack-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleConsistencyConfig(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRackRoleConsistencyConfig(rackRoleName, rackRoleSlug),
			},
		},
	})
}

func TestAccRackRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-rack-role-id")
	slug := testutil.RandomSlug("tf-test-rr-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccRackRoleResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccRackRoleResourceConfig_full(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  color       = %[4]q
}
`, name, slug, description, color)
}

func TestAccConsistency_RackRole_LiteralNames(t *testing.T) {
	t.Parallel()

	rackRoleName := testutil.RandomName("tf-test-rack-role-lit")
	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-lit")
	description := testutil.RandomName("description")
	color := "4caf50"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", color),
				),
			},
			{
				Config:   testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
				),
			},
		},
	})
}

func testAccRackRoleConsistencyConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccRackRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	rackRoleName := testutil.RandomName("tf-test-rack-role-extdel")
	rackRoleSlug := testutil.RandomSlug("tf-test-rr-ed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					roles, _, err := client.DcimAPI.DcimRackRolesList(context.Background()).Slug([]string{rackRoleSlug}).Execute()
					if err != nil || roles == nil || len(roles.Results) == 0 {
						t.Fatalf("Failed to find rack role for external deletion: %v", err)
					}
					roleID := roles.Results[0].Id
					_, err = client.DcimAPI.DcimRackRolesDestroy(context.Background(), roleID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack role: %v", err)
					}
					t.Logf("Successfully externally deleted rack role with ID: %d", roleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
