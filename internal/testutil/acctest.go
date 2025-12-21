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

	// sharedClient is a singleton API client used across tests.

	sharedClient *netbox.APIClient

	sharedClientOnce sync.Once

	sharedClientErr error
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

// GenerateSlug generates a unique slug with a given prefix.

// This is an alias for RandomSlug for readability.

func GenerateSlug(prefix string) string {

	return RandomSlug(prefix)

}

// RandomVID generates a random VLAN ID between 2 and 4094.

// Range is limited to avoid reserved VLAN IDs.

func RandomVID() int32 {

	return int32(acctest.RandIntRange(2, 4094)) // #nosec G115 -- test value in safe range

}

// RandomIPv4Prefix generates a random private IPv4 prefix.

// Uses 10.x.x.0/24 format to avoid conflicts.

func RandomIPv4Prefix() string {

	// Use 10.x.x.0/24 format

	second := acctest.RandIntRange(0, 255)

	third := acctest.RandIntRange(0, 255)

	return fmt.Sprintf("10.%d.%d.0/24", second, third)

}

// RandomIPv6Prefix generates a random IPv6 prefix using ULA (Unique Local Address).

// Uses fd00:xxxx:xxxx::/48 format.

func RandomIPv6Prefix() string {

	// Use fd00:xxxx:xxxx::/48 format (ULA)

	segment1 := acctest.RandIntRange(0, 65535)

	segment2 := acctest.RandIntRange(0, 65535)

	return fmt.Sprintf("fd00:%04x:%04x::/48", segment1, segment2)

}

// RandomIPv4Address generates a random private IPv4 address with CIDR notation.

// Uses 10.x.x.x/32 format to avoid conflicts.

func RandomIPv4Address() string {

	second := acctest.RandIntRange(0, 255)

	third := acctest.RandIntRange(0, 255)

	fourth := acctest.RandIntRange(1, 254)

	return fmt.Sprintf("10.%d.%d.%d/32", second, third, fourth)

}

// RandomIPv6Address generates a random IPv6 address with CIDR notation using ULA.

// Uses fd00:xxxx:xxxx::x/128 format.

func RandomIPv6Address() string {

	segment1 := acctest.RandIntRange(0, 65535)

	segment2 := acctest.RandIntRange(0, 65535)

	host := acctest.RandIntRange(1, 65535)

	return fmt.Sprintf("fd00:%04x:%04x::%x/128", segment1, segment2, host)

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

	t *testing.T
}

// NewCleanupResource creates a new cleanup helper.

func NewCleanupResource(t *testing.T) *CleanupResource {

	t.Helper()

	if os.Getenv("TF_ACC") == "" {

		t.Skip("TF_ACC must be set for acceptance tests")

	}

	client, err := GetSharedClient()

	if err != nil {

		t.Fatalf("Failed to get shared client for cleanup: %v", err)

	}

	return &CleanupResource{

		client: client,

		t: t,
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

// RegisterInterfaceCleanup registers a cleanup function that will delete

// an interface by name and device after the test completes.

func (c *CleanupResource) RegisterInterfaceCleanup(name string, deviceName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimInterfacesList(ctx).Name([]string{name}).Device([]*string{&deviceName}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list interfaces with name %s on device %s: %v", name, deviceName, err)

			return

		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: interface with name %s on device %s not found (already deleted)", name, deviceName)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimInterfacesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete interface %d (name: %s, device: %s): %v", id, name, deviceName, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted interface %d (name: %s, device: %s)", id, name, deviceName)

		}

	})

}

// RegisterVRFCleanup registers a cleanup function that will delete

// a VRF by name after the test completes.

func (c *CleanupResource) RegisterVRFCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVrfsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VRFs with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: VRF with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVrfsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VRF %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VRF %d (name: %s)", id, name)

		}

	})

}

// RegisterVLANGroupCleanup registers a cleanup function that will delete

// a VLAN group by slug after the test completes.

func (c *CleanupResource) RegisterVLANGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VLAN groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: VLAN group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVlanGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VLAN group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VLAN group %d (slug: %s)", id, slug)

		}

	})

}

// RegisterVLANCleanup registers a cleanup function that will delete

// a VLAN by vid after the test completes.

func (c *CleanupResource) RegisterVLANCleanup(vid int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVlansList(ctx).Vid([]int32{vid}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VLANs with VID %d: %v", vid, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: VLAN with VID %d not found (already deleted)", vid)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVlansDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VLAN %d (VID: %d): %v", id, vid, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VLAN %d (VID: %d)", id, vid)

		}

	})

}

// RegisterPrefixCleanup registers a cleanup function that will delete

