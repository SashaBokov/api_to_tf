// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
	"regexp"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RuntimeGroup{}
var _ resource.ResourceWithImportState = &RuntimeGroup{}

func NewRuntimeGroup() resource.Resource {
	return &RuntimeGroup{}
}

// RuntimeGroup defines the resource implementation.
type RuntimeGroup struct {
	client *client.Client
}

// RuntimeGroupModel describes the resource data model.
type RuntimeGroupModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ClusterType          types.String `tfsdk:"cluster_type"`
	Labels               types.Map    `tfsdk:"labels"`
	ControlPlaneEndpoint types.String `tfsdk:"control_plane_endpoint"`
	TelemetryEndpoint    types.String `tfsdk:"telemetry_endpoint"`
}

func (r *RuntimeGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *RuntimeGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runtime group resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the runtime group.",
				Required:            true,
			},
			"Description": schema.StringAttribute{
				MarkdownDescription: "The description of the runtime group in Konnect.",
				Optional:            true,
			},
			"cluster_type": schema.StringAttribute{
				MarkdownDescription: "The ClusterType value of the cluster associated with the Runtime Group.",
				Optional:            true,
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "Labels to facilitate tagged search on runtime groups. Keys must be of length 1-63 characters, and cannot start with 'kong', 'konnect', 'mesh', 'kic'.",
				Optional:            true,
				Validators: []validator.Map{mapvalidator.KeysAre(stringvalidator.RegexMatches(regexp.MustCompile("^(?!kong|konnect|mesh|kic).{1,63}$"),
					"Keys must be of length 1-63 characters, and cannot start with 'kong', 'konnect', 'mesh', 'kic'."))},
				ElementType: types.StringType,
			},
			"control_plane_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"telemetry_endpoint": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service generated identifier for the Runtime Group.",
			},
		},
	}
}

func (r *RuntimeGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RuntimeGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RuntimeGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	labels := make(map[string]string)
	for k, v := range data.Labels.Elements() {
		labels[k] = v.String()
	}

	createReq := client.CreateRuntimeGroupRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		ClusterType: data.ClusterType.ValueString(),
		Labels:      labels,
	}

	createResp, err := r.client.CreateRuntimeGroup(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create, got error: %s", err))
		return
	}

	data.ControlPlaneEndpoint = types.StringValue(createResp.Config.ControlPlaneEndpoint)
	data.TelemetryEndpoint = types.StringValue(createResp.Config.TelemetryEndpoint)
	data.Id = data.Name

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuntimeGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RuntimeGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuntimeGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RuntimeGroupModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RuntimeGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RuntimeGroupModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *RuntimeGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
