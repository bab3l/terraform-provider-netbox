package provider

import (
	"context"
	"os"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure NetboxProvider satisfies various provider interfaces.
var _ provider.Provider = &NetboxProvider{}

// NetboxProvider defines the provider implementation.
type NetboxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NetboxProviderModel describes the provider data model.
type NetboxProviderModel struct {
	ServerURL types.String `tfsdk:"server_url"`
	APIToken  types.String `tfsdk:"api_token"`
	Insecure  types.Bool   `tfsdk:"insecure"`
}

func (p *NetboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netbox"
	resp.Version = p.version
}

func (p *NetboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Netbox provider is used to interact with Netbox, an open-source web application designed to help manage and document computer networks. This provider allows you to manage Netbox resources such as sites, devices, racks, IP addresses, and more using Terraform.",

		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "The base URL of your Netbox instance (e.g., `https://netbox.example.com`). Can also be set via the `NETBOX_SERVER_URL` environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token for authenticating with Netbox. Generate this token in your Netbox user profile. Can also be set via the `NETBOX_API_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip TLS certificate verification. Defaults to false. Can also be set via the `NETBOX_INSECURE` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *NetboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NetboxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	serverURL := os.Getenv("NETBOX_SERVER_URL")
	apiToken := os.Getenv("NETBOX_API_TOKEN")
	insecure := os.Getenv("NETBOX_INSECURE") == "true"

	if !data.ServerURL.IsNull() {
		serverURL = data.ServerURL.ValueString()
	}

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if !data.Insecure.IsNull() {
		insecure = data.Insecure.ValueBool()
	}

	// Validation
	if serverURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing Netbox Server URL",
			"The provider cannot create the Netbox API client as there is a missing or empty value for the Netbox server URL. "+
				"Set the server_url value in the configuration or use the NETBOX_SERVER_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Netbox API Token",
			"The provider cannot create the Netbox API client as there is a missing or empty value for the Netbox API token. "+
				"Set the api_token value in the configuration or use the NETBOX_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "netbox_server_url", serverURL)
	ctx = tflog.SetField(ctx, "netbox_api_token", apiToken)
	ctx = tflog.SetField(ctx, "netbox_insecure", insecure)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "netbox_api_token")

	tflog.Debug(ctx, "Creating Netbox client")

	// Create a new Netbox client using go-netbox
	cfg := netbox.NewConfiguration()
	cfg.Servers = netbox.ServerConfigurations{
		{
			URL: serverURL,
		},
	}

	// Set up authentication
	cfg.DefaultHeader = map[string]string{
		"Authorization": "Token " + apiToken,
	}

	// Handle insecure connections
	if insecure {
		// Note: go-netbox uses a standard HTTP client, so TLS verification
		// would need to be configured on the HTTP client if needed
		tflog.Debug(ctx, "Insecure mode enabled - TLS verification disabled")
	}

	client := netbox.NewAPIClient(cfg)

	// Make the Netbox client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Netbox client", map[string]any{"success": true})
}

