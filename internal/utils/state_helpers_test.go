// Package utils provides utility functions for the Terraform provider.
package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int32
		wantError bool
	}{
		{
			name:      "valid positive ID",
			input:     "123",
			want:      123,
			wantError: false,
		},
		{
			name:      "valid ID with value 1",
			input:     "1",
			want:      1,
			wantError: false,
		},
		{
			name:      "valid large ID",
			input:     "2147483647", // max int32
			want:      2147483647,
			wantError: false,
		},
		{
			name:      "valid ID with leading zeros",
			input:     "00123",
			want:      123,
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			want:      0,
			wantError: true,
		},
		{
			name:      "non-numeric string",
			input:     "abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "mixed string",
			input:     "123abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "negative ID",
			input:     "-123",
			want:      -123,
			wantError: false,
		},
		{
			name:      "overflow int32",
			input:     "2147483648", // max int32 + 1
			want:      0,
			wantError: true,
		},
		{
			name:      "underflow int32",
			input:     "-2147483649", // min int32 - 1
			want:      0,
			wantError: true,
		},
		{
			name:      "float string",
			input:     "123.45",
			want:      0,
			wantError: true,
		},
		{
			name:      "string with spaces",
			input:     " 123 ",
			want:      0,
			wantError: true,
		},
		{
			name:      "zero",
			input:     "0",
			want:      0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseID(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseID(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseID(%q) unexpected error: %v", tt.input, err)
				return
			}

			if got != tt.want {
				t.Errorf("ParseID(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestMustParseID(t *testing.T) {
	// Test valid case
	t.Run("valid ID", func(t *testing.T) {
		got := MustParseID("123")
		if got != 123 {
			t.Errorf("MustParseID(%q) = %d, want %d", "123", got, 123)
		}
	})

	// Test panic case
	t.Run("invalid ID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustParseID(%q) did not panic", "invalid")
			}
		}()
		MustParseID("invalid")
	})
}

