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

func TestAccCustomLinkDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("cl-id")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_custom_link.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "name", name),
				),
			},
		},
	})
}

func TestAccCustomLinkDataSource_byID(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "link_text", "View Details"),
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "object_types.#", "1"),
				),
			},
		},
	})
}

func TestAccCustomLinkDataSource_byName(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkDataSourceConfig_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_custom_link.test", "id"),
				),
			},
		},
	})
}

func testAccCustomLinkDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
}

data "netbox_custom_link" "test" {
  id = netbox_custom_link.test.id
}
`, name)
}

func testAccCustomLinkDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
}

data "netbox_custom_link" "test" {
  name = netbox_custom_link.test.name
}
`, name)
}
