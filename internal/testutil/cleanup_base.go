// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"os"
	"testing"

	"github.com/bab3l/go-netbox"
)

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
