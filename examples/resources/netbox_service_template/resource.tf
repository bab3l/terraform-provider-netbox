resource "netbox_service_template" "test" {
  name        = "SSH"
  protocol    = "tcp"
  ports       = [22]
  description = "Standard SSH service template"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "default_monitoring"
      value = "true"
    },
    {
      name  = "standard_port"
      value = "22"
    },
    {
      name  = "security_category"
      value = "remote-access"
    }
  ]

  tags = [
    "service-template",
    "standard",
    "ssh"
  ]
}
