package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestJournalEntryResource(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	if r == nil {
		t.Fatal("Expected non-nil resource")
	}
}

func TestJournalEntryResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", schemaResponse.Diagnostics)
	}
}

func TestJournalEntryResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_journal_entry")
}

func TestJournalEntryResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewJournalEntryResource()
	testutil.ValidateResourceConfigure(t, r)
}
