package testutil

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// ValidationErrorTestConfig defines the configuration for validation error tests.
type ValidationErrorTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// InvalidConfig returns a Terraform configuration that should fail validation
	InvalidConfig func() string

	// ExpectedError is a regex pattern that the error message should match
	ExpectedError *regexp.Regexp

	// ExpectedErrorStrings are string patterns (any one should match)
	// Use this for simpler error matching when regex is overkill
	ExpectedErrorStrings []string

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// RunValidationErrorTest executes a test expecting a validation error.
// This verifies that invalid configurations are properly rejected.
func RunValidationErrorTest(t *testing.T, config ValidationErrorTestConfig) {
	t.Helper()

	step := resource.TestStep{
		Config: config.InvalidConfig(),
	}

	if config.ExpectedError != nil {
		step.ExpectError = config.ExpectedError
	} else if len(config.ExpectedErrorStrings) > 0 {
		// Build regex from strings (any match)
		pattern := ""
		for i, s := range config.ExpectedErrorStrings {
			if i > 0 {
				pattern += "|"
			}
			pattern += regexp.QuoteMeta(s)
		}
		step.ExpectError = regexp.MustCompile(pattern)
	}

	testCase := resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{step},
	}

	if config.CheckDestroy != nil {
		testCase.CheckDestroy = config.CheckDestroy
	}

	resource.Test(t, testCase)
}

// MultiValidationErrorTestConfig tests multiple invalid configurations.
type MultiValidationErrorTestConfig struct {
	// ResourceName is the Terraform resource type (e.g., "netbox_device")
	ResourceName string

	// TestCases maps test names to their configurations
	TestCases map[string]ValidationErrorCase

	// CheckDestroy function to verify resource cleanup (optional)
	CheckDestroy resource.TestCheckFunc
}

// ValidationErrorCase represents a single validation error test case.
type ValidationErrorCase struct {
	// Config returns the invalid Terraform configuration
	Config func() string

	// ExpectedError is a regex pattern that the error message should match
	ExpectedError *regexp.Regexp
}

// RunMultiValidationErrorTest runs multiple validation error tests as subtests.
func RunMultiValidationErrorTest(t *testing.T, config MultiValidationErrorTestConfig) {
	t.Helper()

	for name, tc := range config.TestCases {
		t.Run(name, func(t *testing.T) {
			testCase := resource.TestCase{
				PreCheck:                 func() { TestAccPreCheck(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.Config(),
						ExpectError: tc.ExpectedError,
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

// Common validation error patterns for reuse.
var (
	// ErrPatternRequired matches "required" validation errors.
	ErrPatternRequired = regexp.MustCompile(`(?i)required|must be specified|cannot be empty`)

	// ErrPatternInvalidValue matches invalid value errors.
	ErrPatternInvalidValue = regexp.MustCompile(`(?i)invalid|not valid|must be one of`)

	// ErrPatternInvalidFormat matches format validation errors.
	ErrPatternInvalidFormat = regexp.MustCompile(`(?i)invalid format|malformed|parse error|parseprefix|parseaddr|not a valid|network prefix|address with prefix|Internal Server Error|KeyError`)

	// ErrPatternInvalidIP matches IP address validation errors.
	ErrPatternInvalidIP = regexp.MustCompile(`(?i)invalid.*ip|invalid.*address|not a valid.*ip`)

	// ErrPatternInvalidURL matches URL validation errors.
	ErrPatternInvalidURL = regexp.MustCompile(`(?i)invalid.*url|malformed.*url`)

	// ErrPatternInvalidEnum matches enum validation errors.
	ErrPatternInvalidEnum = regexp.MustCompile(`(?i)must be one of|invalid.*value|expected.*got|not a valid choice`)

	// ErrPatternNotFound matches "not found" errors for invalid references.
	ErrPatternNotFound = regexp.MustCompile(`(?i)not found|does not exist|no.*found`)

	// ErrPatternConflict matches conflict errors.
	ErrPatternConflict = regexp.MustCompile(`(?i)conflict|mutually exclusive|cannot.*together`)

	// ErrPatternRange matches out-of-range errors.
	ErrPatternRange = regexp.MustCompile(`(?i)out of range|must be between|exceeds|minimum|maximum|less than or equal|greater than or equal`)

	// ErrPatternInconsistent matches provider inconsistency errors.
	ErrPatternInconsistent = regexp.MustCompile(`(?i)inconsistent result|produced an unexpected|was.*but now`)
)
