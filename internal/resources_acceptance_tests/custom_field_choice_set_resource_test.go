package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: CustomFieldChoiceSet resources are safe to create with t.Parallel()
// because they do not have a "required" flag and do not persist state
// that interferes with other tests.

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldChoiceSetCleanupByName(name)

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldChoiceSetCleanupByName(name)
	cleanup.RegisterCustomFieldChoiceSetCleanupByName(updatedName)

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

func TestAccConsistency_CustomFieldChoiceSet_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldChoiceSetCleanupByName(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldChoiceSetConsistencyLiteralNamesConfig(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},

			{

				Config: testAccCustomFieldChoiceSetConsistencyLiteralNamesConfig(name),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_custom_field_choice_set.test", "id"),
				),
			},
		},
	})

}

func testAccCustomFieldChoiceSetConsistencyLiteralNamesConfig(name string) string {

	return fmt.Sprintf(`

resource "netbox_custom_field_choice_set" "test" {

  name                  = %q

  order_alphabetically = true

  extra_choices = [

    { value = "critical", label = "Critical" },

    { value = "high",     label = "High" },

    { value = "low",      label = "Low" },

  ]

}

`, name)

}

func TestAccCustomFieldChoiceSetResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("custom-field-choice-set-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field_choice_set.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),
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
