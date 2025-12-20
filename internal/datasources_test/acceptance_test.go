package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const invalidProviderData = "invalid"

func TestAccSiteDataSource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-site-ds")

	slug := testutil.RandomSlug("tf-test-site-ds")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_site.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_site.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_site.test", "slug", slug),
				),
			},
		},
	})

}

func testAccSiteDataSourceConfig(name, slug string) string {

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







resource "netbox_site" "test" {



  name   = %q



  slug   = %q



  status = "active"



}







data "netbox_site" "test" {



  slug = netbox_site.test.slug



}







`, name, slug)

}

func TestAccTenantDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tenant-ds")

	slug := testutil.RandomSlug("tf-test-tenant-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant.test", "slug", slug),
				),
			},
		},
	})

}

func testAccTenantDataSourceConfig(name, slug string) string {

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







resource "netbox_tenant" "test" {



  name = %q



  slug = %q



}







data "netbox_tenant" "test" {



  slug = netbox_tenant.test.slug



}







`, name, slug)

}

func TestAccSiteGroupDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-sg-ds")

	slug := testutil.RandomSlug("tf-test-sg-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_site_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_site_group.test", "slug", slug),
				),
			},
		},
	})

}

func testAccSiteGroupDataSourceConfig(name, slug string) string {

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







resource "netbox_site_group" "test" {



  name = %q



  slug = %q



}







data "netbox_site_group" "test" {



  slug = netbox_site_group.test.slug



}







`, name, slug)

}

func TestAccTenantGroupDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tg-ds")

	slug := testutil.RandomSlug("tf-test-tg-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})

}

func testAccTenantGroupDataSourceConfig(name, slug string) string {

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







resource "netbox_tenant_group" "test" {



  name = %q



  slug = %q



}







data "netbox_tenant_group" "test" {



  slug = netbox_tenant_group.test.slug



}







`, name, slug)

}

func TestAccManufacturerDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-mfr-ds")

	slug := testutil.RandomSlug("tf-test-mfr-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckManufacturerDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccManufacturerDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})

}

func testAccManufacturerDataSourceConfig(name, slug string) string {

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







data "netbox_manufacturer" "test" {



  slug = netbox_manufacturer.test.slug



}







`, name, slug)

}

func TestAccPlatformDataSource_basic(t *testing.T) {

	// Generate unique names for both manufacturer and platform

	// Platform requires a manufacturer, so we create both

	mfrName := testutil.RandomName("tf-test-mfr-for-plat-ds")

	mfrSlug := testutil.RandomSlug("tf-test-mfr-pds")

	platName := testutil.RandomName("tf-test-plat-ds")

	platSlug := testutil.RandomSlug("tf-test-plat-ds")

	// Register cleanup for both resources

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPlatformCleanup(platSlug)

	cleanup.RegisterManufacturerCleanup(mfrSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckPlatformDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccPlatformDataSourceConfig(platName, platSlug, mfrName, mfrSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", platName),

					resource.TestCheckResourceAttr("data.netbox_platform.test", "slug", platSlug),
				),
			},
		},
	})

}

func testAccPlatformDataSourceConfig(platName, platSlug, mfrName, mfrSlug string) string {

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







resource "netbox_manufacturer" "test_mfr" {



  name = %q



  slug = %q



}







resource "netbox_platform" "test" {



  name         = %q



  slug         = %q



  manufacturer = netbox_manufacturer.test_mfr.slug



}







data "netbox_platform" "test" {



  slug = netbox_platform.test.slug



}







`, mfrName, mfrSlug, platName, platSlug)

}

func TestAccRegionDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-region-ds")

	slug := testutil.RandomSlug("tf-test-region-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRegionDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_region.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_region.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_region.test", "slug", slug),
				),
			},
		},
	})

}

func testAccRegionDataSourceConfig(name, slug string) string {

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







resource "netbox_region" "test" {



  name = %q



  slug = %q



}







data "netbox_region" "test" {



  slug = netbox_region.test.slug



}







`, name, slug)

}

func TestAccLocationDataSource_basic(t *testing.T) {

	// Generate unique names

	siteName := testutil.RandomName("tf-test-loc-ds-site")

	siteSlug := testutil.RandomSlug("tf-test-loc-ds-s")

	name := testutil.RandomName("tf-test-location-ds")

	slug := testutil.RandomSlug("tf-test-location-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterLocationCleanup(slug)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccLocationDataSourceConfig(siteName, siteSlug, name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_location.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_location.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_location.test", "slug", slug),
				),
			},
		},
	})

}

func testAccLocationDataSourceConfig(siteName, siteSlug, name, slug string) string {

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







resource "netbox_site" "test" {



  name   = %q



  slug   = %q



  status = "active"



}







resource "netbox_location" "test" {



  name = %q



  slug = %q



  site = netbox_site.test.id



}







data "netbox_location" "test" {



  slug = netbox_location.test.slug



}







`, siteName, siteSlug, name, slug)

}

