# Development Configuration

This file shows how to configure Terraform to use your locally built provider during development.

## 1. Create a `.terraformrc` file in your home directory

Windows: `%APPDATA%\terraform.rc`
Linux/macOS: `~/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "bab3l/netbox" = "C:\\GitRoot\\terraform-provider-netbox"
  }

  # For all other providers, install them directly as normal.
  direct {}
}
```

## 2. Build the provider

```bash
go build -o terraform-provider-netbox.exe .
```

## 3. Create a test Terraform configuration

Create a `test.tf` file:

```hcl
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  server_url = "https://your-netbox-instance.com"
  api_token  = "your-api-token"
}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}
```

## 4. Test the provider

```bash
terraform init
terraform plan
terraform apply
```

## Environment Variables for Testing

Set these environment variables for easier testing:

```bash
export NETBOX_SERVER_URL="https://your-netbox-instance.com"
export NETBOX_API_TOKEN="your-api-token"
export TF_LOG=DEBUG  # Enable debug logging
```

## Debugging

To debug the provider:

1. Build with debug flags:
   ```bash
   go build -gcflags="all=-N -l" -o terraform-provider-netbox.exe .
   ```

2. Run the provider in debug mode:
   ```bash
   .\terraform-provider-netbox.exe -debug
   ```

3. Use the provided TF_REATTACH_PROVIDERS environment variable in your Terraform commands.
