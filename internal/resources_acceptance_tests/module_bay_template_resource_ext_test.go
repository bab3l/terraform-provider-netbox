package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestAccModuleBayTemplateResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-mbt-ext")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-mbt-ext")
	dtModel := testutil.RandomName("tf-test-dt-mbt-ext")
	dtSlug := testutil.RandomSlug("tf-test-dt-mbt-ext")
	templateName := testutil.RandomName("tf-test-mbt-ext")
	position := testutil.RandomName("pos")
	description := testutil.RandomName("desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_module_bay_template",
		BaseConfig: func() string {
			return testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName)
		},
		ConfigWithFields: func() string {
			return testAccModuleBayTemplateResourceConfig_withPositionAndDescription(mfgName, mfgSlug, dtModel, dtSlug, templateName, position, description)
		},
		OptionalFields: map[string]string{
			"position":    position,
			"description": description,
		},
		RequiredFields: map[string]string{
			"name": templateName,
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckModuleBayTemplateDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
	})
}

func testAccModuleBayTemplateResourceConfig_withPositionAndDescription(mfgName, mfgSlug, dtModel, dtSlug, templateName, position, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_module_bay_template" "test" {
  device_type  = netbox_device_type.test.id
  name         = %q
  position     = %q
  description  = %q
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName, position, description)
}
