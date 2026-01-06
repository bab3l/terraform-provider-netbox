//go:build customfields
// +build customfields

// Package resources_acceptance_tests_customfields contains acceptance tests
// that create/delete custom fields. These tests must run serially to avoid
// conflicts since custom fields are global per content type in NetBox.
//
// Custom fields in NetBox are defined globally for each object type (e.g.,
// dcim.device, ipam.aggregate). When multiple tests run in parallel and
// attempt to create or delete custom fields for the same object type, they
// can cause race conditions and database deadlocks.
//
// To run these tests:
//
//	go test -tags=customfields ./internal/resources_acceptance_tests_customfields/... -v
//
// Or using the Makefile:
//
//	make test-acceptance-customfields
//
// These tests are excluded from normal test runs to:
//  1. Speed up development cycles (saves 60-90 minutes)
//  2. Prevent parallel execution conflicts
//  3. Allow CI to run them separately in serial mode
package resources_acceptance_tests_customfields
