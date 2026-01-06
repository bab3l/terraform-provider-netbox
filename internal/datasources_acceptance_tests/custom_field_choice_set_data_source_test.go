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

func TestAccCustomFieldChoiceSetDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cfcs-ds-id")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_custom_field_choice_set.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "name", name),
				),
			},
		},
	})
}

func TestAccCustomFieldChoiceSetDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},
		},
	})
}

func TestAccCustomFieldChoiceSetDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetDataSourceConfig_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_custom_field_choice_set.test", "id"),
				),
			},
		},
	})
}

func testAccCustomFieldChoiceSetDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name = "%s"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}

data "netbox_custom_field_choice_set" "test" {
  id = netbox_custom_field_choice_set.test.id
}
`, name)
}

func testAccCustomFieldChoiceSetDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name = "%s"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}

data "netbox_custom_field_choice_set" "test" {
  name = netbox_custom_field_choice_set.test.name
}
`, name)
}
