# Terraform Provider Netbox Documentation Generation

## Canonical commands

- Regenerate examples and docs: `make docs`
- Check for docs/example drift without keeping changes: `make docs-check`

The `make docs` target delegates to `go generate ./...`, which runs the commands embedded in [main.go](../main.go):

1. `terraform fmt -recursive ./examples/`
2. `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs`

## What is generated

- Formatted Terraform example files under `examples/`
- Generated provider, resource, and data source documentation under `docs/`

Custom rendering is driven by templates in `templates/`.

## Contributor workflow

Rerun `make docs` whenever you change:

- resource schemas
- data source schemas
- provider schema or descriptions
- documentation templates under `templates/`
- example Terraform snippets that should stay formatted

Before opening a PR, run `make docs-check` to confirm there is no generated-doc drift.

## CI and release policy

- CI verifies generated docs and example formatting with `make docs-check`.
- Release automation validates docs/example cleanliness before publishing artifacts.
- Release tooling must not mutate repository state during the release itself.

## Troubleshooting

- If `make docs` fails on `terraform fmt`, ensure Terraform is installed and on `PATH`.
- If generated docs do not reflect schema changes, rerun `make docs` after verifying the schema code builds.
- If template changes are ignored, confirm the files are under `templates/` and rerun generation.

## References

- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs)
- [Terraform Provider Documentation Guide](https://developer.hashicorp.com/terraform/plugin/docs)
