# Terraform Provider for Netbox - Documentation

This directory contains comprehensive documentation for all resources and data sources available in the Terraform Netbox provider.

## Quick Navigation

### Provider Configuration
- [Provider](./index.md) - Provider configuration and authentication

### Resources

#### Organization & Location Management
- [netbox_site](./resources/site.md) - Physical locations (data centers, offices, facilities)
- [netbox_site_group](./resources/site_group.md) - Hierarchical organization of sites

#### Device & Infrastructure Management
*Coming soon*
- `netbox_device_type` - Device type definitions
- `netbox_device_role` - Device role classifications  
- `netbox_device` - Physical and virtual devices
- `netbox_rack` - Equipment racks
- `netbox_rack_role` - Rack role classifications

#### Network Management
*Coming soon*
- `netbox_ip_address` - IP address assignments
- `netbox_prefix` - IP prefixes and subnets
- `netbox_vlan` - VLAN definitions
- `netbox_vrf` - Virtual routing and forwarding instances

#### Circuit & Provider Management
*Coming soon*
- `netbox_circuit` - Network circuits
- `netbox_circuit_type` - Circuit type definitions
- `netbox_provider` - Service providers

#### Tenancy & Organization
*Coming soon*
- `netbox_tenant` - Multi-tenancy support
- `netbox_tenant_group` - Tenant organization
- `netbox_region` - Geographic regions

### Data Sources

*Data sources will be available in future releases to provide read-only access to Netbox data.*

#### Examples of Planned Data Sources
- `netbox_site` - Look up site information
- `netbox_device_type` - Query device type details
- `netbox_ip_address` - Find available IP addresses
- `netbox_prefix` - Search for prefixes

## Resource Features

All resources in this provider include:

### Common Features
- **Tags Support**: Organize and categorize resources
- **Custom Fields**: Store additional structured metadata
- **Validation**: Comprehensive field validation with helpful error messages
- **Import Support**: Import existing Netbox resources into Terraform
- **State Management**: Proper Terraform state handling

### Authentication & Configuration
- Environment variable support
- API token authentication
- TLS/SSL configuration options
- Connection validation

## Getting Started

1. **Provider Setup**: Start with the [provider configuration](./index.md)
2. **Basic Resources**: Begin with [sites](./resources/site.md) and [site groups](./resources/site_group.md)
3. **Examples**: Check the `../examples/` directory for complete configurations
4. **Testing**: See the testing guides for validation approaches

## Examples Directory Structure

```
examples/
├── provider/                    # Provider configuration examples
├── resources/
│   ├── netbox_site/            # Site resource examples
│   ├── netbox_site_group/      # Site group examples
│   └── ...                     # Additional resource examples
└── complete/                   # End-to-end configurations
```

## Validation & Error Handling

This provider includes comprehensive validation:

- **Field Validation**: Type checking, format validation, length limits
- **Business Logic**: Netbox-specific rules and constraints
- **Relationship Validation**: Parent-child and reference validation
- **Clear Error Messages**: Helpful diagnostics for troubleshooting

## Best Practices

### Resource Organization
- Use site groups for hierarchical organization
- Apply consistent tagging strategies
- Leverage custom fields for metadata
- Plan import strategies for existing infrastructure

### State Management
- Use remote state for team collaboration
- Implement proper backup strategies
- Consider state splitting for large infrastructures
- Use workspaces for environment separation

### Security
- Store API tokens securely (environment variables or secrets management)
- Use least-privilege API tokens
- Enable TLS verification in production
- Audit resource changes regularly

## Contributing

See the main [README](../README.md) for contribution guidelines. When adding new resources:

1. Follow the established patterns in existing resources
2. Include comprehensive validation
3. Add unit tests for all functionality
4. Generate documentation using `tfplugindocs`
5. Provide examples for common use cases

## Support

For issues and questions:
- Check existing documentation and examples
- Review validation error messages
- Test with the Netbox API directly
- Open issues with detailed reproduction steps
