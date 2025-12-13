# Contributing to terraform-provider-netbox

Thank you for your interest in contributing!

## License and Contributor Terms
- This project is licensed under the Mozilla Public License 2.0 (MPL-2.0).
- By contributing, you agree that your contributions are licensed under MPL-2.0.
- Include the MPL-2.0 notice where appropriate for new files.

## Development Prerequisites
- Go 1.24.x
- Terraform 1.0+
- NetBox v4.1.x (for acceptance tests)
- Docker (for local acceptance tests)

## Workflow
- Fork the repo and create a feature branch.
- Run `make dev` or:
  - `go fmt ./...`
  - `go vet ./...`
  - `go test ./internal/... -v`
- For acceptance tests, use `scripts/run-acceptance-tests.ps1` (Windows) or `.sh`.
- Ensure docs are generated: `c:\GitRoot\terraform-plugin-docs\tfplugindocs.exe generate --provider-dir=. --rendered-website-dir=docs`.

## Commit and PR Guidelines
- Use clear, descriptive commit messages.
- Include tests for new functionality (unit + acceptance when applicable).
- Update examples and docs as needed.
- Link related issues in PR description.

## Code Style
- Use terraform-plugin-framework (not terraform-plugin-sdk).
- Prefer diagnostics for error handling.
- Use `tflog` for structured logging.
- Follow patterns defined in `RESOURCE-IMPLEMENTATION-TODO.md`.

## Security
- Do not include secrets or tokens in code or tests.
- Report vulnerabilities via the process in `SECURITY.md`.

## Questions
Open a GitHub discussion or issue for design questions before large changes.
