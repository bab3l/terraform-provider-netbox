package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceRoleDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-devicerole-ds-id")
	slug := testutil.RandomSlug("tf-test-dr-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckDeviceRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "name", name),
				),
			},
		},
	})
}

func TestAccDeviceRoleDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-devicerole-ds")
	slug := testutil.RandomSlug("tf-test-dr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckDeviceRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "vm_role", "true"),
				),
			},
		},
	})
}

func testAccDeviceRoleDataSourceConfig(name, slug string) string {
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

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

data "netbox_device_role" "test" {
  slug = netbox_device_role.test.slug
}
`, name, slug)

}
func TestAccDeviceRoleDataSource_byName(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-devicerole-ds")
	slug := testutil.RandomSlug("tf-test-dr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckDeviceRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccDeviceRoleDataSourceConfigByName(name, slug string) string {
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

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

data "netbox_device_role" "test" {
  name = netbox_device_role.test.name
}
`, name, slug)
}

func TestAccDeviceRoleDataSource_byID(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-devicerole-ds")
	slug := testutil.RandomSlug("tf-test-dr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckDeviceRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleDataSourceConfigByID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_device_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccDeviceRoleDataSourceConfigByID(name, slug string) string {
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

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

data "netbox_device_role" "test" {
  id = netbox_device_role.test.id
}
`, name, slug)
}
