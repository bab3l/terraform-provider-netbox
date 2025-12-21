package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Cable resource is simpler than other resources (e.g., ASN, Aggregate) because:

// - Only a_terminations and b_terminations are required (complex nested objects)

// - Other fields are simple scalars (type, status, color, label, description, comments, etc.)

// - No complex reference validation or state drift issues

// Therefore, a single comprehensive test that validates core functionality (creating a cable

// with terminations and import/export) is sufficient to ensure the resource works correctly.

func TestAccCableResource_basic(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("test-site-cable")

	siteSlug := testutil.GenerateSlug(siteName)

	deviceName := testutil.RandomName("test-device-cable")

	interfaceNameA := "eth2"

	interfaceNameB := "eth1"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),

					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},

			{

				ResourceName: "netbox_cable.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccCableResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %q

  slug = %q

  status = "active"

}

resource "netbox_manufacturer" "test" {

  name = "Test Manufacturer"

  slug = "test-manufacturer"

}

resource "netbox_device_role" "test" {

  name = "Test Device Role"

  slug = "test-device-role"

}

resource "netbox_device_type" "test" {

  model = "Test Device Type"

  slug  = "test-device-type"

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_device" "test_a" {

  name           = "%s-a"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_device" "test_b" {

  name           = "%s-b"

  device_type    = netbox_device_type.test.id

  role           = netbox_device_role.test.id

  site           = netbox_site.test.id

}

resource "netbox_interface" "test_a" {

  name      = %q

  device    = netbox_device.test_a.id

  type      = "1000base-t"

}

resource "netbox_interface" "test_b" {

  name      = %q

  device    = netbox_device.test_b.id

  type      = "1000base-t"

}

resource "netbox_cable" "test" {

  status = "connected"

  type   = "cat6"

  a_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_a.id

    }

  ]

  b_terminations = [

    {

      object_type = "dcim.interface"

      object_id   = netbox_interface.test_b.id

    }

  ]

}

`, siteName, siteSlug, deviceName, deviceName, interfaceNameA, interfaceNameB)

}
