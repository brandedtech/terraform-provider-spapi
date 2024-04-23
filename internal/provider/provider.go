package provider

import (
	"context"
	"os"

	sp "github.com/brandedtech/sp-api-sdk/pkg/selling-partner"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &SPAPIProvider{}
)

type SPAPIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type SPAPIProviderModel struct {
	LWAClientID     types.String `tfsdk:"lwa_client_id"`
	LWAClientSecret types.String `tfsdk:"lwa_client_secret"`
	RefreshToken    types.String `tfsdk:"refresh_token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SPAPIProvider{
			version: version,
		}
	}
}

func (p *SPAPIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "spapi"
	resp.Version = p.version
}

func (p *SPAPIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"lwa_client_id": schema.StringAttribute{
				Optional: true,
			},
			"lwa_client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"refresh_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *SPAPIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Selling Partner client")

	var config SPAPIProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.LWAClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("lwa_client_id"),
			"Unknown LWA client ID",
			"The provider cannot create the SP-API client as there is an unknown configuration value for the SP-API LWA client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_LWA_CLIENT_ID environment variable.",
		)
	}

	if config.LWAClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("lwa_client_secret"),
			"Unknown LWA client secret",
			"The provider cannot create the SP-API client as there is an unknown configuration value for the SP-API LWA client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_LWA_CLIENT_SECRET environment variable.",
		)
	}

	if config.RefreshToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Unknown refresh token",
			"The provider cannot create the SP-API client as there is an unknown configuration value for the SP-API refresh token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_REFRESH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	lwaClientID := os.Getenv("SP_API_LWA_CLIENT_ID")
	lwaClientSecret := os.Getenv("SP_API_LWA_CLIENT_SECRET")
	refreshToken := os.Getenv("SP_API_REFRESH_TOKEN")

	if !config.LWAClientID.IsNull() {
		lwaClientID = config.LWAClientID.ValueString()
	}

	if !config.LWAClientSecret.IsNull() {
		lwaClientSecret = config.LWAClientSecret.ValueString()
	}

	if !config.RefreshToken.IsNull() {
		refreshToken = config.RefreshToken.ValueString()
	}

	if lwaClientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("lwa_client_id"),
			"SP-API LWA client ID is not set",
			"The provider cannot create the SP-API client as there is no SP-API LWA client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_LWA_CLIENT_ID environment variable.",
		)
	}

	if lwaClientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("lwa_client_secret"),
			"SP-API LWA client secret is not set",
			"The provider cannot create the SP-API client as there is no SP-API LWA client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_LWA_CLIENT_SECRET environment variable.",
		)
	}

	if refreshToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"SP-API refresh token is not set",
			"The provider cannot create the SP-API client as there is no SP-API refresh token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SP_API_REFRESH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	sellingPartner, err := sp.NewSellingPartner(&sp.Config{
		ClientID:     lwaClientID,
		ClientSecret: lwaClientSecret,
		RefreshToken: refreshToken,
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create SP-API client",
			"The provider cannot create the SP-API client as there was an error creating the SP-API client. \n\n"+
				"HashiCups Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = sellingPartner
	resp.ResourceData = sellingPartner
}

func (p *SPAPIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNotificationDestinationResource,
		NewNotificationSubscriptionResource,
	}
}

func (p *SPAPIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNotificationDestinationsDatasource,
	}
}

func (p *SPAPIProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}
