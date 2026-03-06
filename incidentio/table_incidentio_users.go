package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// UserRole represents a base or custom role assigned to a user.
type UserRole struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}

// FullUser represents a user record from the /v2/users endpoint.
type FullUser struct {
	BaseRole    *UserRole  `json:"base_role"`
	CustomRoles []UserRole `json:"custom_roles"`
	Email       string     `json:"email"`
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	SlackUserID string     `json:"slack_user_id"`
}

// usersResponse is the envelope returned by GET /v2/users.
type usersResponse struct {
	Users          []FullUser     `json:"users"`
	PaginationMeta PaginationMeta `json:"pagination_meta"`
}

//// TABLE DEFINITION

func tableIncidentioUsers() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_users",
		Description: "List users in your incident.io account.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "email", Require: plugin.Optional},
				{Name: "slack_user_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the user."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Full name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "System role of the user (viewer, responder, administrator)."},
			{Name: "slack_user_id", Type: proto.ColumnType_STRING, Description: "Slack user ID of the user."},
			{Name: "base_role_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("BaseRole.ID"), Description: "ID of the user's base role."},
			{Name: "base_role_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("BaseRole.Name"), Description: "Name of the user's base role."},
			{Name: "base_role_slug", Type: proto.ColumnType_STRING, Transform: transform.FromField("BaseRole.Slug"), Description: "Slug of the user's base role."},
			{Name: "custom_roles", Type: proto.ColumnType_JSON, Description: "Custom roles assigned to the user."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	after := ""
	for {
		params := map[string]string{"page_size": "250"}
		if after != "" {
			params["after"] = after
		}
		if v := d.EqualsQualString("email"); v != "" {
			params["email"] = v
		}
		if v := d.EqualsQualString("slack_user_id"); v != "" {
			params["slack_user_id"] = v
		}

		var result usersResponse
		if err := client.get(ctx, "/v2/users", params, &result); err != nil {
			return nil, fmt.Errorf("listing users: %w", err)
		}

		for _, u := range result.Users {
			d.StreamListItem(ctx, u)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.PaginationMeta.After == "" || len(result.Users) == 0 {
			break
		}
		after = result.PaginationMeta.After
	}

	return nil, nil
}

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, nil
	}

	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result struct {
		User FullUser `json:"user"`
	}
	if err := client.get(ctx, fmt.Sprintf("/v2/users/%s", id), nil, &result); err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}

	return result.User, nil
}
