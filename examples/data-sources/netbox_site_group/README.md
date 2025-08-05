# Site Group Data Source Examples

This directory contains examples of how to use the `netbox_site_group` data source to read existing site group information from Netbox.

## Examples

### 1. Basic Usage (`main.tf`)
Demonstrates basic data source usage:
- Reading a site group by ID
- Reading a site group by slug
- Outputting site group information

### 2. Hierarchical Organization (`with_hierarchy.tf`)
Shows advanced hierarchical site group management:
- Creating parent and child site groups
- Reading site groups with data sources
- Understanding parent-child relationships
- Conditional logic for hierarchy validation
- Complex nested outputs

### 3. Site Integration (`with_sites.tf`)
Comprehensive example combining site groups and sites:
- Creating site groups and associated sites
- Using data sources for both resource types
- Relationship validation between sites and groups
- Tag consistency checking
- Using data sources for resource dependencies

## Key Features

The `netbox_site_group` data source supports:

- **Flexible Identification**: Find site groups by either ID or slug
- **Complete Site Group Information**: Access all site group attributes including:
  - Basic info (name, slug, description)
  - Hierarchical data (parent relationships)
  - Metadata (tags, custom fields)
- **Hierarchical Support**: Understand parent-child relationships in site group organization

## Hierarchical Organization

Site groups in Netbox support hierarchical organization:

```hcl
# Parent group
data "netbox_site_group" "regional" {
  slug = "north-america"
}

# Child group
data "netbox_site_group" "country" {
  slug = "united-states"
}

# The child group's parent field will reference the parent group
locals {
  is_child_of_regional = data.netbox_site_group.country.parent == data.netbox_site_group.regional.name
}
```

## Common Use Cases

1. **Hierarchical Organization**: Organize sites into logical groups with parent-child relationships
2. **Reference Existing Groups**: Use data sources to incorporate existing site groups into configurations
3. **Site Assignment**: Use site group data to assign sites to appropriate groups
4. **Validation**: Verify site group configurations and relationships
5. **Conditional Logic**: Make decisions based on site group properties and hierarchy
6. **Reporting**: Generate outputs with hierarchical site group information

## Integration with Sites

Site groups and sites work together for comprehensive organization:

```hcl
# Read a site group
data "netbox_site_group" "production" {
  slug = "production-sites"
}

# Create a site in that group
resource "netbox_site" "new_datacenter" {
  name  = "New Datacenter"
  slug  = "new-dc"
  group = data.netbox_site_group.production.id
}
```

## Usage Tips

- Prefer using `slug` for human-readable configurations
- Use `id` when you need to reference site groups by their exact Netbox ID
- Data sources are read-only - use resources to modify site groups
- Use the `parent` field to understand hierarchical relationships
- Combine data sources with locals for complex conditional logic
- Use outputs to expose site group information for external consumption

## Validation Patterns

The examples demonstrate several validation patterns:

1. **Hierarchy Validation**: Check if parent-child relationships are correct
2. **Tag Consistency**: Ensure related resources have consistent tagging
3. **Relationship Verification**: Validate that sites belong to expected groups

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