// a prefix by CIDR after the test completes.

func (c *CleanupResource) RegisterPrefixCleanup(prefix string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamPrefixesList(ctx).Prefix([]string{prefix}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list prefixes with CIDR %s: %v", prefix, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: prefix with CIDR %s not found (already deleted)", prefix)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamPrefixesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete prefix %d (CIDR: %s): %v", id, prefix, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted prefix %d (CIDR: %s)", id, prefix)

		}

	})

}

// RegisterIPAddressCleanup registers a cleanup function that will delete

// an IP address by address after the test completes.

func (c *CleanupResource) RegisterIPAddressCleanup(address string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamIpAddressesList(ctx).Address([]string{address}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IP addresses with address %s: %v", address, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IP address %s not found (already deleted)", address)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamIpAddressesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IP address %d (address: %s): %v", id, address, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IP address %d (address: %s)", id, address)

		}

	})

}

// RegisterClusterTypeCleanup registers a cleanup function that will delete

// a cluster type by slug after the test completes.

func (c *CleanupResource) RegisterClusterTypeCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list cluster types with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: cluster type with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationClusterTypesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete cluster type %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted cluster type %d (slug: %s)", id, slug)

		}

	})

}

// RegisterClusterCleanup registers a cleanup function that will delete

// a cluster by name after the test completes.

func (c *CleanupResource) RegisterClusterCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list clusters with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: cluster with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationClustersDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete cluster %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted cluster %d (name: %s)", id, name)

		}

	})

}

// RegisterVirtualMachineCleanup registers a cleanup function that will delete

// a virtual machine by name after the test completes.

func (c *CleanupResource) RegisterVirtualMachineCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list virtual machines with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: virtual machine with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual machine %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual machine %d (name: %s)", id, name)

		}

	})

}

// RegisterVMInterfaceCleanup registers a cleanup function that will delete

// a VM interface by name and virtual machine after the test completes.

func (c *CleanupResource) RegisterVMInterfaceCleanup(name string, vmName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationInterfacesList(ctx).Name([]string{name}).VirtualMachine([]string{vmName}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VM interfaces with name %s on VM %s: %v", name, vmName, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: VM interface with name %s on VM %s not found (already deleted)", name, vmName)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationInterfacesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VM interface %d (name: %s, VM: %s): %v", id, name, vmName, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VM interface %d (name: %s, VM: %s)", id, name, vmName)

		}

	})

}

// RegisterProviderCleanup registers a cleanup function that will delete

// a circuit provider by slug after the test completes.

func (c *CleanupResource) RegisterProviderCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list providers with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: provider with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsProvidersDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete provider %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted provider %d (slug: %s)", id, slug)

		}

	})

}

// RegisterCircuitTypeCleanup registers a cleanup function that will delete

// a circuit type by slug after the test completes.

func (c *CleanupResource) RegisterCircuitTypeCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit types with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: circuit type with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitTypesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit type %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit type %d (slug: %s)", id, slug)

		}

	})

}

// RegisterCircuitCleanup registers a cleanup function that will delete

// a circuit by CID after the test completes.

func (c *CleanupResource) RegisterCircuitCleanup(cid string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{cid}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuits with CID %s: %v", cid, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: circuit with CID %s not found (already deleted)", cid)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit %d (CID: %s): %v", id, cid, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit %d (CID: %s)", id, cid)

		}

	})

}

// RegisterRouteTargetCleanup registers a cleanup function that will delete

// a route target by name after the test completes.

func (c *CleanupResource) RegisterRouteTargetCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRouteTargetsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list route targets with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: route target with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamRouteTargetsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete route target %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted route target %d (name: %s)", id, name)

		}

	})

}

// RegisterVirtualDiskCleanup registers a cleanup function that will delete

// a virtual disk by name after the test completes.

func (c *CleanupResource) RegisterVirtualDiskCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationVirtualDisksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list virtual disks with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: virtual disk with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationVirtualDisksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual disk %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual disk %d (name: %s)", id, name)

		}

	})

}

// RegisterASNRangeCleanup registers a cleanup function that will delete

// an ASN range by name after the test completes.

func (c *CleanupResource) RegisterASNRangeCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamAsnRangesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list ASN ranges with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: ASN range with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamAsnRangesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete ASN range %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted ASN range %d (name: %s)", id, name)

		}

	})

}

// RegisterDeviceBayTemplateCleanup registers a cleanup function that will delete

// a device bay template by name after the test completes.

func (c *CleanupResource) RegisterDeviceBayTemplateCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimDeviceBayTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list device bay templates with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: device bay template with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimDeviceBayTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete device bay template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted device bay template %d (name: %s)", id, name)

		}

	})

}

