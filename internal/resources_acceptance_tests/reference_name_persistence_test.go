package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReferenceNamePersistence_Device(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	locationName := testutil.RandomName("location")

	locationSlug := testutil.RandomSlug("location")

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	deviceRoleName := testutil.RandomName("device-role")

	deviceRoleSlug := testutil.RandomSlug("device-role")

	deviceName := testutil.RandomName("device")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceReferenceNameConfig(siteName, siteSlug, locationName, locationSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					resource.TestCheckResourceAttr("netbox_device.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_device.test", "location", locationName),

					resource.TestCheckResourceAttr("netbox_device.test", "device_type", deviceTypeName),

					resource.TestCheckResourceAttr("netbox_device.test", "role", deviceRoleName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccDeviceReferenceNameConfig(siteName, siteSlug, locationName, locationSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName),
			},
		},
	})

}

func testAccDeviceReferenceNameConfig(siteName, siteSlug, locationName, locationSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_location" "test" {

  name = "%[3]s"

  slug = "%[4]s"

  site = netbox_site.test.id

}



resource "netbox_manufacturer" "test" {

  name = "%[5]s"

  slug = "%[6]s"

}



resource "netbox_device_type" "test" {

  model = "%[7]s"

  slug  = "%[8]s"

  manufacturer = netbox_manufacturer.test.id

  u_height = 1

}



resource "netbox_device_role" "test" {

  name = "%[9]s"

  slug = "%[10]s"

}



resource "netbox_device" "test" {

  name = "%[11]s"

  site = netbox_site.test.name

  location = netbox_location.test.name

  device_type = netbox_device_type.test.model

  role = netbox_device_role.test.name

}

`, siteName, siteSlug, locationName, locationSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName)

}

func TestAccReferenceNamePersistence_Rack(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	rackTypeName := testutil.RandomName("rack-type")

	rackTypeSlug := testutil.RandomSlug("rack-type")

	rackRoleName := testutil.RandomName("rack-role")

	rackRoleSlug := testutil.RandomSlug("rack-role")

	rackName := testutil.RandomName("rack")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRackReferenceNameConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug, rackRoleName, rackRoleSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),

					resource.TestCheckResourceAttr("netbox_rack.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_rack.test", "rack_type", rackTypeName),

					resource.TestCheckResourceAttr("netbox_rack.test", "role", rackRoleName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRackReferenceNameConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug, rackRoleName, rackRoleSlug, rackName),
			},
		},
	})

}

func testAccRackReferenceNameConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug, rackRoleName, rackRoleSlug, rackName string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_manufacturer" "test" {

  name = "%[3]s"

  slug = "%[4]s"

}



resource "netbox_rack_type" "test" {

  model = "%[5]s"

  slug  = "%[6]s"

  manufacturer = netbox_manufacturer.test.id

  form_factor = "4-post-cabinet"

}



resource "netbox_rack_role" "test" {

  name = "%[7]s"

  slug = "%[8]s"

}



resource "netbox_rack" "test" {

  name = "%[9]s"

  site = netbox_site.test.name

  rack_type = netbox_rack_type.test.model

  role = netbox_rack_role.test.name

  status = "active"

  width = "19"

  u_height = 42

}

`, siteName, siteSlug, manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug, rackRoleName, rackRoleSlug, rackName)

}

func TestAccReferenceNamePersistence_Site(t *testing.T) {

	t.Parallel()

	regionName := testutil.RandomName("region")

	regionSlug := testutil.RandomSlug("region")

	siteGroupName := testutil.RandomName("site-group")

	siteGroupSlug := testutil.RandomSlug("site-group")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteReferenceNameConfig(regionName, regionSlug, siteGroupName, siteGroupSlug, tenantName, tenantSlug, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),

					resource.TestCheckResourceAttr("netbox_site.test", "region", regionName),

					resource.TestCheckResourceAttr("netbox_site.test", "group", siteGroupName),

					resource.TestCheckResourceAttr("netbox_site.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccSiteReferenceNameConfig(regionName, regionSlug, siteGroupName, siteGroupSlug, tenantName, tenantSlug, siteName, siteSlug),
			},
		},
	})

}

func testAccSiteReferenceNameConfig(regionName, regionSlug, siteGroupName, siteGroupSlug, tenantName, tenantSlug, siteName, siteSlug string) string {

	return fmt.Sprintf(`

resource "netbox_region" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_site_group" "test" {

  name = "%[3]s"

  slug = "%[4]s"

}



resource "netbox_tenant" "test" {

  name = "%[5]s"

  slug = "%[6]s"

}



resource "netbox_site" "test" {

  name = "%[7]s"

  slug = "%[8]s"

  region = netbox_region.test.name

  group = netbox_site_group.test.name

  tenant = netbox_tenant.test.name

}

`, regionName, regionSlug, siteGroupName, siteGroupSlug, tenantName, tenantSlug, siteName, siteSlug)

}

