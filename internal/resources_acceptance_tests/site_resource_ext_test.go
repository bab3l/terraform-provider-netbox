package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestAccSiteResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-ext")
	slug := testutil.RandomSlug("tf-test-site-ext")
	facility := testutil.RandomName("facility")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_site",
		BaseConfig: func() string {
			return testAccSiteResourceConfig_minimal(name, slug)
		},
		ConfigWithFields: func() string {
			return testAccSiteResourceConfig_withFacility(name, slug, facility)
		},
		OptionalFields: map[string]string{
			"facility": facility,
		},
		RequiredFields: map[string]string{
			"name": name,
			"slug": slug,
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}

func testAccSiteResourceConfig_withFacility(name, slug, facility string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name     = %[1]q
  slug     = %[2]q
  status   = "active"
  facility = %[3]q
}
`, name, slug, facility)
}