// RegisterRIRCleanup registers a cleanup function that will delete

// an RIR by slug after the test completes.

func (c *CleanupResource) RegisterRIRCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRirsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list RIRs with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: RIR with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamRirsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete RIR %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted RIR %d (slug: %s)", id, slug)

		}

	})

}

// RegisterIKEProposalCleanup registers a cleanup function that will delete

// an IKE proposal by name after the test completes.

func (c *CleanupResource) RegisterIKEProposalCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIkeProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IKE proposals with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IKE proposal with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIkeProposalsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IKE proposal %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IKE proposal %d (name: %s)", id, name)

		}

	})

}

// RegisterIKEPolicyCleanup registers a cleanup function that will delete

// an IKE policy by name after the test completes.

func (c *CleanupResource) RegisterIKEPolicyCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIkePoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IKE policies with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IKE policy with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIkePoliciesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IKE policy %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IKE policy %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecProposalCleanup registers a cleanup function that will delete

// an IPSec proposal by name after the test completes.

func (c *CleanupResource) RegisterIPSecProposalCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec proposals with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec proposal with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecProposalsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec proposal %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec proposal %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecPolicyCleanup registers a cleanup function that will delete

// an IPSec policy by name after the test completes.

func (c *CleanupResource) RegisterIPSecPolicyCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecPoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec policies with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec policy with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecPoliciesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec policy %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec policy %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecProfileCleanup registers a cleanup function that will delete

// an IPSec profile by name after the test completes.

func (c *CleanupResource) RegisterIPSecProfileCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecProfilesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec profiles with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec profile with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecProfilesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec profile %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec profile %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelGroupCleanup registers a cleanup function that will delete

// a tunnel group by name after the test completes.

func (c *CleanupResource) RegisterTunnelGroupCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnTunnelGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnel groups with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: tunnel group with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnTunnelGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tunnel group %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tunnel group %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelCleanup registers a cleanup function that will delete

// a tunnel by name after the test completes.

func (c *CleanupResource) RegisterTunnelCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnels with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: tunnel with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnTunnelsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tunnel %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tunnel %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelTerminationCleanup registers a cleanup function that will delete

// a tunnel termination by ID after the test completes.

func (c *CleanupResource) RegisterTunnelTerminationCleanup(tunnelName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First find the tunnel by name

		tunnelList, resp, err := c.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{tunnelName}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnels with name %s: %v", tunnelName, err)

			return

		}

		if resp.StatusCode != 200 || tunnelList.Count == 0 {

			c.t.Logf("Cleanup: tunnel with name %s not found (already deleted)", tunnelName)

			return

		}

		tunnelID := tunnelList.Results[0].GetId()

		// List terminations for this tunnel

		termList, _, err := c.client.VpnAPI.VpnTunnelTerminationsList(ctx).TunnelId([]int32{tunnelID}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnel terminations for tunnel %d: %v", tunnelID, err)

			return

		}

		for _, term := range termList.Results {

			_, err := c.client.VpnAPI.VpnTunnelTerminationsDestroy(ctx, term.GetId()).Execute()

			if err != nil {

				c.t.Logf("Cleanup: failed to delete tunnel termination %d: %v", term.GetId(), err)

			} else {

				c.t.Logf("Cleanup: successfully deleted tunnel termination %d", term.GetId())

			}

		}

	})

}

// RegisterCircuitGroupCleanup registers a cleanup function that will delete

// a circuit group by name after the test completes.

func (c *CleanupResource) RegisterCircuitGroupCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit groups with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: circuit group with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit group %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit group %d (name: %s)", id, name)

		}

	})

}

// RegisterCircuitGroupAssignmentCleanup registers a cleanup function that will delete

// a circuit group assignment by ID after the test completes.

func (c *CleanupResource) RegisterCircuitGroupAssignmentCleanup(groupName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First find the circuit group by name

		groupList, resp, err := c.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{groupName}).Execute()

		if err != nil || resp.StatusCode != 200 || groupList.Count == 0 {

			c.t.Logf("Cleanup: circuit group with name %s not found, cannot cleanup assignments", groupName)

			return

		}

		groupID := groupList.Results[0].GetId()

		// Find all assignments for this group

		assignmentList, _, err := c.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsList(ctx).GroupId([]int32{groupID}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit group assignments for group %d: %v", groupID, err)

			return

		}

		for _, assignment := range assignmentList.Results {

			_, err := c.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsDestroy(ctx, assignment.GetId()).Execute()

			if err != nil {

				c.t.Logf("Cleanup: failed to delete circuit group assignment %d: %v", assignment.GetId(), err)

			} else {

				c.t.Logf("Cleanup: successfully deleted circuit group assignment %d", assignment.GetId())

			}

		}

	})

}

