// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"fmt"
	"os"
	"strings"
	"sync"

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
func TestAccPreCheck(t interface {
	Fatal(args ...interface{})
	Skip(args ...interface{})
	Helper()
}) {
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
