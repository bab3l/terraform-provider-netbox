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

func TestAccRearPortTemplateResource_basic(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := testutil.RandomName("rear-port")

	portType := "8p8c"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),
				),
			},

			{

				ResourceName: "netbox_rear_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccRearPortTemplateResource_full(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	name := testutil.RandomName("rear-port")

	portType := "lc"

	label := testutil.RandomName("label")

	color := "aa1409"

	description := testutil.RandomName("description")

	positions := int32(2)

	updatedName := testutil.RandomName("rear-port")

	updatedLabel := testutil.RandomName("label")

	updatedDescription := testutil.RandomName("description")

	updatedPositions := int32(4)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description, positions),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "color", color),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", positions)),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", description),
				),
			},

			{

				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, portType, updatedLabel, color, updatedDescription, updatedPositions),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", updatedPositions)),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", updatedDescription),
				),
			},

			{

				ResourceName: "netbox_rear_port_template.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})

}

func TestAccConsistency_RearPortTemplate(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("devicetype")

	deviceTypeSlug := testutil.RandomSlug("devicetype")

	portName := testutil.RandomName("rear-port")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})

}

func TestAccConsistency_RearPortTemplate_LiteralNames(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	resourceName := testutil.RandomName("rear_port")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", resourceName),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeSlug),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})

}

func testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType string) string {

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

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  type        = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType)

}

func testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description string, positions int32) string {

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

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %q

  type        = %q

  label       = %q

  color       = %q

  positions   = %d

  description = %q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, positions, description)

}

func testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}

resource "netbox_device_type" "test" {

  model = "%[3]s"

  slug = "%[4]s"

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.model

  name = "%[5]s"

  type = "8p8c"

  positions = 1

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)

}

func testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  model        = %q

  slug         = %q

  manufacturer = netbox_manufacturer.test.id

}

resource "netbox_rear_port_template" "test" {

  device_type = %q

  name = %q

  type = "8p8c"

  positions = 1

  depends_on = [netbox_device_type.test]

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)

}
