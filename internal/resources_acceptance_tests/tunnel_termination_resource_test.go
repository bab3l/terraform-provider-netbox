package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelTerminationResource_basic(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-for-term")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_full(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-full")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	ipAddress := testutil.RandomIPv4Address()
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_full(tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "hub"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "outside_ip"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config: testAccTunnelTerminationResourceConfig_fullUpdate(tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "custom_fields.0.value", "updated_value"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_IDPreservation(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tnl-term-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_update(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-for-term-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_withRole(tunnelName, "peer"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "peer"),
				),
			},
			{
				Config: testAccTunnelTerminationResourceConfig_withRole(tunnelName, "hub"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "hub"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_import(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-for-term-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
			},
			{
				ResourceName:            "netbox_tunnel_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tunnel"},
			},
		},
	})
}

func TestAccTunnelTerminationResource_externalDeletion(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-for-term-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// Find tunnel termination by filtering for tunnel name
					items, _, err := client.VpnAPI.VpnTunnelTerminationsList(context.Background()).Tunnel([]string{tunnelName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tunnel termination for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnTunnelTerminationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tunnel termination: %v", err)
					}
					t.Logf("Successfully externally deleted tunnel termination with ID: %d", itemID)
				},
				Config: testAccTunnelTerminationResourceConfig_basic(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
				),
			},
		},
	})
}

func testAccTunnelTerminationResourceConfig_full(tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tag" "tag1" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[6]q
  object_types = ["vpn.tunneltermination"]
  type         = "text"
}

resource "netbox_ip_address" "outside" {
  address = %[7]q
  status  = "active"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "hub"
  outside_ip       = netbox_ip_address.outside.id

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}
`, tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, cfName, ipAddress)
}

func testAccTunnelTerminationResourceConfig_fullUpdate(tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, ipAddress, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tag" "tag1" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_custom_field" "test_field" {
  name         = %[6]q
  object_types = ["vpn.tunneltermination"]
  type         = "text"
}

resource "netbox_ip_address" "outside" {
  address = %[7]q
  status  = "active"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
  outside_ip       = netbox_ip_address.outside.id

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}
`, tunnelName, tagName1, tagSlug1, tagName2, tagSlug2, cfName, ipAddress)
}

func TestAccConsistency_TunnelTermination_LiteralNames(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tunnel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationConsistencyLiteralNamesConfig(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "id"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTunnelTerminationConsistencyLiteralNamesConfig(tunnelName),
			},
		},
	})
}

func testAccTunnelTerminationResourceConfig_basic(tunnelName string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}
`
}

func testAccTunnelTerminationConsistencyLiteralNamesConfig(tunnelName string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}
`, tunnelName)
}

func testAccTunnelTerminationResourceConfig_withRole(tunnelName, role string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "` + role + `"
}
`
}

func TestAccTunnelTerminationResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	tunnelName := testutil.RandomName("tf-test-tunnel-opt")
	ipAddress := testutil.RandomIPv4Address()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "ipsec-tunnel"
}

resource "netbox_ip_address" "test" {
  address = %[2]q
  status  = "active"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "hub"
  termination_id   = 123
  outside_ip       = netbox_ip_address.test.id
}
`, tunnelName, ipAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "hub"),
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "termination_id", "123"),
					resource.TestCheckResourceAttrSet("netbox_tunnel_termination.test", "outside_ip"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "ipsec-tunnel"
}

resource "netbox_ip_address" "test" {
  address = %[2]q
  status  = "active"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "hub"
}
`, tunnelName, ipAddress),
				Check: resource.ComposeAggregateTestCheckFunc(
					// role must remain as it's required by the API
					resource.TestCheckResourceAttr("netbox_tunnel_termination.test", "role", "hub"),
					// These fields should be removed
					resource.TestCheckNoResourceAttr("netbox_tunnel_termination.test", "termination_id"),
					resource.TestCheckNoResourceAttr("netbox_tunnel_termination.test", "outside_ip"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_tunnel_termination",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_tunnel": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "test-device-type"
  slug         = "test-device-type"
  manufacturer = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_tunnel_termination" "test" {
  # tunnel missing
  termination_type = "dcim.device"
  termination_id   = netbox_device.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_termination_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tunnel" "test" {
  name          = "test-tunnel"
  encapsulation = "ipsec-transport"
}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "test-device-type"
  slug         = "test-device-type"
  manufacturer = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_tunnel_termination" "test" {
  tunnel = netbox_tunnel.test.id
  # termination_type missing
  termination_id = netbox_device.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
