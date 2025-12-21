package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateResource_basic(t *testing.T) {

	t.Parallel()
	name := acctest.RandomWithPrefix("test-service-template")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
	name := acctest.RandomWithPrefix("test-service-template")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
