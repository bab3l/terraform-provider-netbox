// Package testutil provides utilities for acceptance testing of the Netbox provider.
package testutil

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var (
	// sharedClient is a singleton API client used across tests
	sharedClient     *netbox.APIClient
	sharedClientOnce sync.Once
	sharedClientErr  error
)

// GetSharedClient returns a shared Netbox API client for use in acceptance tests.
// The client is created once and reused across all tests.
func GetSharedClient() (*netbox.APIClient, error) {
	sharedClientOnce.Do(func() {
		serverURL := os.Getenv("NETBOX_SERVER_URL")
		apiToken := os.Getenv("NETBOX_API_TOKEN")

		if serverURL == "" || apiToken == "" {
			sharedClientErr = fmt.Errorf("NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables must be set")
			return
		}

		cfg := netbox.NewConfiguration()
		cfg.Servers = netbox.ServerConfigurations{
			{URL: serverURL},
		}
		cfg.DefaultHeader = map[string]string{
			"Authorization": "Token " + apiToken,
		}

		sharedClient = netbox.NewAPIClient(cfg)
	})

	return sharedClient, sharedClientErr
}

// RandomName generates a unique resource name with a given prefix.
// This helps avoid conflicts between test runs.
func RandomName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum))
}

// RandomSlug generates a unique slug with a given prefix.
// Slugs are lowercase with hyphens.
func RandomSlug(prefix string) string {
	return fmt.Sprintf("%s-%s", strings.ToLower(prefix), strings.ToLower(acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)))
}

// TestAccPreCheck validates the necessary test environment variables exist.
// It should be called at the beginning of each acceptance test.
func TestAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC must be set for acceptance tests")
	}

	if os.Getenv("NETBOX_SERVER_URL") == "" {
		t.Fatal("NETBOX_SERVER_URL must be set for acceptance tests")
	}

	if os.Getenv("NETBOX_API_TOKEN") == "" {
		t.Fatal("NETBOX_API_TOKEN must be set for acceptance tests")
	}
}

// CleanupResource is a helper to register cleanup functions that will run
// even if the test fails. Use this with t.Cleanup() to ensure resources
// are deleted via the API as a fallback.
type CleanupResource struct {
	client *netbox.APIClient
	t      *testing.T
}

// NewCleanupResource creates a new cleanup helper.
func NewCleanupResource(t *testing.T) *CleanupResource {
	t.Helper()
	client, err := GetSharedClient()
	if err != nil {
		t.Fatalf("Failed to get shared client for cleanup: %v", err)
	}
	return &CleanupResource{
		client: client,
		t:      t,
	}
}

// RegisterSiteGroupCleanup registers a cleanup function that will delete
// a site group by slug after the test completes.
func (c *CleanupResource) RegisterSiteGroupCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Find the site group by slug
		list, resp, err := c.client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			// Log but don't fail - resource might already be deleted
			c.t.Logf("Cleanup: failed to list site groups with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: site group with slug %s not found (already deleted)", slug)
			return
		}

		// Delete the site group
		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimSiteGroupsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete site group %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted site group %d (slug: %s)", id, slug)
		}
	})
}

// RegisterSiteCleanup registers a cleanup function that will delete
// a site by slug after the test completes.
func (c *CleanupResource) RegisterSiteCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimSitesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list sites with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: site with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimSitesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete site %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted site %d (slug: %s)", id, slug)
		}
	})
}

// RegisterTenantGroupCleanup registers a cleanup function that will delete
// a tenant group by slug after the test completes.
func (c *CleanupResource) RegisterTenantGroupCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list tenant groups with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: tenant group with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.TenancyAPI.TenancyTenantGroupsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete tenant group %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted tenant group %d (slug: %s)", id, slug)
		}
	})
}

// RegisterTenantCleanup registers a cleanup function that will delete
// a tenant by slug after the test completes.
func (c *CleanupResource) RegisterTenantCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list tenants with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: tenant with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.TenancyAPI.TenancyTenantsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete tenant %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted tenant %d (slug: %s)", id, slug)
		}
	})
}

// RegisterManufacturerCleanup registers a cleanup function that will delete
// a manufacturer by slug after the test completes.
func (c *CleanupResource) RegisterManufacturerCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list manufacturers with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: manufacturer with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimManufacturersDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete manufacturer %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted manufacturer %d (slug: %s)", id, slug)
		}
	})
}

// RegisterPlatformCleanup registers a cleanup function that will delete
// a platform by slug after the test completes.
func (c *CleanupResource) RegisterPlatformCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list platforms with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: platform with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimPlatformsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete platform %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted platform %d (slug: %s)", id, slug)
		}
	})
}

// RegisterRegionCleanup registers a cleanup function that will delete
// a region by slug after the test completes.
func (c *CleanupResource) RegisterRegionCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRegionsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list regions with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: region with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimRegionsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete region %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted region %d (slug: %s)", id, slug)
		}
	})
}

// RegisterLocationCleanup registers a cleanup function that will delete
// a location by slug after the test completes.
func (c *CleanupResource) RegisterLocationCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimLocationsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list locations with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: location with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimLocationsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete location %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted location %d (slug: %s)", id, slug)
		}
	})
}

// RegisterRackCleanup registers a cleanup function that will delete
// a rack by name after the test completes.
func (c *CleanupResource) RegisterRackCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRacksList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list racks with name %s: %v", name, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: rack with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimRacksDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete rack %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted rack %d (name: %s)", id, name)
		}
	})
}

// RegisterDeviceRoleCleanup registers a cleanup function that will delete
// a device role by slug after the test completes.
func (c *CleanupResource) RegisterDeviceRoleCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list device roles with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: device role with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimDeviceRolesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete device role %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted device role %d (slug: %s)", id, slug)
		}
	})
}

// RegisterRackRoleCleanup registers a cleanup function that will delete
// a rack role by slug after the test completes.
func (c *CleanupResource) RegisterRackRoleCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list rack roles with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: rack role with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimRackRolesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete rack role %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted rack role %d (slug: %s)", id, slug)
		}
	})
}

// RegisterDeviceTypeCleanup registers a cleanup function that will delete
// a device type by slug after the test completes.
func (c *CleanupResource) RegisterDeviceTypeCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list device types with slug %s: %v", slug, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: device type with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimDeviceTypesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete device type %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted device type %d (slug: %s)", id, slug)
		}
	})
}

// RegisterDeviceCleanup registers a cleanup function that will delete
// a device by name after the test completes.
func (c *CleanupResource) RegisterDeviceCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDevicesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list devices with name %s: %v", name, err)
			return
		}
		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: device with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimDevicesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete device %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted device %d (name: %s)", id, name)
		}
	})
}
