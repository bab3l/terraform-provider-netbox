resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_location" "test" {
  name        = "Test Location"
  slug        = "test-location"
  site        = netbox_site.test.slug
  description = "Server room location"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "floor"
      value = "2"
    },
    {
      name  = "room_number"
      value = "SR-201"
    },
    {
      name  = "access_level"
      value = "restricted"
    }
  ]

  tags = [
    "server-room",
    "restricted-access"
  ]
}
