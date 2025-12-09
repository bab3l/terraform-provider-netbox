// Package testutil provides utilities for acceptance testing of the Netbox provider.
package testutil

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// CheckSiteGroupDestroy verifies that a site group has been destroyed.
// Use this as the CheckDestroy function in resource.TestCase.
func CheckSiteGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_site_group" {
			continue
		}

		// Try to find the resource by slug
		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			// API error could mean resource doesn't exist, which is expected
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("site group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckSiteDestroy verifies that a site has been destroyed.
func CheckSiteDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_site" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimSitesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("site with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckTenantGroupDestroy verifies that a tenant group has been destroyed.
func CheckTenantGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_tenant_group" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("tenant group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckTenantDestroy verifies that a tenant has been destroyed.
func CheckTenantDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_tenant" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("tenant with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckManufacturerDestroy verifies that a manufacturer has been destroyed.
func CheckManufacturerDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_manufacturer" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("manufacturer with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckPlatformDestroy verifies that a platform has been destroyed.
func CheckPlatformDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_platform" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("platform with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckRegionDestroy verifies that a region has been destroyed.
func CheckRegionDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_region" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimRegionsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("region with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckLocationDestroy verifies that a location has been destroyed.
func CheckLocationDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_location" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimLocationsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("location with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckRackDestroy verifies that a rack has been destroyed.
func CheckRackDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_rack" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimRacksList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("rack with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckDeviceRoleDestroy verifies that a device role has been destroyed.
func CheckDeviceRoleDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_device_role" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("device role with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckRackRoleDestroy verifies that a rack role has been destroyed.
func CheckRackRoleDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_rack_role" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("rack role with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckDeviceTypeDestroy verifies that a device type has been destroyed.
func CheckDeviceTypeDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_device_type" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("device type with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckDeviceDestroy verifies that a device has been destroyed.
func CheckDeviceDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_device" {
			continue
		}

		// Device may not have a unique slug, so check by name or ID
		name := rs.Primary.Attributes["name"]
		if name != "" {
			list, resp, err := client.DcimAPI.DcimDevicesList(ctx).Name([]string{name}).Execute()
			if err != nil {
				continue
			}

			if resp.StatusCode == 200 && list != nil && len(list.Results) > 0 {
				return fmt.Errorf("device with name %s still exists (ID: %d)", name, list.Results[0].GetId())
			}
		}

		// Also check by ID if available
		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.DcimAPI.DcimDevicesRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("device with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// CheckInterfaceDestroy verifies that an interface has been destroyed.
func CheckInterfaceDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_interface" {
			continue
		}

		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.DcimAPI.DcimInterfacesRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("interface with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// ComposeCheckDestroy combines multiple CheckDestroy functions.
// Useful when a test creates multiple resource types.
func ComposeCheckDestroy(checks ...resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, check := range checks {
			if err := check(s); err != nil {
				return err
			}
		}
		return nil
	}
}

// CheckVRFDestroy verifies that a VRF has been destroyed.
func CheckVRFDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_vrf" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.IpamAPI.IpamVrfsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("VRF with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckVLANGroupDestroy verifies that a VLAN group has been destroyed.
func CheckVLANGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_vlan_group" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("VLAN group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckVLANDestroy verifies that a VLAN has been destroyed.
func CheckVLANDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_vlan" {
			continue
		}

		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.IpamAPI.IpamVlansRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("VLAN with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// CheckPrefixDestroy verifies that a prefix has been destroyed.
func CheckPrefixDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_prefix" {
			continue
		}

		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.IpamAPI.IpamPrefixesRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("prefix with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// CheckIPAddressDestroy verifies that an IP address has been destroyed.
func CheckIPAddressDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ip_address" {
			continue
		}

		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.IpamAPI.IpamIpAddressesRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("IP address with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// CheckClusterTypeDestroy verifies that a cluster type has been destroyed.
func CheckClusterTypeDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_cluster_type" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("cluster type with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckClusterDestroy verifies that a cluster has been destroyed.
func CheckClusterDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("cluster with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckVirtualMachineDestroy verifies that a virtual machine has been destroyed.
func CheckVirtualMachineDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_virtual_machine" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("virtual machine with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckVMInterfaceDestroy verifies that a VM interface has been destroyed.
func CheckVMInterfaceDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_vm_interface" {
			continue
		}

		id := rs.Primary.ID
		if id != "" {
			var idInt int32
			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {
				_, resp, err := client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, idInt).Execute()
				if err == nil && resp.StatusCode == 200 {
					return fmt.Errorf("VM interface with ID %s still exists", id)
				}
			}
		}
	}

	return nil
}

