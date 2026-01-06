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
- For acceptance tests, see Testing section below.
- Ensure docs are generated: `c:\GitRoot\terraform-plugin-docs\tfplugindocs.exe generate --provider-dir=. --rendered-website-dir=docs`.

## Testing

### Unit Tests (Fast)
```bash
make test-fast
```
Runs unit tests only (~1-2 minutes).

### Acceptance Tests (Requires NetBox)
Set environment variables:
```bash
export NETBOX_SERVER_URL="http://localhost:8000"
export NETBOX_API_TOKEN="0123456789abcdef0123456789abcdef01234567"
```

#### Option 1: Parallel Tests Only (Recommended for Development)
```bash
make test-acceptance
```
Runs ~150 parallel-safe tests in 30-40 minutes. **This is the fastest option for development cycles.**

#### Option 2: Custom Field Tests Only (Serial Execution)
```bash
make test-acceptance-customfields
```
Runs 41 custom field tests serially in 60-90 minutes. These tests are separated because custom fields are global per content type in NetBox and cause conflicts when run in parallel.

#### Option 3: Full Test Suite
```bash
make test-acceptance-all
```
Runs all acceptance tests (2-3 hours total). Use this before submitting PRs.

### Why Two Test Packages?
Custom field tests are in `internal/resources_acceptance_tests_customfields/` with build tag `customfields`. This separation:
- Speeds up normal test runs (saves 60-90 minutes)
- Prevents database deadlocks from parallel custom field creation/deletion
- Allows CI to run parallel tests concurrently and serial tests sequentially

### Legacy Command
```bash
make testacc
```
Still available - runs all tests with `TestAcc` prefix.

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