// RegisterRearPortTemplateCleanup registers a cleanup function that will delete

// a rear port template by name and device type after the test completes.

func (c *CleanupResource) RegisterRearPortTemplateCleanup(name string, deviceTypeID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		dtID := deviceTypeID

		list, resp, err := c.client.DcimAPI.DcimRearPortTemplatesList(ctx).Name([]string{name}).DeviceTypeId([]*int32{&dtID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: rear port template with name %s and device type %d not found: %v", name, deviceTypeID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRearPortTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rear port template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rear port template %d (name: %s)", id, name)

		}

	})

}

// RegisterFrontPortTemplateCleanup registers a cleanup function that will delete

// a front port template by name and device type after the test completes.

func (c *CleanupResource) RegisterFrontPortTemplateCleanup(name string, deviceTypeID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		dtID := deviceTypeID

		list, resp, err := c.client.DcimAPI.DcimFrontPortTemplatesList(ctx).Name([]string{name}).DeviceTypeId([]*int32{&dtID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: front port template with name %s and device type %d not found: %v", name, deviceTypeID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimFrontPortTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete front port template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted front port template %d (name: %s)", id, name)

		}

	})

}

// RegisterRearPortCleanup registers a cleanup function that will delete

// a rear port by name and device after the test completes.

func (c *CleanupResource) RegisterRearPortCleanup(name string, deviceID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimRearPortsList(ctx).Name([]string{name}).DeviceId([]int32{deviceID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: rear port with name %s and device %d not found: %v", name, deviceID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimRearPortsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rear port %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rear port %d (name: %s)", id, name)

		}

	})

}

// RegisterFrontPortCleanup registers a cleanup function that will delete

// a front port by name and device after the test completes.

func (c *CleanupResource) RegisterFrontPortCleanup(name string, deviceID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimFrontPortsList(ctx).Name([]string{name}).DeviceId([]int32{deviceID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: front port with name %s and device %d not found: %v", name, deviceID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.DcimAPI.DcimFrontPortsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete front port %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted front port %d (name: %s)", id, name)

		}

	})

}

// RegisterFHRPGroupCleanup registers a cleanup function that will delete

// an FHRP group by protocol and group_id after the test completes.

func (c *CleanupResource) RegisterFHRPGroupCleanup(protocol string, groupID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamFhrpGroupsList(ctx).Protocol([]string{protocol}).GroupId([]int32{groupID}).Execute()

		if err != nil || resp.StatusCode != 200 || list.Count == 0 {

			c.t.Logf("Cleanup: FHRP group with protocol %s and group_id %d not found: %v", protocol, groupID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamFhrpGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete FHRP group %d (protocol: %s, group_id: %d): %v", id, protocol, groupID, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted FHRP group %d (protocol: %s, group_id: %d)", id, protocol, groupID)

		}

	})

}

// RegisterJournalEntryCleanup registers a cleanup function that will delete

// a journal entry by ID after the test completes.

func (c *CleanupResource) RegisterJournalEntryCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasJournalEntriesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete journal entry %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted journal entry %d", id)

		}

	})

}

// RegisterCustomFieldChoiceSetCleanup registers a cleanup function that will delete

// a custom field choice set by ID after the test completes.

func (c *CleanupResource) RegisterCustomFieldChoiceSetCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom field choice set %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom field choice set %d", id)

		}

	})

}

// RegisterCustomLinkCleanup registers a cleanup function that will delete

// a custom link by ID after the test completes.

func (c *CleanupResource) RegisterCustomLinkCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasCustomLinksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom link %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom link %d", id)

		}

	})

}

// RegisterWirelessLinkCleanup registers a cleanup function that will delete

// a wireless link by ID after the test completes.

func (c *CleanupResource) RegisterWirelessLinkCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.WirelessAPI.WirelessWirelessLinksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete wireless link %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted wireless link %d", id)

		}

	})

}

// RegisterRackReservationCleanup registers a cleanup function that will delete

// a rack reservation by ID after the test completes.

func (c *CleanupResource) RegisterRackReservationCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimRackReservationsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete rack reservation %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted rack reservation %d", id)

		}

	})

}

// RegisterVirtualDeviceContextCleanup registers a cleanup function that will delete

// a virtual device context by ID after the test completes.

func (c *CleanupResource) RegisterVirtualDeviceContextCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimVirtualDeviceContextsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual device context %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual device context %d", id)

		}

	})

}

