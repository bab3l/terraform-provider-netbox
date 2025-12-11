<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Terraform Provider for Netbox Development Instructions

This workspace contains a Terraform provider for Netbox using the modern terraform-plugin-framework.

## Key Technologies and Patterns:
- **Framework**: Use terraform-plugin-framework (NOT terraform-plugin-sdk)
- **API Client**: Integration with go-netbox wrapper (available locally at ../go-netbox during development)
- **Go Version**: Go 1.24.5
- **Testing**: Use terraform-plugin-testing framework for acceptance tests

## Code Organization:
- `internal/provider/`: Core provider implementation
- `internal/resources/`: Resource implementations (CRUD operations)
- `internal/datasources/`: Data source implementations (read-only)
- `examples/`: Terraform configuration examples
- `docs/`: Generated documentation

## Development Guidelines:
1. **Schema Definition**: Use terraform-plugin-framework schema types (types.String, types.Bool, etc.)
2. **Error Handling**: Use diagnostics for validation and runtime errors
3. **Logging**: Use tflog for structured logging
4. **Testing**: Write both unit tests and acceptance tests
5. **Documentation**: Generate docs using terraform-plugin-docs

## Resource/Data Source Pattern:
- Implement Configure() method to get client from provider
- Use structured models for Terraform state and API data
- Handle Create, Read, Update, Delete operations with proper error handling
- Use diagnostics for validation and error reporting

## Local Development:
- The go-netbox dependency uses a local replace directive pointing to ../go-netbox
- Run `make dev` for the full development cycle (format, vet, test, build)
- Use `make install` to install provider locally for testing with Terraform
