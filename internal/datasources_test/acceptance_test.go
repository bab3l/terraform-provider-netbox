package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_site" "test" {
  name   = "Test Site DS"
  slug   = "test-site-ds"
  status = "active"
}

data "netbox_site" "test" {
  slug = netbox_site.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "name", "Test Site DS"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "slug", "test-site-ds"),
				),
			},
		},
	})
}

func TestAccTenantDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_tenant" "test" {
  name = "Test Tenant DS"
  slug = "test-tenant-ds"
}

data "netbox_tenant" "test" {
  slug = netbox_tenant.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "name", "Test Tenant DS"),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "slug", "test-tenant-ds"),
				),
			},
		},
	})
}

func TestAccSiteGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_site_group" "test" {
  name = "Test Site Group DS"
  slug = "test-site-group-ds"
}

data "netbox_site_group" "test" {
  slug = netbox_site_group.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site_group.test", "name", "Test Site Group DS"),
					resource.TestCheckResourceAttr("data.netbox_site_group.test", "slug", "test-site-group-ds"),
				),
			},
		},
	})
}

func TestAccTenantGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_tenant_group" "test" {
  name = "Test Tenant Group DS"
  slug = "test-tenant-group-ds"
}

data "netbox_tenant_group" "test" {
  slug = netbox_tenant_group.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", "Test Tenant Group DS"),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", "test-tenant-group-ds"),
				),
			},
		},
	})
}

func TestAccManufacturerDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
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
  name = "Test Manufacturer DS"
  slug = "test-manufacturer-ds"
}

data "netbox_manufacturer" "test" {
  slug = netbox_manufacturer.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", "Test Manufacturer DS"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "slug", "test-manufacturer-ds"),
				),
			},
		},
	})
}

func TestAccPlatformDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_platform" "test" {
  name = "Test Platform DS"
  slug = "test-platform-ds"
}

data "netbox_platform" "test" {
  slug = netbox_platform.test.slug
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", "Test Platform DS"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "slug", "test-platform-ds"),
				),
			},
		},
	})
}
