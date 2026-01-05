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

func TestAccASNRangeDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asnrange-ds")
	slug := testutil.RandomSlug("tf-test-asnrange-ds")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(slug)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeDataSourceConfigByID(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_asn_range.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_id", "name", name),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_id", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_id", "start", "64512"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_id", "end", "64520"),
				),
			},
		},
	})
}

func TestAccASNRangeDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asnrange-ds")
	slug := testutil.RandomSlug("tf-test-asnrange-ds")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(slug)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeDataSourceConfigByName(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_asn_range.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_name", "name", name),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_name", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_name", "start", "64512"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.by_name", "end", "64520"),
				),
			},
		},
	})
}

func TestAccASNRangeDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-asnrange-ds")
	slug := testutil.RandomSlug("tf-test-asnrange-ds")
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(slug)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckASNRangeDestroy,
			testutil.CheckRIRDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeDataSourceConfig(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "slug", slug),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "start", "64512"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "end", "64520"),
				),
			},
		},
	})
}

func TestAccASNRangeDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("asnr-ds-id")
	slug := testutil.GenerateSlug(name)
	rirName := testutil.RandomName("rir-asnr")
	rirSlug := testutil.GenerateSlug(rirName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterASNRangeCleanup(slug)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccASNRangeDataSourceConfig(name, slug, rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_asn_range.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "slug", slug),
				),
			},
		},
	})
}

func testAccASNRangeDataSourceConfig(name, slug, rirName, rirSlug string) string {
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

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = 64512
  end   = 64520
}

data "netbox_asn_range" "test" {
  slug = netbox_asn_range.test.slug
}
`, rirName, rirSlug, name, slug)
}

func testAccASNRangeDataSourceConfigByID(name, slug, rirName, rirSlug string) string {
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

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = 64512
  end   = 64520
}

data "netbox_asn_range" "by_id" {
  id = netbox_asn_range.test.id
}
`, rirName, rirSlug, name, slug)
}

func testAccASNRangeDataSourceConfigByName(name, slug, rirName, rirSlug string) string {
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

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = 64512
  end   = 64520
}

data "netbox_asn_range" "by_name" {
  name = netbox_asn_range.test.name
}
`, rirName, rirSlug, name, slug)
}
