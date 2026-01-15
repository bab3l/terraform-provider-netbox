package testutil

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Shared test constants used across acceptance tests

const (
	Comments      = "Test comments"
	Description1  = "Initial description"
	Description2  = "Updated description"
	RearPortName  = "rear0"
	Color         = "aa1409"
	ColorOrange   = "ff5722"
	InterfaceName = "eth0"

	// Device/Interface status values.
	StatusActive  = "active"
	StatusPlanned = "planned"
	StatusStaged  = "staged"
	StatusFailed  = "failed"
	StatusOffline = "offline"

	// Interface type constants.
	InterfaceType1000BaseT   = "1000base-t"
	InterfaceType10GBaseSFPP = "10gbase-x-sfpp"
	InterfaceType10GBaseT    = "10gbase-t"
	InterfaceType25GBaseSFP  = "25gbase-x-sfp28"

	// Port type constants.
	PortType8P8C = "8p8c"
	PortTypeLC   = "lc"
	PortTypeSC   = "sc"
	PortTypeST   = "st"

	// Power port type constants.
	PowerPortTypeIEC60320C14 = "iec-60320-c14"
	PowerPortTypeIEC60320C20 = "iec-60320-c20"

	// Prefix status constants.
	PrefixStatusContainer  = "container"
	PrefixStatusActive     = "active"
	PrefixStatusReserved   = "reserved"
	PrefixStatusDeprecated = "deprecated"
)

// CheckCustomFieldValue returns a TestCheckFunc that verifies a custom field's name, type, and value.
// It searches through the custom_fields set to find a field with the given name and verifies its properties.
func CheckCustomFieldValue(resourceName, fieldName, fieldType, fieldValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		// Custom fields are stored as a set, we need to find the specific field
		// Check if custom_fields attribute exists
		customFieldsCount, ok := rs.Primary.Attributes["custom_fields.#"]
		if !ok {
			return fmt.Errorf("custom_fields attribute not found on %s", resourceName)
		}

		if customFieldsCount == "0" {
			return fmt.Errorf("no custom fields found on %s", resourceName)
		}

		// Search through all custom fields to find the one with matching name
		// Format is: custom_fields.{index}.{attribute}
		found := false
		for key, value := range rs.Primary.Attributes {
			if value == fieldName {
				// Found the matching name, extract the index
				// key format: custom_fields.0.name, custom_fields.1.name, etc.
				prefix := key[:len(key)-len(".name")]

				// Verify type and value
				actualType := rs.Primary.Attributes[prefix+".type"]
				if actualType != fieldType {
					return fmt.Errorf("custom field %s has type %s, expected %s", fieldName, actualType, fieldType)
				}

				actualValue := rs.Primary.Attributes[prefix+".value"]
				if actualValue != fieldValue {
					return fmt.Errorf("custom field %s has value %s, expected %s", fieldName, actualValue, fieldValue)
				}

				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("custom field %s not found on %s", fieldName, resourceName)
		}

		return nil
	}
}
