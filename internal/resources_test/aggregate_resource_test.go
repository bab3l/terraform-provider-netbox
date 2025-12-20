package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAggregateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	if r == nil {

		t.Fatal("Expected non-nil Aggregate resource")

	}

}

func TestAggregateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	testutil.ValidateResourceSchema(t, schemaResponse.Schema.Attributes, testutil.SchemaValidation{

		Required: []string{"prefix", "rir"},

		Optional: []string{"tenant", "date_added", "description", "comments", "tags", "custom_fields"},

		Computed: []string{"id"},
	})

}

func TestAggregateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	testutil.ValidateResourceMetadata(t, r, "netbox", "netbox_aggregate")

}

func TestAggregateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	testutil.ValidateResourceConfigure(t, r)

}

func TestAccAggregateResource_basic(t *testing.T) {

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	prefix := "192.0.2.0/24"

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

	rirName := testutil.RandomName("tf-test-rir-full")

	rirSlug := testutil.RandomSlug("tf-test-rir-full")

	prefix := "198.51.100.0/24"

	description := "Test aggregate with all fields"

	updatedDescription := "Updated aggregate description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", comments),
				),
			},

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, updatedDescription, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", updatedDescription),
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

	prefix := "192.168.0.0/16"

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "tenant", tenantName),
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

  rir = netbox_rir.test.slug

  tenant = netbox_tenant.test.name

}

`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}

// TestAccConsistency_Aggregate_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_Aggregate_LiteralNames(t *testing.T) {

	t.Parallel()

	prefix := "10.50.0.0/16"

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "tenant", tenantName),
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
