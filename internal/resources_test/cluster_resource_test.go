package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClusterResource(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterResource()

	if r == nil {

		t.Fatal("Expected non-nil Cluster resource")
	}
}

func TestClusterResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "type"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"status", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestClusterResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_cluster"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestClusterResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewClusterResource().(*resources.ClusterResource)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccClusterResource_basic(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")

	clusterName := testutil.RandomName("tf-test-cluster")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),
				),
			},
		},
	})
}

func TestAccClusterResource_full(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-full")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-full")

	clusterName := testutil.RandomName("tf-test-cluster-full")

	description := "Test cluster with all fields"

	comments := "Test comments"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),

					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "description", description),

					resource.TestCheckResourceAttr("netbox_cluster.test", "comments", comments),
				),
			},
		},
	})
}

func TestAccClusterResource_update(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-update")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-update")

	clusterName := testutil.RandomName("tf-test-cluster-update")

	updatedName := testutil.RandomName("tf-test-cluster-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterCleanup(updatedName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
				),
			},

			{

				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, updatedName, "Updated description", "Updated comments"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.slug
}

`, clusterTypeName, clusterTypeSlug, clusterName)
}

func testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name        = %q
  type        = netbox_cluster_type.test.slug
  status      = "active"
  description = %q
  comments    = %q
}

`, clusterTypeName, clusterTypeSlug, clusterName, description, comments)
}

func TestAccClusterResource_import(t *testing.T) {

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-import")

	clusterTypeSlug := clusterTypeName

	clusterName := testutil.RandomName("tf-test-cluster-import")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_import(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),
				),
			},

			{

				ResourceName: "netbox_cluster.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func testAccClusterResourceConfig_import(clusterTypeName, clusterTypeSlug, clusterName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.slug
}

`, clusterTypeName, clusterTypeSlug, clusterName)
}

func TestAccConsistency_Cluster(t *testing.T) {

	t.Parallel()

	clusterName := testutil.RandomName("cluster")

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	groupName := testutil.RandomName("group")

	groupSlug := testutil.RandomSlug("group")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterConsistencyConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),

					resource.TestCheckResourceAttr("netbox_cluster.test", "group", groupSlug),

					resource.TestCheckResourceAttr("netbox_cluster.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccClusterConsistencyConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccClusterConsistencyConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_cluster_group" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_site" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_tenant" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_cluster" "test" {
  name = "%[1]s"
  type = netbox_cluster_type.test.slug
  group = netbox_cluster_group.test.slug
  site = netbox_site.test.name
  tenant = netbox_tenant.test.name
}

`, clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug)
}
