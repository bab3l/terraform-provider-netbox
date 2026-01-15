package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderAccountResource_basic(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	accountID := testutil.RandomName("acct")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
				),
			},
			{
				ResourceName:      "netbox_provider_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderAccountResource_full(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	accountID := testutil.RandomName("acct")
	accountName := testutil.RandomName("tf-test-acct")
	description := testutil.RandomName("description")
	updatedDescription := "Updated provider account description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "name", accountName),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", description),
				),
			},
			{
				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccProviderAccountResource_IDPreservation(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-id")
	providerSlug := testutil.RandomSlug("tf-test-provider-id")
	accountID := testutil.RandomName("tf-test-acct-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "circuit_provider"),
				),
			},
		},
	})
}

func TestAccProviderAccountResource_update(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-upd")
	providerSlug := testutil.RandomSlug("tf-test-provider-upd")
	accountID := testutil.RandomName("tf-test-acct-upd")
	accountName := testutil.RandomName("Account Name")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
				),
			},
			{
				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "name", accountName),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccProviderAccountResource_import(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-imp")
	providerSlug := testutil.RandomSlug("tf-test-provider-imp")
	accountID := testutil.RandomName("tf-test-acct-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
			},
			{
				ResourceName:      "netbox_provider_account.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account          = %q
}
`, providerName, providerSlug, accountID)
}

func testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account          = %q
  name             = %q
  description      = %q
}
`, providerName, providerSlug, accountID, accountName, description)
}

func TestAccConsistency_ProviderAccount_LiteralNames(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-lit")
	providerSlug := testutil.RandomSlug("tf-test-provider-lit")
	accountID := testutil.RandomName("acct-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
				),
			},
			{
				Config:   testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
				),
			},
		},
	})
}

func TestAccProviderAccountResource_externalDeletion(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-ext-del")
	providerSlug := testutil.RandomSlug("provider-ext-del")
	accountID := testutil.RandomName("tf-test-account-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List provider accounts filtered by account ID
					items, _, err := client.CircuitsAPI.CircuitsProviderAccountsList(context.Background()).AccountIc([]string{accountID}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find provider account for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsProviderAccountsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete provider account: %v", err)
					}
					t.Logf("Successfully externally deleted provider account with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProviderAccountResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-rem")
	providerSlug := testutil.RandomSlug("tf-test-provider-rem")
	accountID := testutil.RandomName("acct")
	accountName := "Account Name"
	const testDescription = "Description"
	const testComments = "Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourceConfig_fullWithComments(providerName, providerSlug, accountID, accountName, testDescription, testComments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_account.test", "name", accountName),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "comments", testComments),
				),
			},
			{
				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_provider_account.test", "name"),
					resource.TestCheckNoResourceAttr("netbox_provider_account.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_provider_account.test", "comments"),
				),
			},
		},
	})
}

func testAccProviderAccountResourceConfig_fullWithComments(providerName, providerSlug, accountID, accountName, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_account" "test" {
  account          = %q
  circuit_provider = netbox_provider.test.id
  name             = %q
  description      = %q
  comments         = %q
}
`, providerName, providerSlug, accountID, accountName, description, comments)
}

func TestAccProviderAccountResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_provider_account",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_circuit_provider": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider_account" "test" {
  # circuit_provider missing
  account = "12345"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_account": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  # account missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
