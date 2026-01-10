//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
