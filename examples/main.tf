terraform {
  required_providers {
    netbox = {
      source  = "bab3l/netbox"
      version = "~> 0.1.0"
    }
  }
}

provider "netbox" {
  server_url = "https://netbox.example.com"
  api_token  = "your-api-token-here" # Or set NETBOX_API_TOKEN environment variable
  insecure   = false
}

# Example site resource
resource "netbox_site" "example" {
  name        = "Example Site"
  slug        = "example-site"
  status      = "active"
  description = "An example site created with Terraform"
}
