resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
  u_height     = 1
}

resource "netbox_inventory_item_template" "test" {
  name        = "Inventory Item Template"
  device_type = netbox_device_type.test.id

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "component_type"
      value = "power-supply"
    },
    {
      name  = "warranty_years"
      value = "3"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_inventory_item_template.test
  id = "123"

  identity = {
    custom_fields = [
      "component_type:text",
      "warranty_years:integer",
    ]
  }
}