func TestAccReferenceNamePersistence_CircuitGroup_Tenant(t *testing.T) {

	t.Parallel()

	tenantGroupName := testutil.RandomName("tenant-group")

	tenantGroupSlug := testutil.RandomSlug("tenant-group")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	circuitGroupName := testutil.RandomName("circuit-group")

	circuitGroupSlug := testutil.RandomSlug("circuit-group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupReferenceNameConfig(tenantGroupName, tenantGroupSlug, tenantName, tenantSlug, circuitGroupName, circuitGroupSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", circuitGroupName),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_tenant.test", "group", tenantGroupName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccCircuitGroupReferenceNameConfig(tenantGroupName, tenantGroupSlug, tenantName, tenantSlug, circuitGroupName, circuitGroupSlug),
			},
		},
	})

}

func testAccCircuitGroupReferenceNameConfig(tenantGroupName, tenantGroupSlug, tenantName, tenantSlug, circuitGroupName, circuitGroupSlug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant_group" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_tenant" "test" {

  name = "%[3]s"

  slug = "%[4]s"

  group = netbox_tenant_group.test.name

}



resource "netbox_circuit_group" "test" {

  name = "%[5]s"

  slug = "%[6]s"

  tenant = netbox_tenant.test.name

}

`, tenantGroupName, tenantGroupSlug, tenantName, tenantSlug, circuitGroupName, circuitGroupSlug)

}

func TestAccReferenceNamePersistence_ContactGroup(t *testing.T) {

	t.Parallel()

	parentName := testutil.RandomName("parent-contact-group")

	parentSlug := testutil.RandomSlug("parent-contact-group")

	childName := testutil.RandomName("child-contact-group")

	childSlug := testutil.RandomSlug("child-contact-group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.child", "name", childName),

					resource.TestCheckResourceAttr("netbox_contact_group.child", "parent", parentName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccContactGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),
			},
		},
	})

}

func testAccContactGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_group" "parent" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_contact_group" "child" {

  name = "%[3]s"

  slug = "%[4]s"

  parent = netbox_contact_group.parent.name

}

`, parentName, parentSlug, childName, childSlug)

}

func TestAccReferenceNamePersistence_Region(t *testing.T) {

	t.Parallel()

	parentName := testutil.RandomName("parent-region")

	parentSlug := testutil.RandomSlug("parent-region")

	childName := testutil.RandomName("child-region")

	childSlug := testutil.RandomSlug("child-region")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRegionReferenceNameConfig(parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_region.child", "name", childName),

					resource.TestCheckResourceAttr("netbox_region.child", "parent", parentName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRegionReferenceNameConfig(parentName, parentSlug, childName, childSlug),
			},
		},
	})

}

func testAccRegionReferenceNameConfig(parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

resource "netbox_region" "parent" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_region" "child" {

  name = "%[3]s"

  slug = "%[4]s"

  parent = netbox_region.parent.name

}

`, parentName, parentSlug, childName, childSlug)

}

func TestAccReferenceNamePersistence_SiteGroup(t *testing.T) {

	t.Parallel()

	parentName := testutil.RandomName("parent-site-group")

	parentSlug := testutil.RandomSlug("parent-site-group")

	childName := testutil.RandomName("child-site-group")

	childSlug := testutil.RandomSlug("child-site-group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_site_group.child", "name", childName),

					resource.TestCheckResourceAttr("netbox_site_group.child", "parent", parentName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccSiteGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),
			},
		},
	})

}

func testAccSiteGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

resource "netbox_site_group" "parent" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_site_group" "child" {

  name = "%[3]s"

  slug = "%[4]s"

  parent = netbox_site_group.parent.name

}

`, parentName, parentSlug, childName, childSlug)

}

func TestAccReferenceNamePersistence_TenantGroup(t *testing.T) {

	t.Parallel()

	parentName := testutil.RandomName("parent-tenant-group")

	parentSlug := testutil.RandomSlug("parent-tenant-group")

	childName := testutil.RandomName("child-tenant-group")

	childSlug := testutil.RandomSlug("child-tenant-group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tenant_group.child", "name", childName),

					resource.TestCheckResourceAttr("netbox_tenant_group.child", "parent", parentName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccTenantGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug),
			},
		},
	})

}

