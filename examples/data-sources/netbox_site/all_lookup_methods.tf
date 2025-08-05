# Site Data Source - All Lookup Methods Example

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

# Create a test site for demonstration
resource "netbox_site" "demo_site" {
  name        = "Demo Datacenter Site"
  slug        = "demo-datacenter"
  status      = "active"
  description = "Demonstration site for testing all lookup methods"
  facility    = "Demo Building"

  tags = [
    {
      name = "demo"
      slug = "demo"
    }
  ]
}

# Method 1: Lookup by ID (most efficient)
data "netbox_site" "by_id" {
  id = netbox_site.demo_site.id
}

# Method 2: Lookup by slug (human-readable and unique)
data "netbox_site" "by_slug" {
  slug = netbox_site.demo_site.slug
}

# Method 3: Lookup by name (most human-readable)
data "netbox_site" "by_name" {
  name = netbox_site.demo_site.name
}

# Lookup existing sites using different methods
data "netbox_site" "existing_by_id" {
  id = "1"
}

data "netbox_site" "existing_by_slug" {
  slug = "headquarters"
}

data "netbox_site" "existing_by_name" {
  name = "Corporate Headquarters"
}

# Verify all methods return the same data for our demo site
locals {
  demo_site_results = {
    by_id = {
      id   = data.netbox_site.by_id.id
      name = data.netbox_site.by_id.name
      slug = data.netbox_site.by_id.slug
    }
    by_slug = {
      id   = data.netbox_site.by_slug.id
      name = data.netbox_site.by_slug.name
      slug = data.netbox_site.by_slug.slug
    }
    by_name = {
      id   = data.netbox_site.by_name.id
      name = data.netbox_site.by_name.name
      slug = data.netbox_site.by_name.slug
    }
  }

  # Validate that all methods return the same site
  results_match = (
    local.demo_site_results.by_id.id == local.demo_site_results.by_slug.id &&
    local.demo_site_results.by_slug.id == local.demo_site_results.by_name.id
  )
}

output "lookup_methods_demo" {
  description = "Demonstration of all three lookup methods"
  value = {
    created_site = {
      id   = netbox_site.demo_site.id
      name = netbox_site.demo_site.name
      slug = netbox_site.demo_site.slug
    }

    lookup_results = local.demo_site_results

    validation = {
      all_methods_return_same_site = local.results_match
      message                      = local.results_match ? "✅ All lookup methods work correctly" : "❌ Lookup methods returned different results"
    }
  }
}

output "existing_sites" {
  description = "Information about existing sites using different lookup methods"
  value = {
    by_id = {
      id          = data.netbox_site.existing_by_id.id
      name        = data.netbox_site.existing_by_id.name
      slug        = data.netbox_site.existing_by_id.slug
      status      = data.netbox_site.existing_by_id.status
      description = data.netbox_site.existing_by_id.description
    }

    by_slug = {
      id          = data.netbox_site.existing_by_slug.id
      name        = data.netbox_site.existing_by_slug.name
      slug        = data.netbox_site.existing_by_slug.slug
      status      = data.netbox_site.existing_by_slug.status
      description = data.netbox_site.existing_by_slug.description
    }

    by_name = {
      id          = data.netbox_site.existing_by_name.id
      name        = data.netbox_site.existing_by_name.name
      slug        = data.netbox_site.existing_by_name.slug
      status      = data.netbox_site.existing_by_name.status
      description = data.netbox_site.existing_by_name.description
    }
  }
}

# Best practices for choosing lookup methods
locals {
  lookup_recommendations = {
    use_id_when = [
      "You know the exact Netbox ID",
      "Performance is critical",
      "Working with API responses that include IDs",
      "Building programmatic integrations"
    ]

    use_slug_when = [
      "You want human-readable configurations",
      "You control the slug naming convention",
      "You need guaranteed uniqueness",
      "Building reusable modules"
    ]

    use_name_when = [
      "You only know the display name",
      "Working with user-provided input",
      "Names are guaranteed unique in your environment",
      "Maximum readability is important"
    ]

    warnings = {
      name_lookup = "Names may not be unique in Netbox - use with caution in production"
      efficiency  = "ID lookup is most efficient, followed by slug, then name"
    }
  }
}

output "lookup_best_practices" {
  description = "Guidance on choosing the right lookup method"
  value       = local.lookup_recommendations
}
