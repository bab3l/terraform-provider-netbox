package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevicePrimaryIPResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("10.0.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", "", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_device_primary_ip.test", "primary_ip6"),
					resource.TestCheckNoResourceAttr("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
		},
	})
}

func TestAccDevicePrimaryIPResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("10.1.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	oob := fmt.Sprintf("10.2.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4)
	cleanup.RegisterIPAddressCleanup(oob)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", oob, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
		},
	})
}

func TestAccDevicePrimaryIPResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4a := fmt.Sprintf("10.3.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	ip4b := fmt.Sprintf("10.4.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	oob := fmt.Sprintf("10.5.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4a)
	cleanup.RegisterIPAddressCleanup(ip4b)
	cleanup.RegisterIPAddressCleanup(oob)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4a, "", oob, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4b, "", oob, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
		},
	})
}

func TestAccDevicePrimaryIPResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("10.6.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", "", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_device_primary_ip.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "primary_ip4", "primary_ip6", "oob_ip"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_device_primary_ip.test", "device"),
					testutil.ReferenceFieldCheck("netbox_device_primary_ip.test", "primary_ip4"),
					testutil.ReferenceFieldCheck("netbox_device_primary_ip.test", "primary_ip6"),
					testutil.ReferenceFieldCheck("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
			{
				Config:   testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", "", false, false),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDevicePrimaryIPResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("10.7.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4)

	testutil.RunExternalDeletionTest(t, testutil.ExternalDeletionTestConfig{
		ResourceName: "netbox_device_primary_ip",
		Config: func() string {
			return testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", "", false, false)
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			client, err := testutil.GetSharedClient()
			if err != nil {
				return err
			}
			deviceID64, err := strconv.ParseInt(id, 10, 32)
			if err != nil {
				return err
			}
			patch := netbox.NewPatchedWritableDeviceWithConfigContextRequest()
			patch.SetPrimaryIp4Nil()
			patch.SetPrimaryIp6Nil()
			patch.SetOobIpNil()
			_, _, err = client.DcimAPI.DcimDevicesPartialUpdate(ctx, int32(deviceID64)).
				PatchedWritableDeviceWithConfigContextRequest(*patch).
				Execute()
			return err
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
			testutil.CheckIPAddressDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
	})
}

func TestAccDevicePrimaryIPResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	deviceName := testutil.RandomName("tf-test-device")
	interfaceName := testutil.RandomName("eth")
	ip4 := fmt.Sprintf("10.8.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))
	oob := fmt.Sprintf("10.9.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterIPAddressCleanup(ip4)
	cleanup.RegisterIPAddressCleanup(oob)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", oob, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", oob, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckNoResourceAttr("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
			{
				Config: testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, "", oob, false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "primary_ip4"),
					resource.TestCheckResourceAttrSet("netbox_device_primary_ip.test", "oob_ip"),
				),
			},
		},
	})
}

func testAccDevicePrimaryIPResourceConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, ip6, oob string, setPrimaryIP6, setOob bool) string {
	primaryIP6Resource := ""
	primaryIP6Attr := ""
	if ip6 != "" {
		primaryIP6Resource = fmt.Sprintf(`
resource "netbox_ip_address" "test_v6" {
  address = %q
  status  = "active"
}
`, ip6)
		if setPrimaryIP6 {
			primaryIP6Attr = "\n  primary_ip6 = netbox_ip_address.test_v6.id"
		}
	}

	oobResource := ""
	oobAttr := ""
	if oob != "" {
		oobResource = fmt.Sprintf(`
resource "netbox_ip_address" "test_oob" {
  address              = %q
  status               = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.test.id
}
`, oob)
		if setOob {
			oobAttr = "\n  oob_ip = netbox_ip_address.test_oob.id"
		}
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

resource "netbox_interface" "test" {
	device = netbox_device.test.id
  name   = %q
  type   = "1000base-t"
}

resource "netbox_ip_address" "test_v4" {
  address              = %q
  status               = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.test.id
}
%s%s
resource "netbox_device_primary_ip" "test" {
	device      = netbox_device.test.id
  primary_ip4 = netbox_ip_address.test_v4.id%s%s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, ip4, primaryIP6Resource, oobResource, primaryIP6Attr, oobAttr)
}
