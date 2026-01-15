package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactResource_basic(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(randomName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResource(randomName, "john.doe@example.com", "+1-555-0100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_contact.test", "email", "john.doe@example.com"),
					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_contact.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactResource(randomName, "jane.doe@example.com", "+1-555-0200"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_contact.test", "email", "jane.doe@example.com"),
					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0200"),
				),
			},
		},
	})
}

func TestAccContactResource_full(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(randomName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResourceFull(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttr("netbox_contact.test", "title", "Network Engineer"),
					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),
					resource.TestCheckResourceAttr("netbox_contact.test", "email", "engineer@example.com"),
					resource.TestCheckResourceAttr("netbox_contact.test", "address", "123 Main Street, City, Country"),
					resource.TestCheckResourceAttr("netbox_contact.test", "link", "https://example.com/profile"),
					resource.TestCheckResourceAttr("netbox_contact.test", "description", "Test contact description"),
					resource.TestCheckResourceAttr("netbox_contact.test", "comments", "Test contact comments"),
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},
		},
	})
}

func TestAccConsistency_Contact(t *testing.T) {
	t.Parallel()

	contactName := testutil.RandomName("contact")
	contactGroupName := testutil.RandomName("contactgroup")
	contactGroupSlug := testutil.RandomSlug("contactgroup")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(contactGroupSlug)
	cleanup.RegisterContactCleanup(contactName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "group", contactGroupName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug),
			},
		},
	})
}

func TestAccConsistency_Contact_LiteralNames(t *testing.T) {
	t.Parallel()

	contactName := testutil.RandomName("contact")
	contactGroupName := testutil.RandomName("contactgroup")
	contactGroupSlug := testutil.RandomSlug("contactgroup")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(contactGroupSlug)
	cleanup.RegisterContactCleanup(contactName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "group", contactGroupName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug),
			},
		},
	})
}

func TestAccContactResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	contactName := testutil.RandomName("tf-test-contact-update")
	updatedName := testutil.RandomName("tf-test-contact-updated")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("contact-upd"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResource(contactName, contactEmail, "+1-555-0100"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
				),
			},
			{
				Config: testAccContactResource(updatedName, contactEmail, "+1-555-0100"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccContactResource_IDPreservation(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)

	contactName := testutil.RandomName("tf-test-contact-id")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("contact-id"))
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResource(contactName, contactEmail, "+1-555-0100"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "email", contactEmail),
					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),
				),
			},
		},
	})
}

func testAccContactResource(name, email, phone string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name  = %[1]q
  email = %[2]q
  phone = %[3]q
}
`, name, email, phone)
}

func testAccContactResourceFull(name string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name        = %[1]q
  title       = "Network Engineer"
  phone       = "+1-555-0100"
  email       = "engineer@example.com"
  address     = "123 Main Street, City, Country"
  link        = "https://example.com/profile"
  description = "Test contact description"
  comments    = "Test contact comments"
}
`, name)
}

func testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_contact" "test" {
  name = "%[1]s"
  group = netbox_contact_group.test.name
}
`, contactName, contactGroupName, contactGroupSlug)
}

func testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_contact" "test" {
  name = "%[1]s"
  group = "%[2]s"
  depends_on = [netbox_contact_group.test]
}
`, contactName, contactGroupName, contactGroupSlug)
}

func TestAccContactResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	contactName := testutil.RandomName("tf-test-contact-del")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("contact-del"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResource(contactName, contactEmail, "+1-555-0100"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("netbox_contact.test", "email", contactEmail),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyContactsList(context.Background()).Email([]string{contactEmail}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact: %v", err)
					}
					t.Logf("Successfully externally deleted contact with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccContactResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	contactName := testutil.RandomName("tf-test-contact-rem")
	contactGroupName := testutil.RandomName("tf-test-cg-rem")
	contactGroupSlug := testutil.RandomSlug("tf-test-cg-rem")
	email := "test@example.com"
	phone := "+1234567890"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactGroupCleanup(contactGroupSlug)
	cleanup.RegisterContactCleanup(contactName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckResourceAttrSet("netbox_contact.test", "group"),
				),
			},
			{
				Config: testAccContactResource(contactName, email, phone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),
					resource.TestCheckNoResourceAttr("netbox_contact.test", "group"),
				),
			},
		},
	})
}

func TestAccContactResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_contact",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact" "test" {
  # name missing
  title = "Test Title"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
