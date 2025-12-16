package datasources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIPSecProfileDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProfileDataSource()

	if d == nil {

		t.Fatal("Expected non-nil IPSec profile data source")

	}

}

func TestIPSecProfileDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProfileDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check that key attributes exist

	requiredAttrs := []string{"id", "name", "mode", "description", "comments", "tags", "ike_policy", "ipsec_policy"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestIPSecProfileDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProfileDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_profile"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestIPSecProfileDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProfileDataSource().(*datasources.IPSecProfileDataSource)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with invalid provider data")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIPSecProfileDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-profile-ds")

	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-ds")

	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProfileDataSourceByID(randomName, ikePolicyName, ipsecPolicyName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})

}

func TestAccIPSecProfileDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-profile-ds")

	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-ds")

	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProfileDataSourceByName(randomName, ikePolicyName, ipsecPolicyName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_profile.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})

}

func testAccIPSecProfileDataSourceByID(name, ikePolicyName, ipsecPolicyName string) string {

	return fmt.Sprintf(`















resource "netbox_ike_policy" "test" {















  name    = %[2]q















  version = 2















}































resource "netbox_ipsec_policy" "test" {















  name = %[3]q















}































resource "netbox_ipsec_profile" "test" {















  name         = %[1]q















  mode         = "esp"















  ike_policy   = netbox_ike_policy.test.id















  ipsec_policy = netbox_ipsec_policy.test.id















}































data "netbox_ipsec_profile" "test" {















  id = netbox_ipsec_profile.test.id















}















`, name, ikePolicyName, ipsecPolicyName)

}

func testAccIPSecProfileDataSourceByName(name, ikePolicyName, ipsecPolicyName string) string {

	return fmt.Sprintf(`















resource "netbox_ike_policy" "test" {















  name    = %[2]q















  version = 2















}































resource "netbox_ipsec_policy" "test" {















  name = %[3]q















}































resource "netbox_ipsec_profile" "test" {















  name         = %[1]q















  mode         = "esp"















  ike_policy   = netbox_ike_policy.test.id















  ipsec_policy = netbox_ipsec_policy.test.id















}































data "netbox_ipsec_profile" "test" {















  name = netbox_ipsec_profile.test.name















}















`, name, ikePolicyName, ipsecPolicyName)

}
