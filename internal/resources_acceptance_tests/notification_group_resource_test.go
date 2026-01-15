package resources_acceptance_tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNotificationGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-basic")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_notification_group.test", "description"),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "group_ids.#", "0"),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "user_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccNotificationGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-full")
	description := "Test notification group with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_full(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccNotificationGroupResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-update")
	updatedName := testutil.RandomName("tf-test-notifgroup-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)
	cleanup.RegisterNotificationGroupCleanup(updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_notification_group.test", "description"),
				),
			},
			{
				Config: testAccNotificationGroupResourceConfig_full(updatedName, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccNotificationGroupResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_notification_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccNotificationGroupResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccNotificationGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
				),
			},
		},
	})
}

func TestAccConsistency_NotificationGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
				),
			},
			{
				Config:   testAccNotificationGroupConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccNotificationGroupConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name = %q
}
`, name)
}

func TestAccNotificationGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-notifgroup-ext-del")
	var notificationGroupID string

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckNotificationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["netbox_notification_group.test"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						notificationGroupID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					if notificationGroupID == "" {
						t.Fatal("notification group ID not captured from previous step")
					}
					// Delete the notification group externally via API
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("failed to get API client: %s", err)
					}
					id, err := strconv.Atoi(notificationGroupID)
					if err != nil {
						t.Fatalf("failed to convert ID to int: %s", err)
					}
					//nolint:gosec // ID from NetBox API is always a valid positive integer
					_, err = client.ExtrasAPI.ExtrasNotificationGroupsDestroy(context.Background(), int32(id)).Execute()
					if err != nil {
						t.Fatalf("failed to delete notification group externally: %s", err)
					}
				},
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_notification_group.test", "id"),
				),
			},
		},
	})
}

// TestAccNotificationGroupResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccNotificationGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	const testDescription = "Test Description"

	name := testutil.RandomName("notification-group-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_withDescription(name, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_notification_group.test", "description", testDescription),
				),
			},
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_notification_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_notification_group.test", "description"),
				),
			},
		},
	})
}

func testAccNotificationGroupResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name = %[1]q
}
`, name)
}

func testAccNotificationGroupResourceConfig_full(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}

func testAccNotificationGroupResourceConfig_withDescription(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}

func TestAccNotificationGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_notification_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_notification_group" "test" {
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
