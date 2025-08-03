# Example Terraform configuration demonstrating tags and custom fields
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

# Example site with tags and custom fields
resource "netbox_site" "main_datacenter" {
  name        = "Main Data Center"
  slug        = "main-dc"
  status      = "active"
  description = "Primary data center facility"
  facility    = "Building A - Floor 2"
  comments    = "This is our primary datacenter with redundant power and cooling"

  # Tags for categorization
  tags = [
    {
      name = "production"
      slug = "production"
    },
    {
      name = "datacenter"
      slug = "datacenter"
    },
    {
      name = "critical"
      slug = "critical"
    }
  ]

  # Custom fields for additional metadata
  custom_fields = [
    {
      name  = "cost_center"
      type  = "text"
      value = "IT-Infrastructure"
    },
    {
      name  = "power_capacity"
      type  = "integer"
      value = "500"
    },
    {
      name  = "backup_generator"
      type  = "boolean"
      value = "true"
    },
    {
      name  = "cooling_type"
      type  = "select"
      value = "chilled-water"
    },
    {
      name  = "certifications"
      type  = "multiselect"
      value = "SOC2,ISO27001,PCI-DSS"
    },
    {
      name  = "additional_info"
      type  = "json"
      value = "{\"emergency_contact\": \"+1-555-0123\", \"access_hours\": \"24/7\"}"
    }
  ]
}

# Example site with minimal configuration
resource "netbox_site" "branch_office" {
  name = "Branch Office NYC"
  slug = "branch-nyc"

  tags = [
    {
      name = "branch"
      slug = "branch"
    }
  ]

  custom_fields = [
    {
      name  = "cost_center"
      type  = "text"
      value = "Sales-East"
    }
  ]
}
