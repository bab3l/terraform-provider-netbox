# Example demonstrating validation features
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your-api-token-here"
}

# Valid site configuration
resource "netbox_site" "valid_site" {
  name        = "Production Datacenter"        # Valid: 1-100 chars
  slug        = "prod-dc"                      # Valid: lowercase, letters, numbers, hyphens
  status      = "active"                       # Valid: one of the allowed values
  description = "Primary production facility"  # Valid: under 200 chars
  facility    = "Building A"                   # Valid: under 50 chars
  comments    = "24/7 operations center"       # Valid: under 1000 chars

  tags = [
    {
      name = "production"     # Valid: 1-100 chars
      slug = "production"     # Valid: slug format
    },
    {
      name = "critical"
      slug = "critical"
    }
  ]

  custom_fields = [
    {
      name  = "cost_center"   # Valid: starts with letter, alphanumeric + underscore
      type  = "text"          # Valid: supported type
      value = "IT-INFRA-001"  # Valid: under 1000 chars
    },
    {
      name  = "rack_count"
      type  = "integer"
      value = "42"            # Valid: integer as string
    },
    {
      name  = "has_generator"
      type  = "boolean"
      value = "true"          # Valid: boolean as string
    },
    {
      name  = "metadata"
      type  = "json"
      value = "{\"contact\": \"admin@example.com\"}"  # Valid: JSON string
    }
  ]
}

# Examples that would trigger validation errors:

# resource "netbox_site" "invalid_slug" {
#   name = "Test Site"
#   slug = "Test Site!"  # ERROR: Contains spaces and special characters
#   # Error: Slug 'Test Site!' contains invalid character ' '. Only lowercase letters, numbers, hyphens, and underscores are allowed.
# }

# resource "netbox_site" "invalid_status" {
#   name   = "Test Site"
#   slug   = "test-site"
#   status = "unknown"    # ERROR: Not in allowed list
#   # Error: Attribute status value must be one of: ["planned" "staging" "active" "decommissioning" "retired"], got: "unknown"
# }

# resource "netbox_site" "name_too_long" {
#   name = "This is an extremely long site name that exceeds the maximum allowed length of 100 characters and will cause a validation error"  # ERROR: > 100 chars
#   slug = "long-name"
#   # Error: Attribute name string length must be between 1 and 100, got: 150
# }

# resource "netbox_site" "invalid_custom_field" {
#   name = "Test Site"
#   slug = "test-site"
#   
#   custom_fields = [
#     {
#       name  = "123invalid"   # ERROR: Starts with number
#       type  = "text"
#       value = "test"
#       # Error: Attribute custom_fields[0].name value must start with a letter and contain only letters, numbers, and underscores, got: "123invalid"
#     },
#     {
#       name  = "test_field"
#       type  = "invalid_type"  # ERROR: Not in allowed list
#       value = "test"
#       # Error: Attribute custom_fields[1].type value must be one of: ["text" "longtext" "integer" "boolean" "date" "url" "json" "select" "multiselect" "object" "multiobject" "multiple" "selection"], got: "invalid_type"
#     }
#   ]
# }

# resource "netbox_site" "invalid_tag" {
#   name = "Test Site"
#   slug = "test-site"
#   
#   tags = [
#     {
#       name = "production"
#       slug = "_invalid-slug_"  # ERROR: Starts and ends with underscore
#       # Error: Slug '_invalid-slug_' cannot start or end with hyphens or underscores
#     }
#   ]
# }
