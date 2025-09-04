package resources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManufacturerResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
		  terraform {		import "github.com/bab3l/terraform-provider-netbox/internal/resources"
			required_providers {
			  netbox = {
				source = "bab3l/netbox"
				version = ">= 0.1.0"
			  }
			}
		  }

		  provider "netbox" {}

		  resource "netbox_manufacturer" "test" {
			name = "Test Manufacturer"
			slug = "test-manufacturer"
		  }
		`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("netbox_manufacturer.test", "name", "Test Manufacturer"),
				),
			},
		},
	})
}
