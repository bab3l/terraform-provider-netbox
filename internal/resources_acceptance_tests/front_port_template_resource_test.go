package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const frontPortTypeStandard = "8p8c"

func TestAccFrontPortTemplateResource_basic(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	frontPortName := testutil.RandomName("front-port")
	portType := frontPortTypeStandard
	rearPortName := testutil.RandomName("rear-port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", portType),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port", rearPortName),
				),
			},
			{
				ResourceName:            "netbox_front_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type", "rear_port"},
			},
		},
	})
}

func TestAccFrontPortTemplateResource_full(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	frontPortName := testutil.RandomName("front-port")
	portType := "lc"
	label := testutil.RandomName("label")
	description := testutil.RandomName("description")
	rearPortPosition := int32(1)
	updatedFrontPortName := testutil.RandomName("front-port-upd")
	updatedLabel := testutil.RandomName("label-upd")
	updatedDescription := testutil.RandomName("description-upd")
	updatedRearPortPosition := int32(2)
	rearPortName := testutil.RandomName("rear-port")
	color := testutil.Color

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType, label, color, description, rearPortPosition),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", portType),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "color", color),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", fmt.Sprintf("%d", rearPortPosition)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", description),
				),
			},
			{
				Config: testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, updatedFrontPortName, portType, updatedLabel, color, updatedDescription, updatedRearPortPosition),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", updatedFrontPortName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "rear_port_position", fmt.Sprintf("%d", updatedRearPortPosition)),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "description", updatedDescription),
				),
			},
			{
				ResourceName:            "netbox_front_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type", "rear_port", "display_name"},
			},
		},
	})
}

func TestAccConsistency_FrontPortTemplate_LiteralNames(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	resourceName := testutil.RandomName("front_port")
	rearPortTemplateName := testutil.RandomName("rear-port")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName, rearPortTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName, rearPortTemplateName),
			},
		},
	})
}

func TestAccFrontPortTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-id")
	manufacturerSlug := testutil.RandomSlug("mfr-id")
	deviceTypeName := testutil.RandomName("dt-id")
	deviceTypeSlug := testutil.RandomSlug("dt-id")
	frontPortName := testutil.RandomName("front-port-id")
	portType := frontPortTypeStandard
	rearPortName := testutil.RandomName("rear-port-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port_template.test", "type", portType),
				),
			},
		},
	})
}

func testAccFrontPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType string) string {
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
  positions   = 2
}

resource "netbox_front_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
  rear_port   = netbox_rear_port_template.test.name
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, portType, frontPortName, portType)
}

func testAccFrontPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, frontPortName, portType, label, color, description string, rearPortPosition int32) string {
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
  positions   = 4
}

resource "netbox_front_port_template" "test" {
  device_type        = netbox_device_type.test.id
  name               = %q
  type               = %q
  rear_port          = netbox_rear_port_template.test.name
  rear_port_position = %d
  label              = %q
  color              = %q
  description        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, portType, frontPortName, portType, rearPortPosition, label, color, description)
}

func testAccFrontPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName, rearPortTemplateName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "rear" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
  type        = "8p8c"
  positions   = 1
}

resource "netbox_front_port_template" "test" {
  device_type = %[4]q
  name = %[6]q
  type = "8p8c"
  rear_port = netbox_rear_port_template.rear.name
  rear_port_position = 1

  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortTemplateName, resourceName)
}
