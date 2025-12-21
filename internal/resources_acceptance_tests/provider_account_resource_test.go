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

func TestAccProviderAccountResource_basic(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider")

	providerSlug := testutil.RandomSlug("tf-test-provider")

	accountID := testutil.RandomName("acct")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
				),
			},

			{

				ResourceName: "netbox_provider_account.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccProviderAccountResource_full(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider-full")

	providerSlug := testutil.RandomSlug("tf-test-provider-full")

	accountID := testutil.RandomName("acct")

	accountName := testutil.RandomName("tf-test-acct")

	description := "Test provider account with all fields"

	updatedDescription := "Updated provider account description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "name", accountName),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", description),
				),
			},

			{

				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

resource "netbox_provider_account" "test" {

  circuit_provider = netbox_provider.test.id

  account          = %q

}

`, providerName, providerSlug, accountID)

}

func testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

resource "netbox_provider_account" "test" {

  circuit_provider = netbox_provider.test.id

  account          = %q

  name             = %q

  description      = %q

}

`, providerName, providerSlug, accountID, accountName, description)

}
