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

func TestAccVLANGroupResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vlangrp")
	slug := testutil.GenerateSlug("tf-test-vlangrp")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccVLANGroupResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vlangrp-full")
	slug := testutil.GenerateSlug("tf-test-vlangrp-full")
	description := "Test VLAN Group with all fields"
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccVLANGroupResource_update(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vlangrp-upd")
	slug := testutil.GenerateSlug("tf-test-vlangrp-upd")
	updatedName := testutil.RandomName("tf-test-vlangrp-updated")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},
			{
				Config: testAccVLANGroupResourceConfig_full(updatedName, slug, "Updated description"),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccVLANGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-vlang-extdel")
	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find VLAN group by slug
					items, _, err := client.IpamAPI.IpamVlanGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil {
						t.Fatalf("Failed to list VLAN groups: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("VLAN group not found with slug: %s", slug)
					}

					// Delete the VLAN group
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamVlanGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete VLAN group: %v", err)
					}

					t.Logf("Successfully externally deleted VLAN group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVLANGroupResource_import(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vlangrp")
	slug := testutil.GenerateSlug("tf-test-vlangrp")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_vlan_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_VLANGroup_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("vg")
	slug := testutil.RandomSlug("vg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccVLANGroupConsistencyLiteralNamesConfig(name, slug),
			},
		},
	})
}

func TestAccVLANGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlangrp-id")
	slug := testutil.GenerateSlug("tf-test-vlangrp-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVLANGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccVLANGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_vlan_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccVLANGroupConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_vlan_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccVLANGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_vlan_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}
