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

func TestIPSecPolicyDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecPolicyDataSource()

	if d == nil {

		t.Fatal("Expected non-nil IPSec policy data source")

	}

}

func TestIPSecPolicyDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecPolicyDataSource()

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

	requiredAttrs := []string{"id", "name", "description", "comments", "tags", "proposals", "pfs_group"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestIPSecPolicyDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecPolicyDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_policy"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestIPSecPolicyDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecPolicyDataSource().(*datasources.IPSecPolicyDataSource)

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

func TestAccIPSecPolicyDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-policy-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecPolicyDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_policy.test", "id"),
				),
			},
		},
	})

}

func TestAccIPSecPolicyDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-policy-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecPolicyDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_policy.test", "id"),
				),
			},
		},
	})

}

func testAccIPSecPolicyDataSourceByID(name string) string {

	return fmt.Sprintf(`







resource "netbox_ipsec_policy" "test" {







  name = %[1]q







}















data "netbox_ipsec_policy" "test" {







  id = netbox_ipsec_policy.test.id







}







`, name)

}

func testAccIPSecPolicyDataSourceByName(name string) string {

	return fmt.Sprintf(`







resource "netbox_ipsec_policy" "test" {







  name = %[1]q







}















data "netbox_ipsec_policy" "test" {







  name = netbox_ipsec_policy.test.name







}







`, name)

}
