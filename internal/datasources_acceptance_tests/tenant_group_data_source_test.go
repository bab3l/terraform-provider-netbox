package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantGroupDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tg-ds")

	slug := testutil.RandomSlug("tf-test-tg-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})

}

func testAccTenantGroupDataSourceConfig(name, slug string) string {

	return fmt.Sprintf(`

terraform {

  required_providers {

    netbox = {

      source = "bab3l/netbox"

      version = ">= 0.1.0"

    }

  }

}

provider "netbox" {}

resource "netbox_tenant_group" "test" {

  name = %q

  slug = %q

}

data "netbox_tenant_group" "test" {

  slug = netbox_tenant_group.test.slug

}

`, name, slug)

}
