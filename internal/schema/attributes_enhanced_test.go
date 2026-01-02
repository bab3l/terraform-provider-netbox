package schema

import (
	"testing"
)

// TestReferenceAttributeWithDiffSuppress tests the enhanced ReferenceAttribute function
// that includes diff suppression logic.
func TestReferenceAttributeWithDiffSuppress(t *testing.T) {
	tests := []struct {
		name             string
		targetResource   string
		description      string
		wantOptional     bool
		wantRequired     bool
		wantDiffSuppress bool
		wantMarkdown     string
	}{
		{
			name:             "basic_tenant_reference",
			targetResource:   "tenant",
			description:      "ID or slug of the tenant.",
			wantOptional:     true,
			wantRequired:     false,
			wantDiffSuppress: true,
			wantMarkdown:     "ID or slug of the tenant.",
		},
		{
			name:             "site_reference_with_custom_desc",
			targetResource:   "site",
			description:      "Custom site description for testing.",
			wantOptional:     true,
			wantRequired:     false,
			wantDiffSuppress: true,
			wantMarkdown:     "Custom site description for testing.",
		},
		{
			name:             "device_type_reference_auto_desc",
			targetResource:   "device type",
			description:      "", // Should auto-generate
			wantOptional:     true,
			wantRequired:     false,
			wantDiffSuppress: true,
			wantMarkdown:     "ID or slug of the device type.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := ReferenceAttributeWithDiffSuppress(tt.targetResource, tt.description)

			// Check optional/required flags
			if attr.Optional != tt.wantOptional {
				t.Errorf("ReferenceAttributeWithDiffSuppress().Optional = %v, want %v", attr.Optional, tt.wantOptional)
			}

			if attr.Required != tt.wantRequired {
				t.Errorf("ReferenceAttributeWithDiffSuppress().Required = %v, want %v", attr.Required, tt.wantRequired)
			}

			// Check markdown description
			if attr.MarkdownDescription != tt.wantMarkdown {
				t.Errorf("ReferenceAttributeWithDiffSuppress().MarkdownDescription = %q, want %q",
					attr.MarkdownDescription, tt.wantMarkdown)
			}

			// Verify plan modifiers are present
			if len(attr.PlanModifiers) == 0 {
				t.Error("ReferenceAttributeWithDiffSuppress() should include plan modifiers for diff suppression")
			}
		})
	}
}

// TestRequiredReferenceAttributeWithDiffSuppress tests the enhanced RequiredReferenceAttribute function.
func TestRequiredReferenceAttributeWithDiffSuppress(t *testing.T) {
	tests := []struct {
		name           string
		targetResource string
		description    string
		wantOptional   bool
		wantRequired   bool
		wantMarkdown   string
	}{
		{
			name:           "required_device_type",
			targetResource: "device type",
			description:    "Required device type reference.",
			wantOptional:   false,
			wantRequired:   true,
			wantMarkdown:   "Required device type reference.",
		},
		{
			name:           "required_site_auto_desc",
			targetResource: "site",
			description:    "", // Should auto-generate
			wantOptional:   false,
			wantRequired:   true,
			wantMarkdown:   "ID or slug of the site. Required.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := RequiredReferenceAttributeWithDiffSuppress(tt.targetResource, tt.description)

			// Check optional/required flags
			if attr.Optional != tt.wantOptional {
				t.Errorf("RequiredReferenceAttributeWithDiffSuppress().Optional = %v, want %v", attr.Optional, tt.wantOptional)
			}

			if attr.Required != tt.wantRequired {
				t.Errorf("RequiredReferenceAttributeWithDiffSuppress().Required = %v, want %v", attr.Required, tt.wantRequired)
			}

			// Check markdown description
			if attr.MarkdownDescription != tt.wantMarkdown {
				t.Errorf("RequiredReferenceAttributeWithDiffSuppress().MarkdownDescription = %q, want %q",
					attr.MarkdownDescription, tt.wantMarkdown)
			}
		})
	}
}

// TestBackwardCompatibility ensures existing ReferenceAttribute function still works.
func TestBackwardCompatibility(t *testing.T) {
	// Test that the original ReferenceAttribute function still works as expected
	attr := ReferenceAttribute("tenant", "ID or slug of the tenant.")

	if !attr.Optional {
		t.Error("ReferenceAttribute() should create optional attribute")
	}

	if attr.Required {
		t.Error("ReferenceAttribute() should not create required attribute")
	}

	if attr.MarkdownDescription != "ID or slug of the tenant." {
		t.Errorf("ReferenceAttribute().MarkdownDescription = %q, want %q",
			attr.MarkdownDescription, "ID or slug of the tenant.")
	}
}
