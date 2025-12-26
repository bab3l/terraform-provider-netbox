package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("service-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.0", "80"),
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
				),
			},
			{
				Config: testAccServiceTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated description"),
				),
			},
			{
				ResourceName:      "netbox_service_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("service-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "3"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}
`, name)
}

func testAccServiceTemplateResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "udp"
  ports       = [53, 123]
  description = "Updated description"
}
`, name+"-updated")
}

func TestAccServiceTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("service-template-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "tcp"
  ports       = [80, 443, 8080]
  description = "Test description"
  comments    = "Test comments"
}
`, name)
}
