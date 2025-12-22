package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterGroupResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group")

	slug := testutil.RandomSlug("tf-test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_cluster_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_ClusterGroup_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-lit")

	slug := testutil.RandomSlug("tf-test-cluster-group-lit")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterGroupConsistencyLiteralNamesConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},

			{

				Config: testAccClusterGroupConsistencyLiteralNamesConfig(name, slug),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
				),
			},
		},
	})

}

func testAccClusterGroupConsistencyLiteralNamesConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_group" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccClusterGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_group" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}