func TestAccRackDataSource_basic(t *testing.T) {

	// Generate unique names

	siteName := testutil.RandomName("tf-test-rack-ds-site")

	siteSlug := testutil.RandomSlug("tf-test-rack-ds-s")

	rackName := testutil.RandomName("tf-test-rack-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackCleanup(rackName)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),

		Steps: []resource.TestStep{

			{

				Config: testAccRackDataSourceConfig(siteName, siteSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_rack.test", "name", rackName),
				),
			},
		},
	})

}

func testAccRackDataSourceConfig(siteName, siteSlug, rackName string) string {

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







resource "netbox_site" "test" {



  name   = %q



  slug   = %q



  status = "active"



}







resource "netbox_rack" "test" {



  name = %q



  site = netbox_site.test.id



}







data "netbox_rack" "test" {



  name = netbox_rack.test.name



}







`, siteName, siteSlug, rackName)

}

func TestAccRackRoleDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-rackrole-ds")

	slug := testutil.RandomSlug("tf-test-rr-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackRoleDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_rack_role.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "slug", slug),
				),
			},
		},
	})

}

func testAccRackRoleDataSourceConfig(name, slug string) string {

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







resource "netbox_rack_role" "test" {



  name = %q



  slug = %q



}







data "netbox_rack_role" "test" {



  slug = netbox_rack_role.test.slug



}







`, name, slug)

}

func TestAccDeviceRoleDataSource_basic(t *testing.T) {

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

func TestAccDeviceTypeDataSource_basic(t *testing.T) {

	// Generate unique names

	model := testutil.RandomName("tf-test-devicetype-ds")

	slug := testutil.RandomSlug("tf-test-dt-ds")

	manufacturerName := testutil.RandomName("tf-test-mfr-ds")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceTypeCleanup(slug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeDataSourceConfig(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "manufacturer", manufacturerSlug),

					resource.TestCheckResourceAttr("data.netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func testAccDeviceTypeDataSourceConfig(model, slug, manufacturerName, manufacturerSlug string) string {

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



  manufacturer = netbox_manufacturer.test.slug



  model        = %q



  slug         = %q



}







data "netbox_device_type" "test" {



  slug = netbox_device_type.test.slug



}







`, manufacturerName, manufacturerSlug, model, slug)

}

// Route Target Data Source Tests

func TestAccRouteTargetDataSource_basic(t *testing.T) {

	// Generate unique name - route targets have 21 char max, use format like "65000:400-<random>"

	name := fmt.Sprintf("65000:400-%s", testutil.RandomSlug("ds")[:8])

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRouteTargetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRouteTargetDataSourceConfig(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_route_target.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_route_target.test", "name", name),
				),
			},
		},
	})

}

func testAccRouteTargetDataSourceConfig(name string) string {

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







resource "netbox_route_target" "test" {



  name = %q



}







data "netbox_route_target" "test" {



  name = netbox_route_target.test.name



}







`, name)

}

// Virtual Disk Data Source Tests

func TestAccVirtualDiskDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-vdisk-ds")

	clusterTypeName := testutil.RandomName("tf-test-ct")

	clusterTypeSlug := testutil.RandomSlug("tf-test-ct")

	clusterName := testutil.RandomName("tf-test-cluster")

	vmName := testutil.RandomName("tf-test-vm")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVirtualDiskCleanup(name)

	cleanup.RegisterVirtualMachineCleanup(vmName)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckVirtualDiskDestroy,

			testutil.CheckVirtualMachineDestroy,

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualDiskDataSourceConfig(name, clusterTypeName, clusterTypeSlug, clusterName, vmName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_virtual_disk.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_virtual_disk.test", "size", "100"),
				),
			},
		},
	})

}

func testAccVirtualDiskDataSourceConfig(name, clusterTypeName, clusterTypeSlug, clusterName, vmName string) string {

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







resource "netbox_cluster_type" "test" {



  name = %q



  slug = %q



}







resource "netbox_cluster" "test" {



  name = %q



  type = netbox_cluster_type.test.id



}







resource "netbox_virtual_machine" "test" {



  name    = %q



  cluster = netbox_cluster.test.id







  # Ignore disk changes since Netbox auto-computes this from virtual_disks







  lifecycle {







    ignore_changes = [disk]



  }



}







resource "netbox_virtual_disk" "test" {



  virtual_machine = netbox_virtual_machine.test.id



  name            = %q







  size            = 100



}







data "netbox_virtual_disk" "test" {







  id = netbox_virtual_disk.test.id



}







`, clusterTypeName, clusterTypeSlug, clusterName, vmName, name)

}

// ASN Range Data Source Tests

func TestAccASNRangeDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asnrange-ds")

	slug := testutil.RandomSlug("tf-test-asnrange-ds")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup

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

// Device Bay Template Data Source Tests

func TestAccDeviceBayTemplateDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-dbt-ds")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceBayTemplateCleanup(name)

	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceBayTemplateDestroy,

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttrSet("data.netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})

}

func testAccDeviceBayTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {

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



  model          = %q



  slug           = %q



  manufacturer   = netbox_manufacturer.test.slug



  subdevice_role = "parent"



}







resource "netbox_device_bay_template" "test" {







  device_type = netbox_device_type.test.id



  name        = %q



}







data "netbox_device_bay_template" "test" {







  id = netbox_device_bay_template.test.id



}







`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)

}