func (p *NetboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSiteResource,
		resources.NewSiteGroupResource,
		resources.NewTenantGroupResource,
		resources.NewTenantResource,
		resources.NewPlatformResource,
		resources.NewManufacturerResource,
		resources.NewRegionResource,
		resources.NewLocationResource,
		resources.NewRackResource,
		resources.NewRackRoleResource,
		resources.NewDeviceRoleResource,
		resources.NewDeviceTypeResource,
		resources.NewDeviceResource,
		resources.NewInterfaceResource,
		resources.NewVRFResource,
		resources.NewVLANGroupResource,
		resources.NewVLANResource,
		resources.NewPrefixResource,
		resources.NewIPAddressResource,
		// Phase 3: Virtualization resources
		resources.NewClusterTypeResource,
		resources.NewClusterResource,
		resources.NewVirtualMachineResource,
		resources.NewVMInterfaceResource,
		// Phase 4: Circuits & Connectivity resources
		resources.NewProviderResource,
		resources.NewCircuitTypeResource,
		resources.NewCircuitResource,
		resources.NewCableResource,
		// Phase 5: Extras & Customization resources
		resources.NewTagResource,
		resources.NewContactResource,
		resources.NewWebhookResource,
		resources.NewConfigContextResource,
		// Phase 6: Additional resources
		resources.NewContactGroupResource,
		resources.NewContactRoleResource,
		resources.NewClusterGroupResource,
		// Phase 7: IPAM & Advanced resources
		resources.NewIPRangeResource,
		resources.NewRIRResource,
		resources.NewAggregateResource,
		resources.NewProviderAccountResource,
		resources.NewCircuitTerminationResource,
		resources.NewCustomFieldResource,
		// Phase 8: Additional IPAM & Infrastructure resources
		resources.NewRoleResource,
		resources.NewASNResource,
		resources.NewProviderNetworkResource,
		resources.NewRackTypeResource,
		resources.NewVirtualChassisResource,
		// Phase 9: DCIM Device Components
		resources.NewDeviceBayResource,
		resources.NewPowerPanelResource,
		resources.NewPowerFeedResource,
		resources.NewServiceResource,
		resources.NewWirelessLANResource,
		resources.NewWirelessLANGroupResource,
		resources.NewWirelessLinkResource,
		resources.NewInventoryItemResource,
		resources.NewInventoryItemRoleResource,
		resources.NewModuleTypeResource,
		resources.NewModuleResource,
		resources.NewModuleBayResource,
		resources.NewConsolePortResource,
		resources.NewConsoleServerPortResource,
		resources.NewPowerPortResource,
		resources.NewPowerOutletResource,
		resources.NewConsolePortTemplateResource,
		resources.NewConsoleServerPortTemplateResource,
		resources.NewPowerPortTemplateResource,
		resources.NewPowerOutletTemplateResource,
		resources.NewInterfaceTemplateResource,
		resources.NewConfigTemplateResource,
		// Phase 10: Additional IPAM resources
		resources.NewRouteTargetResource,
		// Phase 11: Virtualization additions
		resources.NewVirtualDiskResource,
		// Phase 12: Additional IPAM
		resources.NewASNRangeResource,
		// Phase 13: DCIM Templates
		resources.NewDeviceBayTemplateResource,
		// Phase 14: VPN resources
		resources.NewIKEProposalResource,
		resources.NewIKEPolicyResource,
		resources.NewIPSecProposalResource,
		resources.NewIPSecPolicyResource,
		resources.NewIPSecProfileResource,
		resources.NewTunnelGroupResource,
		resources.NewTunnelResource,
		resources.NewTunnelTerminationResource,
		// Phase 15: L2VPN resources
		resources.NewL2VPNResource,
		resources.NewL2VPNTerminationResource,
		// Phase 16: Circuit Groups
		resources.NewCircuitGroupResource,
		resources.NewCircuitGroupAssignmentResource,
		// Phase 17: Front/Rear Ports
		resources.NewRearPortTemplateResource,
		resources.NewFrontPortTemplateResource,
		resources.NewRearPortResource,
		resources.NewFrontPortResource,
		// Phase 18: FHRP Groups
		resources.NewFHRPGroupResource,
		// Phase 19: Documentation (Extras)
		resources.NewJournalEntryResource,
		// Phase 20: Custom Field Extensions
		resources.NewCustomFieldChoiceSetResource,
		resources.NewCustomLinkResource,
		// Phase 21: Event Rules and Automation
		resources.NewEventRuleResource,
		resources.NewNotificationGroupResource,
		// Phase 22: DCIM Templates & Infrastructure
		resources.NewRackReservationResource,
		resources.NewVirtualDeviceContextResource,
		resources.NewModuleBayTemplateResource,
		// Note: CableTerminationResource removed - use netbox_cable with embedded terminations instead
		resources.NewInventoryItemTemplateResource,
		// Phase 23: Contact Assignments
		resources.NewContactAssignmentResource,
		// Phase 24: Service Templates
		resources.NewServiceTemplateResource,
		// Phase 25: FHRP Group Assignments
		resources.NewFHRPGroupAssignmentResource,
		// Phase 26: Export Templates
		resources.NewExportTemplateResource,
	}
}