// CheckProviderDestroy verifies that a circuit provider has been destroyed.
func CheckProviderDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_provider" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("provider with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckCircuitTypeDestroy verifies that a circuit type has been destroyed.
func CheckCircuitTypeDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_circuit_type" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("circuit type with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckCircuitDestroy verifies that a circuit has been destroyed.
func CheckCircuitDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_circuit" {
			continue
		}

		cid := rs.Primary.Attributes["cid"]
		if cid == "" {
			continue
		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{cid}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("circuit with CID %s still exists (ID: %d)", cid, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckRouteTargetDestroy verifies that a route target has been destroyed.
func CheckRouteTargetDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_route_target" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.IpamAPI.IpamRouteTargetsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("route target with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckVirtualDiskDestroy verifies that a virtual disk has been destroyed.
func CheckVirtualDiskDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_virtual_disk" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VirtualizationAPI.VirtualizationVirtualDisksList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("virtual disk with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckASNRangeDestroy verifies that an ASN range has been destroyed.
func CheckASNRangeDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_asn_range" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.IpamAPI.IpamAsnRangesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("ASN range with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckDeviceBayTemplateDestroy verifies that a device bay template has been destroyed.
func CheckDeviceBayTemplateDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_device_bay_template" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimDeviceBayTemplatesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("device bay template with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckRIRDestroy verifies that an RIR has been destroyed.
func CheckRIRDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_rir" {
			continue
		}

		slug := rs.Primary.Attributes["slug"]
		if slug == "" {
			continue
		}

		list, resp, err := client.IpamAPI.IpamRirsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("RIR with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckIKEProposalDestroy verifies that an IKE proposal has been destroyed.
func CheckIKEProposalDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ike_proposal" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnIkeProposalsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("IKE proposal with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckIKEPolicyDestroy verifies that an IKE policy has been destroyed.
func CheckIKEPolicyDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ike_policy" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnIkePoliciesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("IKE policy with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckIPSecProposalDestroy verifies that an IPSec proposal has been destroyed.
func CheckIPSecProposalDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ipsec_proposal" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnIpsecProposalsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("IPSec proposal with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckIPSecPolicyDestroy verifies that an IPSec policy has been destroyed.
func CheckIPSecPolicyDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ipsec_policy" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnIpsecPoliciesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("IPSec policy with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckIPSecProfileDestroy verifies that an IPSec profile has been destroyed.
func CheckIPSecProfileDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_ipsec_profile" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnIpsecProfilesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("IPSec profile with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckTunnelGroupDestroy verifies that a tunnel group has been destroyed.
func CheckTunnelGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_tunnel_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnTunnelGroupsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("tunnel group with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckTunnelDestroy verifies that a tunnel has been destroyed.
func CheckTunnelDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_tunnel" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.VpnAPI.VpnTunnelsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("tunnel with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckTunnelTerminationDestroy verifies that a tunnel termination has been destroyed.
func CheckTunnelTerminationDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_tunnel_termination" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.VpnAPI.VpnTunnelTerminationsRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("tunnel termination with ID %d still exists", id)
		}
	}

	return nil
}

// CheckCircuitGroupDestroy verifies that a circuit group has been destroyed.
func CheckCircuitGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_circuit_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 && list.Count > 0 {
			return fmt.Errorf("circuit group with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckCircuitGroupAssignmentDestroy verifies that a circuit group assignment has been destroyed.
func CheckCircuitGroupAssignmentDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_circuit_group_assignment" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.CircuitsAPI.CircuitsCircuitGroupAssignmentsRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("circuit group assignment with ID %d still exists", id)
		}
	}

	return nil
}

// CheckRearPortTemplateDestroy verifies that a rear port template has been destroyed.
func CheckRearPortTemplateDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_rear_port_template" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.DcimAPI.DcimRearPortTemplatesRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("rear port template with ID %d still exists", id)
		}
	}

	return nil
}

// CheckFrontPortTemplateDestroy verifies that a front port template has been destroyed.
func CheckFrontPortTemplateDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_front_port_template" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.DcimAPI.DcimFrontPortTemplatesRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("front port template with ID %d still exists", id)
		}
	}

	return nil
}

// CheckRearPortDestroy verifies that a rear port has been destroyed.
func CheckRearPortDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_rear_port" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.DcimAPI.DcimRearPortsRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("rear port with ID %d still exists", id)
		}
	}

	return nil
}

// CheckFrontPortDestroy verifies that a front port has been destroyed.
func CheckFrontPortDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_front_port" {
			continue
		}

		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			continue
		}

		_, resp, err := client.DcimAPI.DcimFrontPortsRetrieve(ctx, int32(id)).Execute()
		if err == nil && resp.StatusCode == 200 {
			return fmt.Errorf("front port with ID %d still exists", id)
		}
	}

	return nil
}
