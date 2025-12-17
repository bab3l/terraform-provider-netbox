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
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFHRPGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil FHRP Group resource")

	}

}

func TestFHRPGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required attributes

	requiredAttrs := []string{"protocol", "group_id"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Optional attributes

	optionalAttrs := []string{"name", "auth_type", "auth_key", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestFHRPGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_fhrp_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestFHRPGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupResource().(*resources.FHRPGroupResource)

	// Test with nil provider data

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct client type

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with incorrect provider data type

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccFHRPGroupResource_basic(t *testing.T) {

	protocol := "vrrp2"

	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

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

				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
		},
	})

}

func TestAccFHRPGroupResource_full(t *testing.T) {

	protocol := "hsrp"

	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

	name := testutil.RandomName("tf-test-fhrp")

	description := "Test FHRP Group with all fields"

	authType := "plaintext"

	authKey := "secretkey123"

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

				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, name, description, authType, authKey),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", description),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", authType),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", authKey),
				),
			},
		},
	})

}

func TestAccFHRPGroupResource_update(t *testing.T) {

	protocol := "vrrp3"

	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

	updatedName := testutil.RandomName("tf-test-fhrp-updated")

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

				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
				),
			},

			{

				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, updatedName, "Updated description", "md5", "newsecret456"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", "md5"),

					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", "newsecret456"),
				),
			},
		},
	})

}

func TestAccFHRPGroupResource_protocols(t *testing.T) {

	// Test different protocol values

	protocols := []string{"vrrp2", "vrrp3", "carp", "hsrp", "glbp", "clusterxl", "other"}

	for _, protocol := range protocols {

		t.Run(protocol, func(t *testing.T) {

			groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32

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

						Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),

						Check: resource.ComposeTestCheckFunc(

							resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),

							resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
						),
					},
				},
			})

		})

	}

}

func testAccFHRPGroupResourceConfig_basic(protocol string, groupID int32) string {

	return fmt.Sprintf(`































resource "netbox_fhrp_group" "test" {































  protocol = %q































  group_id = %d































}































`, protocol, groupID)

}

func testAccFHRPGroupResourceConfig_full(protocol string, groupID int32, name, description, authType, authKey string) string {

	return fmt.Sprintf(`































resource "netbox_fhrp_group" "test" {































  protocol    = %q































  group_id    = %d































  name        = %q































  description = %q































  auth_type   = %q































  auth_key    = %q































}































`, protocol, groupID, name, description, authType, authKey)

}

func TestAccFHRPGroupResource_import(t *testing.T) {
	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) //nolint:gosec // G115

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
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				ResourceName:      "netbox_fhrp_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
