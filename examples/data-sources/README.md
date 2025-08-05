# Comprehensive Data Sources Example

This example demonstrates the full capabilities of both `netbox_site` and `netbox_site_group` data sources working together in a realistic infrastructure scenario.

## What This Example Demonstrates

### 1. Hierarchical Site Group Structure
- **Global Group**: Top-level organizational group
- **Regional Group**: Geographic organization (North America)
- **Country Group**: Country-specific organization (United States)

### 2. Site Management
- **Corporate Headquarters**: Administrative facility
- **Primary Datacenter**: Production infrastructure facility

### 3. Data Source Capabilities
- Reading resources by both ID and slug
- Accessing all resource attributes
- Understanding hierarchical relationships
- Integrating with existing infrastructure

### 4. Practical Use Cases
- **Hierarchy Validation**: Ensuring parent-child relationships are correct
- **Resource Discovery**: Finding and referencing existing infrastructure
- **Operational Summaries**: Generating useful reports for operations teams
- **Conditional Logic**: Making decisions based on resource properties

## Key Features Showcased

### Data Source Flexibility
```hcl
# Read by ID
data "netbox_site_group" "global_data" {
  id = netbox_site_group.global.id
}

# Read by slug
data "netbox_site_group" "regional_data" {
  slug = netbox_site_group.regional.slug
}
```

### Hierarchical Relationships
```hcl
# Validate parent-child relationships
locals {
  regional_parent_correct = data.netbox_site_group.regional_data.parent == data.netbox_site_group.global_data.name
}
```

### Cross-Resource Validation
```hcl
# Ensure sites belong to correct groups
locals {
  hq_group_correct = data.netbox_site.hq_data.group == data.netbox_site_group.country_data.name
}
```

### Operational Intelligence
```hcl
# Generate lists of resources meeting criteria
locals {
  production_sites = [
    for site_key, site_data in local.all_sites : {
      name = site_data.name
      id   = site_data.id
    }
    if site_data.status == "active"
  ]
}
```

## Expected Outputs

### Infrastructure Overview
Complete hierarchical view of:
- All site groups with their relationships
- All sites with their group assignments
- Existing resources for reference

### Hierarchy Validation
Automated validation showing:
- Whether hierarchy is correctly configured
- Individual relationship checks
- Overall validation status

### Operational Summary
Practical information including:
- Group IDs for creating new resources
- Lists of active sites
- Site counts by status
- Usage recommendations

## Best Practices Demonstrated

1. **Use Both ID and Slug**: Show when to use each lookup method
2. **Validate Relationships**: Automated checking of hierarchical structure
3. **Combine Resources**: Show site groups and sites working together
4. **Error Prevention**: Validate configurations before using them
5. **Operational Focus**: Generate outputs useful for day-to-day operations

## Running the Example

1. Ensure you have existing resources with slugs:
   - `legacy-infrastructure` (site group)
   - `backup-facility` (site)

2. Configure your Netbox provider:
   ```hcl
   provider "netbox" {
     server_url = "https://your-netbox-instance.com"
     api_token  = "your-api-token"
   }
   ```

3. Run Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Extending the Example

You can extend this example by:

1. **Adding More Hierarchy Levels**: Create state/province level groups
2. **Adding More Sites**: Create additional sites in different groups
3. **Adding Validation Rules**: Create more sophisticated validation logic
4. **Adding Conditional Resources**: Use the data sources to conditionally create resources
5. **Adding Integration**: Show integration with other Netbox resources like devices or racks

## Use Cases for This Pattern

- **Infrastructure Discovery**: Map existing Netbox infrastructure
- **Migration Planning**: Understand current state before changes
- **Compliance Checking**: Validate infrastructure meets organizational standards
- **Operational Reporting**: Generate infrastructure summaries
- **Resource Planning**: Understand where to place new resources

This example serves as a foundation for more complex infrastructure management scenarios using the Netbox Terraform provider.
