package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-role")

	slug := testutil.RandomSlug("tf-test-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},

			{

				ResourceName: "netbox_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})

}

func TestAccRoleResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-role-full")
	slug := testutil.RandomSlug("tf-test-role-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated IPAM role description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_full(name, slug, description, 100, tagName1, tagSlug1, tagName2, tagSlug2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "100"),
					resource.TestCheckResourceAttr("netbox_role.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_full(name, slug, description, 100, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},

			{

				Config: testAccRoleResourceConfig_fullUpdate(name, slug, updatedDescription, 200, tagName1, tagSlug1, tagName2, tagSlug2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_role.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "200"),
					resource.TestCheckResourceAttr("netbox_role.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccRoleResourceConfig_fullUpdate(name, slug, updatedDescription, 200, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})

}

func TestAccRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-role-id")
	slug := testutil.RandomSlug("tf-test-role-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
				),
			},
			{
				Config:   testAccRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccRoleResourceConfig_full(name, slug, description string, weight int, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := fmt.Sprintf("test_field_%s", strings.ReplaceAll(slug, "-", "_"))
	return fmt.Sprintf(`

resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[9]q
  object_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[4]d

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}

`, name, slug, description, weight, tagName1, tagSlug1, tagName2, tagSlug2, cfName)

}

func testAccRoleResourceConfig_fullUpdate(name, slug, description string, weight int, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	cfName := fmt.Sprintf("test_field_%s", strings.ReplaceAll(slug, "-", "_"))
	return fmt.Sprintf(`

resource "netbox_tag" "tag1" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag2" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[9]q
  object_types = ["ipam.role"]
  type         = "text"
}

resource "netbox_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  weight      = %[4]d

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}

`, name, slug, description, weight, tagName1, tagSlug1, tagName2, tagSlug2, cfName)

}
func TestAccConsistency_Role_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-role-lit")
	slug := testutil.RandomSlug("tf-test-role-lit")
	description := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", description),
				),
			},
			{
				Config:   testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
				),
			},
		},
	})
}

func testAccRoleConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-role-del")
	slug := testutil.RandomSlug("tf-test-role-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamRolesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete role: %v", err)
					}
					t.Logf("Successfully externally deleted role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
