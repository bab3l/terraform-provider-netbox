package resources_acceptance_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationGroupResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	const testDescription = "Test Description"
	const testPassword = "tfacc-Password123!"

	// Ensure env is present before doing any API calls.
	testutil.TestAccPreCheck(t)

	serverURL := strings.TrimRight(os.Getenv("NETBOX_SERVER_URL"), "/")
	apiToken := os.Getenv("NETBOX_API_TOKEN")
	if serverURL == "" || apiToken == "" {
		t.Fatal("NETBOX_SERVER_URL and NETBOX_API_TOKEN must be set for acceptance tests")
	}

	ctx := context.Background()
	api := netboxRawAPI{serverURL: serverURL, apiToken: apiToken}

	groupName := testutil.RandomName("tf-test-group")
	groupID, err := api.createGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create NetBox auth group for acceptance test: %v", err)
	}

	userName := testutil.RandomName("tf-test-user")
	userID, err := api.createUser(ctx, userName, testPassword)
	if err != nil {
		_ = api.deleteGroup(ctx, groupID)
		t.Fatalf("failed to create NetBox user for acceptance test: %v", err)
	}

	// Ensure cleanup even if Terraform test fails.
	t.Cleanup(func() {
		_ = api.deleteUser(ctx, userID)
		_ = api.deleteGroup(ctx, groupID)
	})

	name := testutil.RandomName("notification-group-remove-ext")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterNotificationGroupCleanup(name)

	resourceRef := "netbox_notification_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationGroupResourceConfig_withDescriptionAndUsersGroups(name, testDescription, groupID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", name),
					resource.TestCheckResourceAttr(resourceRef, "description", testDescription),
					resource.TestCheckResourceAttr(resourceRef, "group_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceRef, "group_ids.*", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr(resourceRef, "user_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceRef, "user_ids.*", fmt.Sprintf("%d", userID)),
				),
			},
			{
				Config: testAccNotificationGroupResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", name),
					resource.TestCheckNoResourceAttr(resourceRef, "description"),
					resource.TestCheckNoResourceAttr(resourceRef, "group_ids.#"),
					resource.TestCheckNoResourceAttr(resourceRef, "user_ids.#"),
				),
			},
			{
				Config: testAccNotificationGroupResourceConfig_withDescriptionAndUsersGroups(name, testDescription, groupID, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRef, "name", name),
					resource.TestCheckResourceAttr(resourceRef, "description", testDescription),
					resource.TestCheckResourceAttr(resourceRef, "group_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceRef, "group_ids.*", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr(resourceRef, "user_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceRef, "user_ids.*", fmt.Sprintf("%d", userID)),
				),
			},
		},
	})
}

type netboxRawAPI struct {
	serverURL string
	apiToken  string
}

func (api netboxRawAPI) doJSON(ctx context.Context, method, path string, payload any) (map[string]interface{}, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, api.serverURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+api.apiToken)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("netbox %s %s returned %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if len(respBody) == 0 {
		return map[string]interface{}{}, nil
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func (api netboxRawAPI) createGroup(ctx context.Context, name string) (int32, error) {
	resp, err := api.doJSON(ctx, http.MethodPost, "/api/users/groups/", map[string]any{
		"name": name,
	})
	if err != nil {
		return 0, err
	}
	id, ok := resp["id"].(float64)
	if !ok {
		return 0, errors.New("group create response missing id")
	}
	return int32(id), nil // #nosec G115 -- NetBox IDs fit in int32
}

func (api netboxRawAPI) deleteGroup(ctx context.Context, id int32) error {
	_, err := api.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/api/users/groups/%d/", id), nil)
	return err
}

func (api netboxRawAPI) createUser(ctx context.Context, username, password string) (int32, error) {
	resp, err := api.doJSON(ctx, http.MethodPost, "/api/users/users/", map[string]any{
		"username":  username,
		"password":  password,
		"is_active": true,
	})
	if err != nil {
		return 0, err
	}
	id, ok := resp["id"].(float64)
	if !ok {
		return 0, errors.New("user create response missing id")
	}
	return int32(id), nil // #nosec G115 -- NetBox IDs fit in int32
}

func (api netboxRawAPI) deleteUser(ctx context.Context, id int32) error {
	_, err := api.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/api/users/users/%d/", id), nil)
	return err
}

func testAccNotificationGroupResourceConfig_withDescriptionAndUsersGroups(name, description string, groupID, userID int32) string {
	return fmt.Sprintf(`
resource "netbox_notification_group" "test" {
  name        = %[1]q
  description = %[2]q
  group_ids   = [%[3]d]
  user_ids    = [%[4]d]
}
`, name, description, groupID, userID)
}
