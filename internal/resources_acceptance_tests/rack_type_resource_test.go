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

func TestAccRackTypeResource_basic(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	model := testutil.RandomName("tf-test-rack-type")

	slug := testutil.RandomSlug("tf-test-rack-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "manufacturer"),
				),
			},

			{

				ResourceName: "netbox_rack_type.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})

}

func TestAccRackTypeResource_full(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-full")

	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")

	model := testutil.RandomName("tf-test-rack-type-full")

	slug := testutil.RandomSlug("tf-test-rack-type-full")

	description := "Test rack type with all fields"

	updatedDescription := "Updated rack type description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description, 42, 19),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "42"),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "width", "19"),
				),
			},

			{

				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, updatedDescription, 48, 19),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "48"),
				),
			},
		},
	})

}

func TestAccConsistency_RackType_LiteralNames(t *testing.T) {

	t.Parallel()

	t.Parallel()

	mfgName := testutil.RandomName("manufacturer")

	mfgSlug := testutil.RandomSlug("manufacturer")

	model := testutil.RandomName("rack-type")

	slug := testutil.RandomSlug("rack-type")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "manufacturer", mfgName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),
			},
		},
	})

}

func testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_rack_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

  form_factor  = "4-post-cabinet"

}

`, mfgName, mfgSlug, model, slug)

}

func testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description string, uHeight, width int) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_rack_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

  description  = %q

  u_height     = %d

  width        = %d

  form_factor  = "4-post-cabinet"

}

`, mfgName, mfgSlug, model, slug, description, uHeight, width)

}

func testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_rack_type" "test" {

  # Use literal string name to mimic existing user state

  manufacturer = %q

  model        = %q

  slug         = %q

  u_height     = 42

  width        = 19

  form_factor  = "4-post-cabinet"

  depends_on = [netbox_manufacturer.test]

}

`, mfgName, mfgSlug, mfgName, model, slug)

}
