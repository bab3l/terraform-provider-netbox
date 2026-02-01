// Package utils provides utility functions for working with Netbox provider data structures.

package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// HandleNotFound executes onNotFound and returns true if the response indicates a 404.
func HandleNotFound(httpResp *http.Response, onNotFound func()) bool {
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		if onNotFound != nil {
			onNotFound()
		}
		return true
	}

	return false
}

// ValidateStatusCode checks that httpResp.StatusCode matches one of expected.
// It adds a friendly diagnostic message and returns false when the status is unexpected.
func ValidateStatusCode(diags *diag.Diagnostics, operation string, httpResp *http.Response, expected ...int) bool {
	if httpResp == nil {
		diags.AddError(
			fmt.Sprintf("Error %s", operation),
			"Missing HTTP response from NetBox API.",
		)
		return false
	}

	status := httpResp.StatusCode
	for _, code := range expected {
		if status == code {
			return true
		}
	}

	diags.AddError(
		fmt.Sprintf("Error %s", operation),
		fmt.Sprintf("Expected HTTP %s, got: %d", formatStatusList(expected), status),
	)
	return false
}

func formatStatusList(codes []int) string {
	if len(codes) == 0 {
		return "unknown"
	}

	parts := make([]string, len(codes))
	for i, code := range codes {
		parts[i] = fmt.Sprintf("%d", code)
	}

	if len(parts) == 1 {
		return parts[0]
	}
	if len(parts) == 2 {
		return parts[0] + " or " + parts[1]
	}

	return strings.Join(parts[:len(parts)-1], ", ") + ", or " + parts[len(parts)-1]
}
