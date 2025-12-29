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

func TestAccTagResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTagResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")
	color := testutil.ColorOrange
	description := testutil.RandomName("description")
	updatedName := testutil.RandomName("tag-updated")
	updatedSlug := testutil.RandomSlug("tag-updated")
	updatedColor := "2196f3"
	updatedDescription := "Updated test tag description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceFull(name, slug, color, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", color),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", description),
				),
			},
			{
				Config: testAccTagResourceFull(updatedName, updatedSlug, updatedColor, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tag.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", updatedSlug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", updatedColor),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccTagResource_withObjectTypes(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tag")
	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceWithObjectTypes(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "object_types.#", "2"),
				),
			},
		},
	})
}

func TestAccTagResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-tag-id")
	slug := testutil.RandomSlug("tf-test-tag-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
				),
			},
		},
	})
}

func testAccTagResourceBasic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_Tag_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tag-lit")
	slug := testutil.RandomSlug("tag-lit")
	color := testutil.ColorOrange
	description := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagConsistencyLiteralNamesConfig(name, slug, color, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tag.test", "color", color),
					resource.TestCheckResourceAttr("netbox_tag.test", "description", description),
				),
			},
			{
				Config:   testAccTagConsistencyLiteralNamesConfig(name, slug, color, description),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
				),
			},
		},
	})
}

func testAccTagConsistencyLiteralNamesConfig(name, slug, color, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name        = %q
  slug        = %q
  color       = %q
  description = %q
}
`, name, slug, color, description)
}
func testAccTagResourceFull(name, slug, color, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name        = %q
  slug        = %q
  color       = %q
  description = %q
}
`, name, slug, color, description)
}

func testAccTagResourceWithObjectTypes(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name         = %q
  slug         = %q
  object_types = ["dcim.device", "dcim.site"]
}
`, name, slug)
}

func TestAccTagResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tag-del")
	slug := testutil.RandomSlug("tag-del")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceBasic(name, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),
					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.ExtrasAPI.ExtrasTagsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tag for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasTagsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tag: %v", err)
					}
					t.Logf("Successfully externally deleted tag with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
