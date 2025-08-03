# Site Data Source Examples

This directory contains examples of how to use the `netbox_site` data source to read existing site information from Netbox.

## Examples

### 1. Basic Usage (`main.tf`)
Demonstrates basic data source usage:
- Reading a site by ID
- Reading a site by slug
- Outputting site information

### 2. Resource Integration (`with_resource.tf`)
Shows how to use data sources alongside resources:
- Creating a new site with a resource
- Reading the created site with data sources
- Comparing resource vs data source outputs

### 3. Site Group Integration (`with_site_group.tf`)
Advanced example combining multiple resource types:
- Creating site groups and sites
- Using data sources for existing infrastructure
- Conditional logic based on site properties
- Complex output structures

## Key Features

The `netbox_site` data source supports:

- **Flexible Identification**: Find sites by either ID or slug
- **Complete Site Information**: Access all site attributes including:
  - Basic info (name, slug, status)
  - Organizational data (region, group, tenant)
  - Descriptive fields (description, comments, facility)
  - Metadata (tags, custom fields)

## Common Use Cases

1. **Reference Existing Infrastructure**: Use data sources to incorporate existing Netbox sites into your Terraform configurations
2. **Cross-Reference Resources**: Link new resources to existing sites
3. **Conditional Logic**: Make decisions based on site properties
4. **Data Validation**: Verify site configurations and statuses
5. **Reporting**: Generate outputs with site information for documentation

## Usage Tips

- Prefer using `slug` for human-readable configurations
- Use `id` when you need to reference sites by their exact Netbox ID
- Data sources are read-only - use resources to modify sites
- Combine data sources with locals for complex conditional logic
- Use outputs to expose site information for external consumption

## Authentication

All examples require proper Netbox authentication:

```hcl
provider "netbox" {
  server_url = "https://your-netbox-instance.com"
  api_token  = "your-api-token"
}
```

Alternatively, use environment variables:
- `NETBOX_SERVER_URL`
- `NETBOX_API_TOKEN`