func TestParseInt32(t *testing.T) {
	// Note: ParseInt32 silently returns 0 on error, so we test that behavior
	tests := []struct {
		name  string
		input string
		want  int32
	}{
		{
			name:  "valid number",
			input: "123",
			want:  123,
		},
		{
			name:  "invalid string returns 0",
			input: "abc",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ParseInt32 takes types.String, so we need to use ParseInt32FromString for string input
			got := ParseInt32FromString(tt.input)
			if got != tt.want {
				t.Errorf("ParseInt32FromString(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseInt32FromString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int32
	}{
		{
			name:  "valid number",
			input: "456",
			want:  456,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
		{
			name:  "negative number",
			input: "-789",
			want:  -789,
		},
		{
			name:  "invalid returns 0",
			input: "not-a-number",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInt32FromString(tt.input)
			if got != tt.want {
				t.Errorf("ParseInt32FromString(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// =====================================================
// STRING FROM API TESTS
// =====================================================

func TestStringFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() string
		current  types.String
		want     types.String
	}{
		{
			name:     "has value with non-empty string",
			hasValue: true,
			getValue: func() string { return "test-value" },
			current:  types.StringNull(),
			want:     types.StringValue("test-value"),
		},
		{
			name:     "has value with empty string returns null",
			hasValue: true,
			getValue: func() string { return "" },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "has value with empty string, current not null returns null",
			hasValue: true,
			getValue: func() string { return "" },
			current:  types.StringValue("old-value"),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current is null stays null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringValue("old-value"),
			want:     types.StringNull(),
		},
		{
			name:     "has value preserves whitespace",
			hasValue: true,
			getValue: func() string { return "  test  " },
			current:  types.StringNull(),
			want:     types.StringValue("  test  "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringFromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("StringFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringFromAPIPreserveEmpty(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() string
		current  types.String
		want     types.String
	}{
		{
			name:     "has value with non-empty string",
			hasValue: true,
			getValue: func() string { return "test-value" },
			current:  types.StringNull(),
			want:     types.StringValue("test-value"),
		},
		{
			name:     "has value with empty string preserves empty",
			hasValue: true,
			getValue: func() string { return "" },
			current:  types.StringNull(),
			want:     types.StringValue(""),
		},
		{
			name:     "no value, current is null stays null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringValue("old-value"),
			want:     types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringFromAPIPreserveEmpty(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("StringFromAPIPreserveEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableStringFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() string
		current  types.String
		want     types.String
	}{
		{
			name:     "has value with non-empty string",
			hasValue: true,
			getValue: func() string { return "nullable-value" },
			current:  types.StringNull(),
			want:     types.StringValue("nullable-value"),
		},
		{
			name:     "has value with empty string returns null",
			hasValue: true,
			getValue: func() string { return "" },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current is null stays null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() string { return "ignored" },
			current:  types.StringValue("old-value"),
			want:     types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableStringFromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("NullableStringFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// INTEGER FROM API TESTS
// =====================================================

func TestInt64FromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() int64
		current  types.Int64
		want     types.Int64
	}{
		{
			name:     "has value returns value",
			hasValue: true,
			getValue: func() int64 { return 42 },
			current:  types.Int64Null(),
			want:     types.Int64Value(42),
		},
		{
			name:     "has value zero returns zero",
			hasValue: true,
			getValue: func() int64 { return 0 },
			current:  types.Int64Null(),
			want:     types.Int64Value(0),
		},
		{
			name:     "has negative value",
			hasValue: true,
			getValue: func() int64 { return -100 },
			current:  types.Int64Null(),
			want:     types.Int64Value(-100),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() int64 { return 999 },
			current:  types.Int64Null(),
			want:     types.Int64Null(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() int64 { return 999 },
			current:  types.Int64Value(123),
			want:     types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int64FromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("Int64FromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64FromInt32API(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() int32
		current  types.Int64
		want     types.Int64
	}{
		{
			name:     "has value returns value as int64",
			hasValue: true,
			getValue: func() int32 { return 42 },
			current:  types.Int64Null(),
			want:     types.Int64Value(42),
		},
		{
			name:     "has value max int32",
			hasValue: true,
			getValue: func() int32 { return 2147483647 },
			current:  types.Int64Null(),
			want:     types.Int64Value(2147483647),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() int32 { return 999 },
			current:  types.Int64Null(),
			want:     types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int64FromInt32API(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("Int64FromInt32API() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableInt64FromAPI(t *testing.T) {
	val42 := int32(42)
	tests := []struct {
		name     string
		hasValue bool
		getValue func() *int32
		current  types.Int64
		want     types.Int64
	}{
		{
			name:     "has non-nil pointer",
			hasValue: true,
			getValue: func() *int32 { return &val42 },
			current:  types.Int64Null(),
			want:     types.Int64Value(42),
		},
		{
			name:     "has nil pointer",
			hasValue: true,
			getValue: func() *int32 { return nil },
			current:  types.Int64Null(),
			want:     types.Int64Null(),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() *int32 { return &val42 },
			current:  types.Int64Null(),
			want:     types.Int64Null(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() *int32 { return &val42 },
			current:  types.Int64Value(123),
			want:     types.Int64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableInt64FromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("NullableInt64FromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// FLOAT FROM API TESTS
// =====================================================

func TestFloat64FromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() float64
		current  types.Float64
		want     types.Float64
	}{
		{
			name:     "has value returns value",
			hasValue: true,
			getValue: func() float64 { return 3.14159 },
			current:  types.Float64Null(),
			want:     types.Float64Value(3.14159),
		},
		{
			name:     "has value zero returns zero",
			hasValue: true,
			getValue: func() float64 { return 0.0 },
			current:  types.Float64Null(),
			want:     types.Float64Value(0.0),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() float64 { return 999.99 },
			current:  types.Float64Null(),
			want:     types.Float64Null(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() float64 { return 999.99 },
			current:  types.Float64Value(123.45),
			want:     types.Float64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Float64FromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("Float64FromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullableFloat64FromAPI(t *testing.T) {
	val := 3.14159
	tests := []struct {
		name     string
		hasValue bool
		getValue func() *float64
		current  types.Float64
		want     types.Float64
	}{
		{
			name:     "has non-nil pointer",
			hasValue: true,
			getValue: func() *float64 { return &val },
			current:  types.Float64Null(),
			want:     types.Float64Value(3.14159),
		},
		{
			name:     "has nil pointer",
			hasValue: true,
			getValue: func() *float64 { return nil },
			current:  types.Float64Null(),
			want:     types.Float64Null(),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() *float64 { return &val },
			current:  types.Float64Null(),
			want:     types.Float64Null(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableFloat64FromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("NullableFloat64FromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// BOOL FROM API TESTS
// =====================================================

func TestBoolFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() bool
		current  types.Bool
		want     types.Bool
	}{
		{
			name:     "has true value",
			hasValue: true,
			getValue: func() bool { return true },
			current:  types.BoolNull(),
			want:     types.BoolValue(true),
		},
		{
			name:     "has false value",
			hasValue: true,
			getValue: func() bool { return false },
			current:  types.BoolNull(),
			want:     types.BoolValue(false),
		},
		{
			name:     "no value, current null stays null",
			hasValue: false,
			getValue: func() bool { return true },
			current:  types.BoolNull(),
			want:     types.BoolNull(),
		},
		{
			name:     "no value, current has value returns null",
			hasValue: false,
			getValue: func() bool { return true },
			current:  types.BoolValue(true),
			want:     types.BoolNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolFromAPI(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("BoolFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// REFERENCE FIELD TESTS
// =====================================================

func TestReferenceIDFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getID    func() int32
		current  types.String
		want     types.String
	}{
		{
			name:     "has value, current null returns ID as string",
			hasValue: true,
			getID:    func() int32 { return 123 },
			current:  types.StringNull(),
			want:     types.StringValue("123"),
		},
		{
			name:     "has value, current unknown returns ID as string",
			hasValue: true,
			getID:    func() int32 { return 456 },
			current:  types.StringUnknown(),
			want:     types.StringValue("456"),
		},
		{
			name:     "has value, current has value preserves current",
			hasValue: true,
			getID:    func() int32 { return 789 },
			current:  types.StringValue("my-slug"),
			want:     types.StringValue("my-slug"),
		},
		{
			name:     "has value zero, current null returns null",
			hasValue: true,
			getID:    func() int32 { return 0 },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current null returns null",
			hasValue: false,
			getID:    func() int32 { return 999 },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
		{
			name:     "no value, current has value preserves current",
			hasValue: false,
			getID:    func() int32 { return 999 },
			current:  types.StringValue("existing"),
			want:     types.StringValue("existing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReferenceIDFromAPI(tt.hasValue, tt.getID, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("ReferenceIDFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequiredReferenceIDFromAPI(t *testing.T) {
	tests := []struct {
		name    string
		getID   func() int32
		current types.String
		want    types.String
	}{
		{
			name:    "current null returns ID as string",
			getID:   func() int32 { return 123 },
			current: types.StringNull(),
			want:    types.StringValue("123"),
		},
		{
			name:    "current unknown returns ID as string",
			getID:   func() int32 { return 456 },
			current: types.StringUnknown(),
			want:    types.StringValue("456"),
		},
		{
			name:    "current has value preserves current",
			getID:   func() int32 { return 789 },
			current: types.StringValue("my-slug"),
			want:    types.StringValue("my-slug"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RequiredReferenceIDFromAPI(tt.getID, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("RequiredReferenceIDFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// ENUM FROM API TESTS
// =====================================================

type TestEnum string

const (
	TestEnumActive   TestEnum = "active"
	TestEnumInactive TestEnum = "inactive"
)

func TestEnumFromAPI(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() TestEnum
		want     types.String
	}{
		{
			name:     "has value returns string",
			hasValue: true,
			getValue: func() TestEnum { return TestEnumActive },
			want:     types.StringValue("active"),
		},
		{
			name:     "no value returns null",
			hasValue: false,
			getValue: func() TestEnum { return TestEnumInactive },
			want:     types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnumFromAPI(tt.hasValue, tt.getValue)
			if !got.Equal(tt.want) {
				t.Errorf("EnumFromAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnumFromAPIWithDefault(t *testing.T) {
	tests := []struct {
		name     string
		hasValue bool
		getValue func() TestEnum
		current  types.String
		want     types.String
	}{
		{
			name:     "has value returns value",
			hasValue: true,
			getValue: func() TestEnum { return TestEnumActive },
			current:  types.StringValue("old"),
			want:     types.StringValue("active"),
		},
		{
			name:     "no value returns current",
			hasValue: false,
			getValue: func() TestEnum { return TestEnumInactive },
			current:  types.StringValue("default"),
			want:     types.StringValue("default"),
		},
		{
			name:     "no value, current null returns null",
			hasValue: false,
			getValue: func() TestEnum { return TestEnumInactive },
			current:  types.StringNull(),
			want:     types.StringNull(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnumFromAPIWithDefault(tt.hasValue, tt.getValue, tt.current)
			if !got.Equal(tt.want) {
				t.Errorf("EnumFromAPIWithDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =====================================================
// REQUEST BUILDING HELPER TESTS
// =====================================================

func TestIsSet(t *testing.T) {
	tests := []struct {
		name  string
		value types.String
		want  bool
	}{
		{
			name:  "non-null value is set",
			value: types.StringValue("test"),
			want:  true,
		},
		{
			name:  "empty string value is set",
			value: types.StringValue(""),
			want:  true,
		},
		{
			name:  "null value is not set",
			value: types.StringNull(),
			want:  false,
		},
		{
			name:  "unknown value is not set",
			value: types.StringUnknown(),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSet(tt.value)
			if got != tt.want {
				t.Errorf("IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSetWithDifferentTypes(t *testing.T) {
	t.Run("Int64 set", func(t *testing.T) {
		if !IsSet(types.Int64Value(42)) {
			t.Error("Int64Value should be set")
		}
	})

	t.Run("Int64 null", func(t *testing.T) {
		if IsSet(types.Int64Null()) {
			t.Error("Int64Null should not be set")
		}
	})

	t.Run("Bool set", func(t *testing.T) {
		if !IsSet(types.BoolValue(true)) {
			t.Error("BoolValue should be set")
		}
	})

	t.Run("Bool null", func(t *testing.T) {
		if IsSet(types.BoolNull()) {
			t.Error("BoolNull should not be set")
		}
	})

	t.Run("Float64 set", func(t *testing.T) {
		if !IsSet(types.Float64Value(3.14)) {
			t.Error("Float64Value should be set")
		}
	})

	t.Run("Float64 null", func(t *testing.T) {
		if IsSet(types.Float64Null()) {
			t.Error("Float64Null should not be set")
		}
	})
}

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name    string
		value   types.String
		wantNil bool
		wantVal string
	}{
		{
			name:    "non-null returns pointer",
			value:   types.StringValue("test"),
			wantNil: false,
			wantVal: "test",
		},
		{
			name:    "empty string returns pointer to empty",
			value:   types.StringValue(""),
			wantNil: false,
			wantVal: "",
		},
		{
			name:    "null returns nil",
			value:   types.StringNull(),
			wantNil: true,
		},
		{
			name:    "unknown returns nil",
			value:   types.StringUnknown(),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringPtr(tt.value)
			if tt.wantNil {
				if got != nil {
					t.Errorf("StringPtr() = %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Errorf("StringPtr() = nil, want %q", tt.wantVal)
				} else if *got != tt.wantVal {
					t.Errorf("StringPtr() = %q, want %q", *got, tt.wantVal)
				}
			}
		})
	}
}

func TestInt32Ptr(t *testing.T) {
	tests := []struct {
		name    string
		value   types.Int64
		wantNil bool
		wantVal int32
	}{
		{
			name:    "non-null returns pointer",
			value:   types.Int64Value(42),
			wantNil: false,
			wantVal: 42,
		},
		{
			name:    "zero returns pointer to zero",
			value:   types.Int64Value(0),
			wantNil: false,
			wantVal: 0,
		},
		{
			name:    "negative returns pointer",
			value:   types.Int64Value(-100),
			wantNil: false,
			wantVal: -100,
		},
		{
			name:    "null returns nil",
			value:   types.Int64Null(),
			wantNil: true,
		},
		{
			name:    "unknown returns nil",
			value:   types.Int64Unknown(),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int32Ptr(tt.value)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Int32Ptr() = %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Errorf("Int32Ptr() = nil, want %d", tt.wantVal)
				} else if *got != tt.wantVal {
					t.Errorf("Int32Ptr() = %d, want %d", *got, tt.wantVal)
				}
			}
		})
	}
}

func TestInt32Value(t *testing.T) {
	tests := []struct {
		name  string
		value types.Int64
		want  int32
	}{
		{
			name:  "non-null returns value",
			value: types.Int64Value(42),
			want:  42,
		},
		{
			name:  "null returns 0",
			value: types.Int64Null(),
			want:  0,
		},
		{
			name:  "unknown returns 0",
			value: types.Int64Unknown(),
			want:  0,
		},
		{
			name:  "negative value",
			value: types.Int64Value(-999),
			want:  -999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int32Value(tt.value)
			if got != tt.want {
				t.Errorf("Int32Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFloat64Ptr(t *testing.T) {
	tests := []struct {
		name    string
		value   types.Float64
		wantNil bool
		wantVal float64
	}{
		{
			name:    "non-null returns pointer",
			value:   types.Float64Value(3.14159),
			wantNil: false,
			wantVal: 3.14159,
		},
		{
			name:    "zero returns pointer to zero",
			value:   types.Float64Value(0.0),
			wantNil: false,
			wantVal: 0.0,
		},
		{
			name:    "null returns nil",
			value:   types.Float64Null(),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Float64Ptr(tt.value)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Float64Ptr() = %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Errorf("Float64Ptr() = nil, want %f", tt.wantVal)
				} else if *got != tt.wantVal {
					t.Errorf("Float64Ptr() = %f, want %f", *got, tt.wantVal)
				}
			}
		})
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name    string
		value   types.Bool
		wantNil bool
		wantVal bool
	}{
		{
			name:    "true returns pointer to true",
			value:   types.BoolValue(true),
			wantNil: false,
			wantVal: true,
		},
		{
			name:    "false returns pointer to false",
			value:   types.BoolValue(false),
			wantNil: false,
			wantVal: false,
		},
		{
			name:    "null returns nil",
			value:   types.BoolNull(),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolPtr(tt.value)
			if tt.wantNil {
				if got != nil {
					t.Errorf("BoolPtr() = %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Errorf("BoolPtr() = nil, want %v", tt.wantVal)
				} else if *got != tt.wantVal {
					t.Errorf("BoolPtr() = %v, want %v", *got, tt.wantVal)
				}
			}
		})
	}
}

func TestParseInt32_TypesString(t *testing.T) {
	tests := []struct {
		name  string
		value types.String
		want  int32
	}{
		{
			name:  "valid number",
			value: types.StringValue("123"),
			want:  123,
		},
		{
			name:  "null returns 0",
			value: types.StringNull(),
			want:  0,
		},
		{
			name:  "unknown returns 0",
			value: types.StringUnknown(),
			want:  0,
		},
		{
			name:  "invalid string returns 0",
			value: types.StringValue("abc"),
			want:  0,
		},
		{
			name:  "negative number",
			value: types.StringValue("-456"),
			want:  -456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInt32(tt.value)
			if got != tt.want {
				t.Errorf("ParseInt32() = %d, want %d", got, tt.want)
			}
		})
	}
}
