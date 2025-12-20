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
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFHRPGroupDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewFHRPGroupDataSource()

	if d == nil {

		t.Fatal("Expected non-nil FHRP Group data source")

	}

}

func TestFHRPGroupDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewFHRPGroupDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required/Optional for lookup

	lookupAttrs := []string{"id", "protocol", "group_id"}

	for _, attr := range lookupAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected lookup attribute %s to exist in schema", attr)

		}

	}

	// Computed attributes

	computedAttrs := []string{"name", "auth_type", "description", "comments"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

}

func TestFHRPGroupDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewFHRPGroupDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_fhrp_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestFHRPGroupDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewFHRPGroupDataSource().(*datasources.FHRPGroupDataSource)

	// Test with nil provider data

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct client type

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with incorrect provider data type

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccFHRPGroupDataSource_byID(t *testing.T) {

	protocol := "vrrp2"

	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

	name := testutil.RandomName("tf-test-fhrp-ds-id")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckFHRPGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccFHRPGroupDataSourceConfig_byID(protocol, groupID, name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "protocol", protocol),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "name", name),
				),
			},
		},
	})

}

func TestAccFHRPGroupDataSource_byProtocolAndGroupID(t *testing.T) {

	protocol := "hsrp"

	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

	name := testutil.RandomName("tf-test-fhrp-ds-lookup")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckFHRPGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccFHRPGroupDataSourceConfig_byProtocolAndGroupID(protocol, groupID, name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "protocol", protocol),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),

					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "name", name),
				),
			},
		},
	})

}

func testAccFHRPGroupDataSourceConfig_byID(protocol string, groupID int32, name string) string {

	return fmt.Sprintf(`







resource "netbox_fhrp_group" "test" {







  protocol = %q







  group_id = %d



  name     = %q



}







data "netbox_fhrp_group" "test" {







  id = netbox_fhrp_group.test.id



}







`, protocol, groupID, name)

}

func testAccFHRPGroupDataSourceConfig_byProtocolAndGroupID(protocol string, groupID int32, name string) string {

	return fmt.Sprintf(`







resource "netbox_fhrp_group" "test" {







  protocol = %q







  group_id = %d



  name     = %q



}







data "netbox_fhrp_group" "test" {







  protocol = netbox_fhrp_group.test.protocol







  group_id = netbox_fhrp_group.test.group_id



}







`, protocol, groupID, name)

}