func testAccTenantGroupReferenceNameConfig(parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant_group" "parent" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_tenant_group" "child" {

  name = "%[3]s"

  slug = "%[4]s"

  parent = netbox_tenant_group.parent.name

}

`, parentName, parentSlug, childName, childSlug)

}

func TestAccReferenceNamePersistence_Location(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	parentName := testutil.RandomName("parent-location")

	parentSlug := testutil.RandomSlug("parent-location")

	childName := testutil.RandomName("child-location")

	childSlug := testutil.RandomSlug("child-location")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccLocationReferenceNameConfig(siteName, siteSlug, parentName, parentSlug, childName, childSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_location.child", "name", childName),

					resource.TestCheckResourceAttr("netbox_location.child", "parent", parentName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccLocationReferenceNameConfig(siteName, siteSlug, parentName, parentSlug, childName, childSlug),
			},
		},
	})

}

func testAccLocationReferenceNameConfig(siteName, siteSlug, parentName, parentSlug, childName, childSlug string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_location" "parent" {

  name = "%[3]s"

  slug = "%[4]s"

  site = netbox_site.test.id

}



resource "netbox_location" "child" {

  name = "%[5]s"

  slug = "%[6]s"

  site = netbox_site.test.id

  parent = netbox_location.parent.name

}

`, siteName, siteSlug, parentName, parentSlug, childName, childSlug)

}

func TestAccReferenceNamePersistence_VRF(t *testing.T) {

	t.Parallel()

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	vrfName := testutil.RandomName("vrf")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVRFReferenceNameConfig(tenantName, tenantSlug, vrfName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_vrf.test", "name", vrfName),

					resource.TestCheckResourceAttr("netbox_vrf.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccVRFReferenceNameConfig(tenantName, tenantSlug, vrfName),
			},
		},
	})

}

func testAccVRFReferenceNameConfig(tenantName, tenantSlug, vrfName string) string {

	return fmt.Sprintf(`

resource "netbox_tenant" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_vrf" "test" {

  name = "%[3]s"

  tenant = netbox_tenant.test.name

}

`, tenantName, tenantSlug, vrfName)

}

func TestAccReferenceNamePersistence_RackType(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	rackTypeName := testutil.RandomName("rack-type")

	rackTypeSlug := testutil.RandomSlug("rack-type")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRackTypeReferenceNameConfig(manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", rackTypeName),

					resource.TestCheckResourceAttr("netbox_rack_type.test", "manufacturer", manufacturerName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRackTypeReferenceNameConfig(manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug),
			},
		},
	})

}

func testAccRackTypeReferenceNameConfig(manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_rack_type" "test" {

  model = "%[3]s"

  slug  = "%[4]s"

  manufacturer = netbox_manufacturer.test.name

  form_factor = "4-post-cabinet"

}

`, manufacturerName, manufacturerSlug, rackTypeName, rackTypeSlug)

}

func TestAccReferenceNamePersistence_RouteTarget(t *testing.T) {

	t.Parallel()

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	routeTargetName := "65000:1"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRouteTargetReferenceNameConfig(tenantName, tenantSlug, routeTargetName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", routeTargetName),

					resource.TestCheckResourceAttr("netbox_route_target.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRouteTargetReferenceNameConfig(tenantName, tenantSlug, routeTargetName),
			},
		},
	})

}

func testAccRouteTargetReferenceNameConfig(tenantName, tenantSlug, routeTargetName string) string {

	return fmt.Sprintf(`

resource "netbox_tenant" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_route_target" "test" {

  name = "%[3]s"

  tenant = netbox_tenant.test.name

}

`, tenantName, tenantSlug, routeTargetName)

}

func TestAccReferenceNamePersistence_VLAN(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	vlanName := testutil.RandomName("vlan")

	vlanVid := 100

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANReferenceNameConfig(siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),

					resource.TestCheckResourceAttr("netbox_vlan.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_vlan.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccVLANReferenceNameConfig(siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
			},
		},
	})

}

func testAccVLANReferenceNameConfig(siteName, siteSlug, tenantName, tenantSlug, vlanName string, vlanVid int) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_tenant" "test" {

  name = "%[3]s"

  slug = "%[4]s"

}



resource "netbox_vlan" "test" {

  name = "%[5]s"

  vid  = %[6]d

  site = netbox_site.test.name

  tenant = netbox_tenant.test.name

}

`, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid)

}

func TestAccReferenceNamePersistence_Platform(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	platformName := testutil.RandomName("platform")

	platformSlug := testutil.RandomSlug("platform")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccPlatformReferenceNameConfig(manufacturerName, manufacturerSlug, platformName, platformSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),

					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccPlatformReferenceNameConfig(manufacturerName, manufacturerSlug, platformName, platformSlug),
			},
		},
	})

}

func testAccPlatformReferenceNameConfig(manufacturerName, manufacturerSlug, platformName, platformSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = "%[1]s"

  slug = "%[2]s"

}



resource "netbox_platform" "test" {

  name = "%[3]s"

  slug = "%[4]s"

  manufacturer = netbox_manufacturer.test.name

}

`, manufacturerName, manufacturerSlug, platformName, platformSlug)

}
