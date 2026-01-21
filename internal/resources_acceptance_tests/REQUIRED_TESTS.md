# Required Acceptance Tests (Concise)

This file summarizes the **minimum acceptance tests** every resource must implement, plus the **helpers** to use.

## Required Tests (all resources)

- `TestAcc{Resource}Resource_basic` – create with required fields only.
- `TestAcc{Resource}Resource_full` – create with all optional fields.
- `TestAcc{Resource}Resource_update` – change existing resource.
- `TestAcc{Resource}Resource_import` – import verification; **must validate reference fields are numeric IDs**.
- `TestAcc{Resource}Resource_externalDeletion` – handles out-of-band deletion.
- `TestAcc{Resource}Resource_removeOptionalFields` – clear optional fields.

## Tag Tests (resources with tags)

Use **helpers only** (no custom tag tests):

- `TestAcc{Resource}Resource_tagLifecycle` – add/replace/remove tags.
- `TestAcc{Resource}Resource_tagOrderInvariance` – order does not drift.

## Helpers (internal/testutil)

Core:

- `RunImportTest()` / `RunSimpleImportTest()`
- `ReferenceFieldCheck()` – validates field is numeric ID after import
- `ValidateReferenceIDs()` – batch validates multiple reference fields
- `RunExternalDeletionTest()`
- `TestRemoveOptionalFields()`
- `RunMultiValidationErrorTest()`

Tags:

- `RunTagLifecycleTest()`
- `RunTagOrderTest()`

## Conventions (brief)

- Test names: `TestAcc{Resource}Resource_{testType}`.
- Config helpers: `testAcc{Resource}ResourceConfig_{variant}`.
- Always register cleanup for created resources.
