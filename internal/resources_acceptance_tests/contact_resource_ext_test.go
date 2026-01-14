package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccContactResource_removeOptionalFields_extended tests adding and removing optional contact fields.
func TestAccContactResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	contactName := testutil.RandomName("tf-test-contact-rem")

	testFields := map[string]string{
		"address": "123 Main St, City, State 12345",
		"email":   "test@example.com",
		"link":    "https://example.com",
		"phone":   "+1234567890",
		"title":   "Network Engineer",
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_contact",
		BaseConfig: func() string {
			return testAccContactResourceConfig_removeOptionalFields_base(contactName)
		},
		ConfigWithFields: func() string {
			return testAccContactResourceConfig_removeOptionalFields_withFields(contactName, testFields)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": contactName,
		},
	})
}

func testAccContactResourceConfig_removeOptionalFields_base(contactName string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name = %[1]q
}
`, contactName)
}

func testAccContactResourceConfig_removeOptionalFields_withFields(contactName string, fields map[string]string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name    = %[1]q
  address = %[2]q
  email   = %[3]q
  link    = %[4]q
  phone   = %[5]q
  title   = %[6]q
}
`, contactName, fields["address"], fields["email"], fields["link"], fields["phone"], fields["title"])
}
