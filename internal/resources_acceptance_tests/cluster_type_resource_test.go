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

func TestAccClusterTypeResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-cluster-type")

	slug := testutil.RandomSlug("tf-test-cluster-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccClusterTypeResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-cluster-type-full")

	slug := testutil.RandomSlug("tf-test-cluster-type-full")

	description := "Test cluster type with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterTypeResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", description),
				),
			},
		},
	})

}

func TestAccClusterTypeResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-cluster-type-update")

	slug := testutil.RandomSlug("tf-test-cluster-type-update")

	updatedName := testutil.RandomName("tf-test-cluster-type-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
				),
			},

			{

				Config: testAccClusterTypeResourceConfig_full(updatedName, slug, "Updated description"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", "Updated description"),
				),
			},
		},
	})

}

func TestAccClusterTypeResource_import(t *testing.T) {

	name := testutil.RandomName("tf-test-cluster-type-import")

	slug := testutil.RandomSlug("tf-test-cluster-type-import")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_cluster_type.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccClusterTypeResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_cluster_type" "test" {



  name = %q



  slug = %q



}



`, name, slug)

}

func testAccClusterTypeResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`



resource "netbox_cluster_type" "test" {



  name        = %q



  slug        = %q



  description = %q



}



`, name, slug, description)

}
