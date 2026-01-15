package testutil

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// LargeValueTestConfig tests handling of large values.
type LargeValueTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// FieldName is the field to test with large values
	FieldName string

	// ConfigWithValue returns config with the field set to the given value
	ConfigWithValue func(value string) string

	// LargeValue is the large value to test (e.g., 10000 character string)
	LargeValue string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunLargeValueTest tests that large values are handled correctly.
func RunLargeValueTest(t *testing.T, config LargeValueTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	steps := []resource.TestStep{
		{
			Config: config.ConfigWithValue(config.LargeValue),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckResourceAttr(resourceRef, config.FieldName, config.LargeValue),
			),
		},
		{
			Config:   config.ConfigWithValue(config.LargeValue),
			PlanOnly: true,
		},
	}

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	}

	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// SpecialCharacterTestConfig tests handling of special characters.
type SpecialCharacterTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// FieldName is the field to test with special characters
	FieldName string

	// ConfigWithValue returns config with the field set to the given value
	ConfigWithValue func(value string) string

	// TestCases maps test names to special character values
	TestCases map[string]string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunSpecialCharacterTests runs tests for various special characters.
func RunSpecialCharacterTests(t *testing.T, config SpecialCharacterTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	for name, value := range config.TestCases {
		t.Run(name, func(t *testing.T) {
			testCase := resource.TestCase{
				PreCheck:                 func() { TestAccPreCheck(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config.ConfigWithValue(value),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(resourceRef, "id"),
							resource.TestCheckResourceAttr(resourceRef, config.FieldName, value),
						),
					},
					{
						Config:   config.ConfigWithValue(value),
						PlanOnly: true,
					},
				},
			}

			if config.CheckDestroy != nil {
				testCase.CheckDestroy = config.CheckDestroy
			}

			resource.Test(t, testCase)
		})
	}
}

// CommonSpecialCharacterValues provides common special character test cases.
var CommonSpecialCharacterValues = map[string]string{
	"unicode_basic":       "Test with émojis 🎉 and ñ",
	"unicode_cjk":         "测试中文字符",
	"unicode_arabic":      "اختبار",
	"unicode_hebrew":      "בדיקה",
	"newlines":            "Line 1\nLine 2\nLine 3",
	"tabs":                "Col1\tCol2\tCol3",
	"quotes_single":       "It's a test with 'quotes'",
	"quotes_double":       `Test with "double quotes"`,
	"backslashes":         `Path\\to\\file`,
	"ampersand":           "Tom & Jerry",
	"angle_brackets":      "Value <tag> content </tag>",
	"special_html":        "5 > 3 && 2 < 4",
	"percent":             "100% complete",
	"hash":                "Issue #123",
	"at_symbol":           "user@example.com",
	"currency":            "Price: $100 or €85 or £75",
	"math_symbols":        "x² + y² = z²",
	"mixed_whitespace":    "  spaced  \t text  ",
	"empty_string":        "",
	"single_space":        " ",
	"leading_whitespace":  "  leading spaces",
	"trailing_whitespace": "trailing spaces  ",
}

// EmptyStringTestConfig tests handling of empty strings vs null.
type EmptyStringTestConfig struct {
	// ResourceName is the Terraform resource type
	ResourceName string

	// FieldName is the optional field to test
	FieldName string

	// ConfigWithEmptyString returns config with field set to ""
	ConfigWithEmptyString func() string

	// ConfigWithoutField returns config without the field
	ConfigWithoutField func() string

	// ExpectEmptyStringCleared indicates if empty string should clear the field
	// Some APIs treat "" the same as null, others preserve it
	ExpectEmptyStringCleared bool

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunEmptyStringTest tests empty string handling.
func RunEmptyStringTest(t *testing.T, config EmptyStringTestConfig) {
	t.Helper()

	resourceRef := fmt.Sprintf("%s.test", config.ResourceName)

	var emptyStringCheck resource.TestCheckFunc
	if config.ExpectEmptyStringCleared {
		emptyStringCheck = resource.TestCheckNoResourceAttr(resourceRef, config.FieldName)
	} else {
		emptyStringCheck = resource.TestCheckResourceAttr(resourceRef, config.FieldName, "")
	}

	steps := []resource.TestStep{
		// Create with empty string
		{
			Config: config.ConfigWithEmptyString(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				emptyStringCheck,
			),
		},
		// Verify no drift
		{
			Config:   config.ConfigWithEmptyString(),
			PlanOnly: true,
		},
		// Remove field entirely
		{
			Config: config.ConfigWithoutField(),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceRef, "id"),
				resource.TestCheckNoResourceAttr(resourceRef, config.FieldName),
			),
		},
	}

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	}

	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// GenerateLargeString creates a string of the specified length.
func GenerateLargeString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = byte('a' + (i % 26))
	}
	return string(result)
}

// GenerateLargeDescription creates a realistic large description.
func GenerateLargeDescription(paragraphs int) string {
	paragraph := "This is a test description paragraph that contains realistic content. " +
		"It includes multiple sentences with various punctuation marks, numbers like 12345, " +
		"and special characters such as dashes - and parentheses (like this). " +
		"The purpose is to test how the system handles moderately large text content.\n\n"

	result := ""
	for i := 0; i < paragraphs; i++ {
		result += paragraph
	}
	return result
}
