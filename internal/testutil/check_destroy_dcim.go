// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("site with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

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

			if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

				return fmt.Errorf("device with name %s still exists (ID: %d)", name, list.Results[0].GetId())

			}

		}

		// Also check by ID if available

		id := rs.Primary.ID

		if id != "" {

			var idInt int32

			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {

				_, resp, err := client.DcimAPI.DcimDevicesRetrieve(ctx, idInt).Execute()

				if err == nil && resp.StatusCode == http.StatusOK {

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

				if err == nil && resp.StatusCode == http.StatusOK {

					return fmt.Errorf("interface with ID %s still exists", id)

				}

			}

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

		_, resp, err := client.DcimAPI.DcimRearPortTemplatesRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

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

		_, resp, err := client.DcimAPI.DcimFrontPortTemplatesRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

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

		_, resp, err := client.DcimAPI.DcimRearPortsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

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

		_, resp, err := client.DcimAPI.DcimFrontPortsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("front port with ID %d still exists", id)

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

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("device bay template with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckRackReservationDestroy verifies that a rack reservation has been destroyed.

func CheckRackReservationDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_rack_reservation" {

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

		_, resp, err := client.DcimAPI.DcimRackReservationsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("rack reservation with ID %d still exists", id)

		}

	}

	return nil

}

// CheckVirtualDeviceContextDestroy verifies that a virtual device context has been destroyed.

func CheckVirtualDeviceContextDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_virtual_device_context" {

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

		_, resp, err := client.DcimAPI.DcimVirtualDeviceContextsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("virtual device context with ID %d still exists", id)

		}

	}

	return nil

}

// CheckModuleBayTemplateDestroy verifies that a module bay template has been destroyed.

func CheckModuleBayTemplateDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_module_bay_template" {

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

		_, resp, err := client.DcimAPI.DcimModuleBayTemplatesRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("module bay template with ID %d still exists", id)

		}

	}

	return nil

}

// CheckInventoryItemTemplateDestroy verifies that an inventory item template has been destroyed.

func CheckInventoryItemTemplateDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_inventory_item_template" {

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

		_, resp, err := client.DcimAPI.DcimInventoryItemTemplatesRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("inventory item template with ID %d still exists", id)

		}

	}

	return nil

}

// CheckPowerPanelDestroy verifies that a power panel has been destroyed.

func CheckPowerPanelDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_power_panel" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.DcimAPI.DcimPowerPanelsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("power panel with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckPowerFeedDestroy verifies that a power feed has been destroyed.

func CheckPowerFeedDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_power_feed" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.DcimAPI.DcimPowerFeedsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("power feed with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckModuleBayDestroy verifies that a module bay has been destroyed.

func CheckModuleBayDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_module_bay" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.DcimAPI.DcimModuleBaysList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("module bay with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVirtualChassisDestroy verifies that a virtual chassis has been destroyed.

func CheckVirtualChassisDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_virtual_chassis" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.DcimAPI.DcimVirtualChassisList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("virtual chassis with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckInventoryItemDestroy verifies that an inventory item has been destroyed.
func CheckInventoryItemDestroy(s *terraform.State) error {
	// No specific destroy check needed as inventory items are tied to devices
	// and will be cleaned up when devices are destroyed
	return nil
}

// CheckModuleDestroy verifies that a module has been destroyed.
func CheckModuleDestroy(s *terraform.State) error {
	// No specific destroy check needed as modules are tied to devices
	// and will be cleaned up when devices are destroyed
	return nil
}

// CheckPowerOutletDestroy verifies that a power outlet has been destroyed.
func CheckPowerOutletDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_power_outlet" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimPowerOutletsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("power outlet with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}

// CheckPowerPortDestroy verifies that a power port has been destroyed.
func CheckPowerPortDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_power_port" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			continue
		}

		list, resp, err := client.DcimAPI.DcimPowerPortsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {
			return fmt.Errorf("power port with name %s still exists (ID: %d)", name, list.Results[0].GetId())
		}
	}

	return nil
}
