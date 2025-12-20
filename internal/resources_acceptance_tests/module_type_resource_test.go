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

func TestAccModuleTypeResource_basic(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	model := testutil.RandomName("tf-test-module-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
				),
			},
			{
				ResourceName:            "netbox_module_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func TestAccModuleTypeResource_full(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	model := testutil.RandomName("tf-test-module-type-full")
	description := "Test module type with all fields"
	updatedDescription := "Updated module type description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_module_type.test", "part_number", "MT-001"),
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", description),
				),
			},
			{
				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_type.test", "description", updatedDescription),
				),
			},
		},
	})
}

func testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}
`, mfgName, mfgSlug, model)
}

func testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  part_number  = "MT-001"
  description  = %q
}
`, mfgName, mfgSlug, model, description)
}
