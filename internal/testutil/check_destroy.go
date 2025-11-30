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
