resource "netbox_circuit_type" "test" {
  name        = "Internet Transit"
  slug        = "internet-transit"
  description = "Internet transit circuits"

  # Partial custom fields management
  # Only the custom fields specified here are managed by Terraform
  # Other custom fields set in NetBox are preserved
  custom_fields = [
    {
      name  = "default_speed_tier"
      value = "1G"
    },
    {
      name  = "billing_model"
      value = "95th-percentile"
    }
  ]

  tags = [
    "circuit-type"
  ]
}
