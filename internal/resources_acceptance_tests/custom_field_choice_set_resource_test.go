package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFieldChoiceSetResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldChoiceSetResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},

			{

				ResourceName: "netbox_custom_field_choice_set.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCustomFieldChoiceSetResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldChoiceSetResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "description", "Test choice set"),

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "order_alphabetically", "true"),

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},
		},
	})

}

func TestAccCustomFieldChoiceSetResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs")

	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldChoiceSetResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),
				),
			},

			{

				Config: testAccCustomFieldChoiceSetResourceConfig_basic(updatedName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", updatedName),
				),
			},
		},
	})

}

func testAccCustomFieldChoiceSetResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field_choice_set" "test" {

  name = "%s"

  extra_choices = [

    { value = "opt1", label = "Option 1" },

    { value = "opt2", label = "Option 2" },

    { value = "opt3", label = "Option 3" },

  ]

}

`, name)

}

func testAccCustomFieldChoiceSetResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field_choice_set" "test" {

  name                 = "%s"

  description          = "Test choice set"

  order_alphabetically = true

  extra_choices = [

    { value = "critical", label = "Critical" },

    { value = "high",     label = "High" },

    { value = "low",      label = "Low" },

  ]

}

`, name)

}
