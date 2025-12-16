// Package utils provides utility functions for the Terraform provider.

package utils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapToCustomFieldModels_Panic(t *testing.T) {

	// Setup

	customFields := map[string]interface{}{

		"int_field": 123.0, // float64 as unmarshaled from JSON

		"bool_field": true,

		"str_field": "some string",

		"multi_field": []interface{}{"a", 1, true},
	}

	stateCFs := []CustomFieldModel{

		{

			Name: types.StringValue("int_field"),

			Type: types.StringValue("integer"),
		},

		{

			Name: types.StringValue("bool_field"),

			Type: types.StringValue("boolean"),
		},

		{

			Name: types.StringValue("str_field"),

			Type: types.StringValue("text"),
		},

		{

			Name: types.StringValue("multi_field"),

			Type: types.StringValue("multiselect"),
		},
	}

	// Execute

	// This should not panic

	result := MapToCustomFieldModels(customFields, stateCFs)

	// Assert

	if len(result) != 4 {

		t.Errorf("Expected 4 results, got %d", len(result))

	}

	for _, cf := range result {

		if cf.Name.ValueString() == "int_field" {

			if cf.Value.ValueString() != "123" {

				t.Errorf("Expected int_field to be '123', got '%s'", cf.Value.ValueString())

			}

		}

		if cf.Name.ValueString() == "bool_field" {

			if cf.Value.ValueString() != "true" {

				t.Errorf("Expected bool_field to be 'true', got '%s'", cf.Value.ValueString())

			}

		}

		if cf.Name.ValueString() == "multi_field" {

			// Expecting "a,1,true" or similar depending on formatting

			if cf.Value.ValueString() != "a,1,true" {

				t.Errorf("Expected multi_field to be 'a,1,true', got '%s'", cf.Value.ValueString())

			}

		}

	}

}

func TestParseDuplicateErrorFromBytes(t *testing.T) {

	tests := []struct {
		name string

		statusCode int

		body string

		wantNil bool

		wantSlug bool

		wantName bool
	}{

		{

			name: "404 error should not be duplicate",

			statusCode: 404,

			body: `{"detail": "Not found."}`,

			wantNil: true,
		},

		{

			name: "400 with slug already exists",

			statusCode: 400,

			body: `{"slug": ["site with this slug already exists."]}`,

			wantNil: false,

			wantSlug: true,
		},

		{

			name: "400 with name already exists",

			statusCode: 400,

			body: `{"name": ["site with this name already exists."]}`,

			wantNil: false,

			wantName: true,
		},

		{

			name: "400 with both slug and name already exist",

			statusCode: 400,

			body: `{"slug": ["site with this slug already exists."], "name": ["site with this name already exists."]}`,

			wantNil: false,

			wantSlug: true,

			wantName: true,
		},

		{

			name: "400 with unrelated validation error",

			statusCode: 400,

			body: `{"status": ["This field is required."]}`,

			wantNil: true,
		},

		{

			name: "200 success should not be duplicate",

			statusCode: 200,

			body: `{"id": 1, "name": "test"}`,

			wantNil: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			resp := &http.Response{

				StatusCode: tt.statusCode,

				Body: io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			result := parseDuplicateErrorFromBytes(resp, []byte(tt.body))

			if tt.wantNil {

				if result != nil {

					t.Errorf("expected nil, got %v", result)

				}

				return

			}

			if result == nil {

				t.Errorf("expected non-nil result")

				return

			}

			if tt.wantSlug {

				if _, ok := result["slug"]; !ok {

					t.Errorf("expected slug in result, got %v", result)

				}

			}

			if tt.wantName {

				if _, ok := result["name"]; !ok {

					t.Errorf("expected name in result, got %v", result)

				}

			}

		})

	}

}

func TestCreateErrorHandler_HandleCreateError_Duplicate(t *testing.T) {

	ctx := context.Background()

	// Test with lookup that finds an existing resource

	handler := CreateErrorHandler{

		ResourceType: "netbox_tenant",

		ResourceName: "test_tenant",

		SlugValue: "my-slug",

		LookupFunc: func(ctx context.Context, slug string) (string, error) {

			// Simulate finding an existing resource with ID 42

			return "42", nil

		},
	}

	// Create a mock HTTP response with duplicate error

	resp := &http.Response{

		StatusCode: 400,

		Body: io.NopCloser(bytes.NewBufferString(`{"slug": ["tenant with this slug already exists."]}`)),
	}

	var diags diag.Diagnostics

	handler.HandleCreateError(ctx, nil, resp, &diags)

	if !diags.HasError() {

		t.Errorf("expected diagnostics to have error")

	}

	// Check that the error message contains import instructions

	for _, d := range diags.Errors() {

		if d.Summary() != "Duplicate netbox_tenant" {

			t.Errorf("expected summary 'Duplicate netbox_tenant', got '%s'", d.Summary())

		}

		detail := d.Detail()

		if !contains(detail, "terraform import") {

			t.Errorf("expected import command in detail, got '%s'", detail)

		}

		if !contains(detail, "42") {

			t.Errorf("expected resource ID 42 in detail, got '%s'", detail)

		}

	}

}

func TestCreateErrorHandler_HandleCreateError_DuplicateWithLookupFailure(t *testing.T) {

	ctx := context.Background()

	// Test with lookup that fails

	handler := CreateErrorHandler{

		ResourceType: "netbox_site",

		ResourceName: "test_site",

		SlugValue: "my-slug",

		LookupFunc: func(ctx context.Context, slug string) (string, error) {

			// Simulate lookup failure

			return "", nil

		},
	}

	// Create a mock HTTP response with duplicate error

	resp := &http.Response{

		StatusCode: 400,

		Body: io.NopCloser(bytes.NewBufferString(`{"slug": ["site with this slug already exists."]}`)),
	}

	var diags diag.Diagnostics

	handler.HandleCreateError(ctx, nil, resp, &diags)

	if !diags.HasError() {

		t.Errorf("expected diagnostics to have error")

	}

	// Check that the error message contains import hints even without ID

	for _, d := range diags.Errors() {

		detail := d.Detail()

		if !contains(detail, "find it in Netbox") {

			t.Errorf("expected 'find it in Netbox' hint in detail, got '%s'", detail)

		}

		if !contains(detail, "my-slug") {

			t.Errorf("expected slug 'my-slug' in detail, got '%s'", detail)

		}

	}

}

func TestCreateErrorHandler_HandleCreateError_NonDuplicate(t *testing.T) {

	ctx := context.Background()

	handler := CreateErrorHandler{

		ResourceType: "netbox_tenant",

		ResourceName: "test_tenant",

		SlugValue: "my-slug",
	}

	// Create a mock HTTP response with a non-duplicate error

	resp := &http.Response{

		StatusCode: 400,

		Body: io.NopCloser(bytes.NewBufferString(`{"group": ["Invalid group ID."]}`)),
	}

	var diags diag.Diagnostics

	handler.HandleCreateError(ctx, nil, resp, &diags)

	if !diags.HasError() {

		t.Errorf("expected diagnostics to have error")

	}

	// Check that it's a normal error, not a duplicate error

	for _, d := range diags.Errors() {

		if d.Summary() == "Duplicate netbox_tenant" {

			t.Errorf("should not be a duplicate error")

		}

	}

}

func contains(s, substr string) bool {

	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))

}

func containsSubstring(s, substr string) bool {

	for i := 0; i <= len(s)-len(substr); i++ {

		if s[i:i+len(substr)] == substr {

			return true

		}

	}

	return false

}