func (p *NetboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSiteDataSource,
		datasources.NewSiteGroupDataSource,
		datasources.NewTenantGroupDataSource,
		datasources.NewTenantDataSource,
		datasources.NewManufacturerDataSource,
		datasources.NewPlatformDataSource,
		datasources.NewRegionDataSource,
		datasources.NewLocationDataSource,
		datasources.NewRackDataSource,
		datasources.NewRackRoleDataSource,
		datasources.NewDeviceRoleDataSource,
		datasources.NewDeviceTypeDataSource,
		datasources.NewDeviceDataSource,
		datasources.NewInterfaceDataSource,
		datasources.NewVRFDataSource,
		datasources.NewVLANGroupDataSource,
		datasources.NewVLANDataSource,
		datasources.NewPrefixDataSource,
		datasources.NewIPAddressDataSource,
		// Phase 3: Virtualization data sources
		datasources.NewClusterTypeDataSource,
		datasources.NewClusterDataSource,
		datasources.NewVirtualMachineDataSource,
		datasources.NewVMInterfaceDataSource,
		// Phase 4: Circuits & Connectivity data sources
		datasources.NewProviderDataSource,
		datasources.NewCircuitTypeDataSource,
		datasources.NewCircuitDataSource,
		datasources.NewCableDataSource,
		// Phase 5: Extras & Customization data sources
		datasources.NewTagDataSource,
		datasources.NewContactDataSource,
		datasources.NewWebhookDataSource,
		datasources.NewConfigContextDataSource,
		// Phase 6: Additional data sources
		datasources.NewContactGroupDataSource,
		datasources.NewContactRoleDataSource,
		datasources.NewClusterGroupDataSource,
		// Phase 7: IPAM & Advanced data sources
		datasources.NewIPRangeDataSource,
		datasources.NewRIRDataSource,
		datasources.NewAggregateDataSource,
		datasources.NewProviderAccountDataSource,
		datasources.NewCircuitTerminationDataSource,
		datasources.NewCustomFieldDataSource,
		// Phase 8: Additional IPAM & Infrastructure data sources
		datasources.NewRoleDataSource,
		datasources.NewASNDataSource,
		datasources.NewProviderNetworkDataSource,
		datasources.NewRackTypeDataSource,
		datasources.NewVirtualChassisDataSource,
		// Phase 9: DCIM Device Components
		datasources.NewDeviceBayDataSource,
		datasources.NewPowerPanelDataSource,
		datasources.NewPowerFeedDataSource,
		datasources.NewServiceDataSource,
		datasources.NewWirelessLANDataSource,
		datasources.NewWirelessLANGroupDataSource,
		datasources.NewWirelessLinkDataSource,
		datasources.NewInventoryItemDataSource,
		datasources.NewInventoryItemRoleDataSource,
		datasources.NewModuleTypeDataSource,
		datasources.NewModuleDataSource,
		datasources.NewModuleBayDataSource,
		datasources.NewConsolePortDataSource,
		datasources.NewConsoleServerPortDataSource,
		datasources.NewPowerPortDataSource,
		datasources.NewPowerOutletDataSource,
		datasources.NewConsolePortTemplateDataSource,
		datasources.NewConsoleServerPortTemplateDataSource,
		datasources.NewPowerPortTemplateDataSource,
		datasources.NewPowerOutletTemplateDataSource,
		datasources.NewInterfaceTemplateDataSource,
		datasources.NewConfigTemplateDataSource,
		// Phase 10: Additional IPAM data sources
		datasources.NewRouteTargetDataSource,
		// Phase 11: Virtualization additions
		datasources.NewVirtualDiskDataSource,
		// Phase 12: Additional IPAM
		datasources.NewASNRangeDataSource,
		// Phase 13: DCIM Templates
		datasources.NewDeviceBayTemplateDataSource,
		// Phase 14: VPN data sources
		datasources.NewIKEProposalDataSource,
		datasources.NewIKEPolicyDataSource,
		datasources.NewIPSecProposalDataSource,
		datasources.NewIPSecPolicyDataSource,
		datasources.NewIPSecProfileDataSource,
		datasources.NewTunnelGroupDataSource,
		datasources.NewTunnelDataSource,
		datasources.NewTunnelTerminationDataSource,
		// Phase 15: L2VPN data sources
		datasources.NewL2VPNDataSource,
		datasources.NewL2VPNTerminationDataSource,
		// Phase 16: Circuit Groups
		datasources.NewCircuitGroupDataSource,
		datasources.NewCircuitGroupAssignmentDataSource,
		// Phase 17: Front/Rear Ports
		datasources.NewRearPortTemplateDataSource,
		datasources.NewFrontPortTemplateDataSource,
		datasources.NewRearPortDataSource,
		datasources.NewFrontPortDataSource,
		// Phase 18: FHRP Groups
		datasources.NewFHRPGroupDataSource,
		// Phase 19: Documentation (Extras)
		datasources.NewJournalEntryDataSource,
		// Phase 20: Custom Field Extensions
		datasources.NewCustomFieldChoiceSetDataSource,
		datasources.NewCustomLinkDataSource,
		// Phase 21: Event Rules and Automation
		datasources.NewEventRuleDataSource,
		datasources.NewNotificationGroupDataSource,
		// Phase 22: DCIM Templates & Infrastructure
		datasources.NewRackReservationDataSource,
		datasources.NewVirtualDeviceContextDataSource,
		datasources.NewModuleBayTemplateDataSource,
		datasources.NewCableTerminationDataSource,
		datasources.NewInventoryItemTemplateDataSource,
		// Phase 23: Users
		datasources.NewUserDataSource,
		// Phase 24: Contact Assignments
		datasources.NewContactAssignmentDataSource,
		// Phase 25: Service Templates
		datasources.NewServiceTemplateDataSource,
		// Phase 26: FHRP Group Assignments
		datasources.NewFHRPGroupAssignmentDataSource,
		// Phase 27: Export Templates
		datasources.NewExportTemplateDataSource,
		// Phase 28: Scripts (read-only)
		datasources.NewScriptDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NetboxProvider{
			version: version,
		}
	}
}
