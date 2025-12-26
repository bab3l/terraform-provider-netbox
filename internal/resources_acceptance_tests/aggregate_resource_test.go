package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateResource_basic(t *testing.T) {

	t.Parallel()
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	// Use a random third octet to ensure uniqueness across test runs
	prefix := fmt.Sprintf("192.0.%d.0/24", acctest.RandIntRange(0, 255))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
				),
			},

			{

				ResourceName: "netbox_aggregate.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})

}

func TestAccAggregateResource_full(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-full")

	rirSlug := testutil.RandomSlug("tf-test-rir-full")

	// Use a random third octet to ensure uniqueness across test runs
	prefix := fmt.Sprintf("198.51.%d.0/24", acctest.RandIntRange(0, 255))

	description := testutil.RandomName("description")

	updatedDescription := "Updated aggregate description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, description, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", testutil.Comments),
				),
			},

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, updatedDescription, testutil.Comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccAggregateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-id")
	rirSlug := testutil.RandomSlug("tf-test-rir-id")
	prefix := fmt.Sprintf("203.0.%d.0/24", acctest.RandIntRange(0, 255))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
				),
			},
		},
	})

}

func testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

resource "netbox_aggregate" "test" {

  prefix = %q

  rir    = netbox_rir.test.id

}

`, rirName, rirSlug, prefix)

}

func testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

resource "netbox_aggregate" "test" {

  prefix      = %q

  rir         = netbox_rir.test.id

  description = %q

  comments    = %q

}

`, rirName, rirSlug, prefix, description, comments)

}

func TestAccConsistency_Aggregate(t *testing.T) {

	t.Parallel()
	// Use random second octet to ensure uniqueness across test runs
	// /16 requires the last two octets to be 0
	prefix := fmt.Sprintf("10.%d.0.0/16", acctest.RandIntRange(0, 255))

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},

			{

				PlanOnly: true,

				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}

resource "netbox_aggregate" "test" {

  prefix = "%[1]s"

  rir = netbox_rir.test.id

  tenant = netbox_tenant.test.name

}

`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}

// TestAccConsistency_Aggregate_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_Aggregate_LiteralNames(t *testing.T) {

	t.Parallel()
	// Use random second octet to ensure uniqueness across test runs
	// /16 requires the last two octets to be 0
	prefix := fmt.Sprintf("10.%d.0.0/16", acctest.RandIntRange(0, 255))

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}

resource "netbox_aggregate" "test" {

  prefix = "%[1]s"

  # Use literal string names to mimic existing user state

  rir = "%[3]s"

  tenant = "%[4]s"

  depends_on = [netbox_rir.test, netbox_tenant.test]

}

`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}