// RegisterModuleBayTemplateCleanup registers a cleanup function that will delete

// a module bay template by ID after the test completes.

func (c *CleanupResource) RegisterModuleBayTemplateCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimModuleBayTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete module bay template %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted module bay template %d", id)

		}

	})

}

// Note: RegisterCableTerminationCleanup removed - netbox_cable_termination resource deprecated

// Use netbox_cable with embedded terminations instead

// RegisterInventoryItemTemplateCleanup registers a cleanup function that will delete

// an inventory item template by ID after the test completes.

func (c *CleanupResource) RegisterInventoryItemTemplateCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.DcimAPI.DcimInventoryItemTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete inventory item template %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted inventory item template %d", id)

		}

	})

}

// RegisterContactGroupCleanup registers a cleanup function that will delete
// a contact group by slug after the test completes.
func (c *CleanupResource) RegisterContactGroupCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyContactGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list contact groups with slug %s: %v", slug, err)
			return
		}

		if resp.StatusCode != 200 || list == nil || len(list.Results) == 0 {
			c.t.Logf("Cleanup: contact group with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.TenancyAPI.TenancyContactGroupsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete contact group %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted contact group %d (slug: %s)", id, slug)
		}
	})
}

func (c *CleanupResource) RegisterClusterGroupCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list cluster groups with slug %s: %v", slug, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: cluster group with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.VirtualizationAPI.VirtualizationClusterGroupsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete cluster group %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted cluster group %d (slug: %s)", id, slug)
		}
	})
}

// RegisterTagCleanup registers a cleanup function that will delete
// a tag by slug after the test completes.
func (c *CleanupResource) RegisterTagCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasTagsList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list tags with slug %s: %v", slug, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: tag with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.ExtrasAPI.ExtrasTagsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete tag %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted tag %d (slug: %s)", id, slug)
		}
	})
}

// RegisterRoleCleanup registers a cleanup function that will delete
// a role by slug after the test completes.
func (c *CleanupResource) RegisterRoleCleanup(slug string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRolesList(ctx).Slug([]string{slug}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list roles with slug %s: %v", slug, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: role with slug %s not found (already deleted)", slug)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.IpamAPI.IpamRolesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete role %d (slug: %s): %v", id, slug, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted role %d (slug: %s)", id, slug)
		}
	})
}

// RegisterWebhookCleanup registers a cleanup function that will delete
// a webhook by name after the test completes.
func (c *CleanupResource) RegisterWebhookCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasWebhooksList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list webhooks with name %s: %v", name, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: webhook with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.ExtrasAPI.ExtrasWebhooksDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete webhook %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted webhook %d (name: %s)", id, name)
		}
	})
}

// RegisterPowerPanelCleanup registers a cleanup function that will delete
// a power panel by name after the test completes.
func (c *CleanupResource) RegisterPowerPanelCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPowerPanelsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list power panels with name %s: %v", name, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: power panel with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimPowerPanelsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete power panel %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted power panel %d (name: %s)", id, name)
		}
	})
}

// RegisterPowerFeedCleanup registers a cleanup function that will delete
// a power feed by name after the test completes.
func (c *CleanupResource) RegisterPowerFeedCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimPowerFeedsList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list power feeds with name %s: %v", name, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: power feed with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimPowerFeedsDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete power feed %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted power feed %d (name: %s)", id, name)
		}
	})
}

// RegisterModuleBayCleanup registers a cleanup function that will delete
// a module bay by name and device ID after the test completes.
func (c *CleanupResource) RegisterModuleBayCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.DcimAPI.DcimModuleBaysList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list module bays with name %s: %v", name, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: module bay with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.DcimAPI.DcimModuleBaysDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete module bay %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted module bay %d (name: %s)", id, name)
		}
	})
}

// RegisterServiceCleanup registers a cleanup function that will delete
// a service by name and device ID after the test completes.
func (c *CleanupResource) RegisterServiceCleanup(name string) {
	c.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamServicesList(ctx).Name([]string{name}).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to list services with name %s: %v", name, err)
			return
		}

		if resp.StatusCode != 200 || list.Count == 0 {
			c.t.Logf("Cleanup: service with name %s not found (already deleted)", name)
			return
		}

		id := list.Results[0].GetId()
		_, err = c.client.IpamAPI.IpamServicesDestroy(ctx, id).Execute()
		if err != nil {
			c.t.Logf("Cleanup: failed to delete service %d (name: %s): %v", id, name, err)
		} else {
			c.t.Logf("Cleanup: successfully deleted service %d (name: %s)", id, name)
		}
	})
}
