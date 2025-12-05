# Example: Basic contact role
resource "netbox_contact_role" "technical" {
  name = "Technical"
  slug = "technical"
}

# Example: Contact role with description
resource "netbox_contact_role" "administrative" {
  name        = "Administrative"
  slug        = "administrative"
  description = "Administrative and management contacts"
}

# Example: Contact role for billing
resource "netbox_contact_role" "billing" {
  name        = "Billing"
  slug        = "billing"
  description = "Billing and finance contacts"
}

# Example: Contact role for emergency
resource "netbox_contact_role" "emergency" {
  name        = "Emergency"
  slug        = "emergency"
  description = "Emergency and on-call contacts"
}

# Example: Contact role with tags
resource "netbox_contact_role" "with_tags" {
  name        = "Support"
  slug        = "support"
  description = "Customer support contacts"
  tags {
    name = "customer-facing"
    slug = "customer-facing"
  }
}
