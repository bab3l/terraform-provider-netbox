# Example: Basic contact group
resource "netbox_contact_group" "basic" {
  name        = "IT Department"
  slug        = "it-department"
  description = "Information Technology department contacts"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "department_code"
      value = "IT-100"
    },
    {
      name  = "team_size"
      value = "15"
    },
    {
      name  = "escalation_email"
      value = "it-escalation@example.com"
    }
  ]

  tags = [
    "contact-group",
    "it-department"
  ]
}

# Example: Contact group with description
resource "netbox_contact_group" "with_description" {
  name        = "Network Operations"
  slug        = "network-operations"
  description = "Network operations and support team"
}

# Example: Hierarchical contact groups (parent-child)
resource "netbox_contact_group" "parent" {
  name = "Engineering"
  slug = "engineering"
}

resource "netbox_contact_group" "child" {
  name        = "DevOps Team"
  slug        = "devops-team"
  parent      = netbox_contact_group.parent.name
  description = "DevOps and infrastructure automation team"
}

# Example: Contact group with tags
resource "netbox_contact_group" "with_tags" {
  name        = "Security Team"
  slug        = "security-team"
  description = "Information security team"
  tags        = ["critical"]
}
