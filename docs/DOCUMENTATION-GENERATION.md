# Terraform Provider Netbox Documentation Generation

## Overview
This project uses the [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) tool to generate documentation for all resources and data sources. The output is placed in the `docs/` directory and is based on Go code schemas and custom templates in the `templates/` directory.

## Key Points
- **Tool:** Documentation is generated using the command:
  ```
  go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir=. --rendered-website-dir=docs
  ```
- **Templates:** Custom Markdown templates for resources and data sources are located in the `templates/` directory. These control the format and content of generated docs.
- **Output:** Generated documentation is written to the `docs/` directory, organized by resource and data source.
- **Schema:** The tool parses Go code schemas (attributes, descriptions, types) and renders them in Markdown, including nested attributes and read-only fields.
- **Customization:** To change documentation output, edit the templates in `templates/` and rerun the generation command.
- **Regeneration:** Always rerun the documentation command after changing resource/data source schemas or templates to keep docs up to date.

## Example Output
- Data source documentation (e.g., `docs/data-sources/site_group.md`) includes:
  - Title, description, and usage
  - Schema: Optional and Read-Only attributes
  - Nested schemas for complex attributes

## Troubleshooting
- If the output directory is incorrect, use the `--rendered-website-dir` flag (not `--output-dir`).
- If templates are not applied, verify the template directory and rerun the command.
- If documentation is missing new fields, ensure Go code schemas are updated and rerun the command.

## References
- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)
- [Terraform Provider Documentation Guide](https://developer.hashicorp.com/terraform/plugin/docs)

---
_Last updated: August 5, 2025_
