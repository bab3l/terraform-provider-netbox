# Example: Basic contact
resource "netbox_contact" "basic" {
  name  = "John Doe"
  email = "john.doe@example.com"
}

# Example: Contact with full details
resource "netbox_contact" "full" {
  name        = "Jane Smith"
  title       = "Network Engineer"
  phone       = "+1-555-0100"
  email       = "jane.smith@example.com"
  address     = "123 Main Street, City, Country"
  link        = "https://example.com/jsmith"
  description = "Primary network contact"
  comments    = "Available Mon-Fri, 9am-5pm"
}

# Example: Contact with tags
resource "netbox_contact" "with_tags" {
  name  = "Bob Wilson"
  email = "bob.wilson@example.com"

  tags = [
    {
      name = "primary"
      slug = "primary"
    },
    {
      name = "on-call"
      slug = "on-call"
    }
  ]
}
