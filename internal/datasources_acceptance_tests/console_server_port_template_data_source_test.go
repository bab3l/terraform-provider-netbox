package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortTemplateDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cspt-ds-id")
	manufacturerName := testutil.RandomName("tf-test-mfr-id")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-id")
	deviceTypeName := testutil.RandomName("tf-test-dt-id")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.by_id", "name", name),
				),
			},
		},
	})
}

func TestAccConsoleServerPortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-cspt-ds")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsoleServerPortTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					// Check by_id lookup
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.by_id", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.by_id", "device_type"),
					// Check by_device_type_and_name lookup
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.by_device_type_and_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.by_device_type_and_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_console_server_port_template.by_device_type_and_name", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port_template.by_device_type_and_name", "device_type"),
					// Verify both lookups return same template
					resource.TestCheckResourceAttrPair("data.netbox_console_server_port_template.by_id", "id", "data.netbox_console_server_port_template.by_device_type_and_name", "id"),
				),
			},
		},
	})
}

func testAccConsoleServerPortTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_console_server_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = "de-9"
}

data "netbox_console_server_port_template" "by_id" {
  id = netbox_console_server_port_template.test.id
}

data "netbox_console_server_port_template" "by_device_type_and_name" {
  device_type = netbox_device_type.test.id
  name        = netbox_console_server_port_template.test.name
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}
