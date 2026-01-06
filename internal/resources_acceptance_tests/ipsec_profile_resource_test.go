package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECProfileResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "ah"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIPSECProfileResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnIpsecProfilesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IPSec profile for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnIpsecProfilesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IPSec profile: %v", err)
					}
					t.Logf("Successfully externally deleted IPSec profile with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccIPSecProfileResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-profile-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},
		},
	})
}
func testAccIPSECProfileResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_ipsec_profile" "test" {
  name           = %q
  mode           = "esp"
  ike_policy     = netbox_ike_policy.test.id
  ipsec_policy   = netbox_ipsec_policy.test.id
}
`, testAccIPSECProfileResourcePrereqs(name), name)
}

func testAccIPSECProfileResourceConfig_full(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_ipsec_profile" "test" {
  name           = %q
  mode           = "ah"
  ike_policy     = netbox_ike_policy.test.id
  ipsec_policy   = netbox_ipsec_policy.test.id
  description    = "Test IPsec profile"
}
`, testAccIPSECProfileResourcePrereqs(name), name)

}

func testAccIPSECProfileResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = "1"
  mode    = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}
`, name+"-ike-policy", name+"-ipsec-policy")
}

func TestAccConsistency_IPSECProfile_LiteralNames(t *testing.T) {
	t.Parallel()

	ikePolicyName := testutil.RandomName("tf-test-ike-policy")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy")
	ipsecProfileName := testutil.RandomName("tf-test-ipsec-profile-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(ipsecProfileName)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, ipsecProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", ipsecProfileName),
				),
			},
			{
				Config:   testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, ipsecProfileName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
				),
			},
		},
	})
}

func testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, profileName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = "1"
  mode    = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}

resource "netbox_ipsec_profile" "test" {
  name         = %q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}
`, ikePolicyName, ipsecPolicyName, profileName)
}
