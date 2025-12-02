// Package testutil provides utilities for acceptance testing of the Netbox provider.
package testutil

import (
	"context"
	"fmt"
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
