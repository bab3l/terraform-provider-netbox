package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFieldResource_basic(t *testing.T) {

	t.Parallel()

	// Custom field names can only contain alphanumeric characters and underscores

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "text"),
				),
			},

			{

				ResourceName: "netbox_custom_field.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCustomFieldResource_full(t *testing.T) {

	t.Parallel()

	// Custom field names can only contain alphanumeric characters and underscores

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	description := "Test custom field with all fields"

	updatedDescription := "Updated custom field description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldResourceConfig_full(name, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "integer"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", description),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "required", "true"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_minimum", "1"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_maximum", "100"),
				),
			},

			{

				Config: testAccCustomFieldResourceConfig_full(name, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccCustomFieldResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "test" {

  name         = %q

  type         = "text"

  object_types = ["dcim.site"]

}

`, name)

}

func testAccCustomFieldResourceConfig_full(name, description string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field" "test" {

  name               = %q

  type               = "integer"

  object_types       = ["dcim.site", "dcim.device"]

  description        = %q

  required           = true

  validation_minimum = 1

  validation_maximum = 100

  weight             = 50

}

`, name, description)

}
