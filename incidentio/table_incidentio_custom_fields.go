package incidentio

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// CustomFieldOption represents one selectable option for a custom field.
type CustomFieldOption struct {
	CustomFieldID string `json:"custom_field_id"`
	ID            string `json:"id"`
	SortKey       int    `json:"sort_key"`
	Value         string `json:"value"`
}

// CustomField represents a custom field definition.
type CustomField struct {
	CatalogTypeID          string              `json:"catalog_type_id"`
	CreatedAt              string              `json:"created_at"`
	Description            string              `json:"description"`
	FieldType              string              `json:"field_type"`
	ID                     string              `json:"id"`
	Name                   string              `json:"name"`
	Options                []CustomFieldOption `json:"options"`
	Required               string              `json:"required"`
	RequiredV2             string              `json:"required_v2"`
	ShowBeforeClosure      bool                `json:"show_before_closure"`
	ShowBeforeCreation     bool                `json:"show_before_creation"`
	ShowBeforeUpdate       bool                `json:"show_before_update"`
	ShowInAnnouncementPost bool                `json:"show_in_announcement_post"`
	UpdatedAt              string              `json:"updated_at"`
}

// customFieldsResponse is the envelope returned by GET /v1/custom_fields.
type customFieldsResponse struct {
	CustomFields []CustomField `json:"custom_fields"`
}

//// TABLE DEFINITION

func tableIncidentioCustomFields() *plugin.Table {
	return &plugin.Table{
		Name:        "incidentio_custom_fields",
		Description: "List custom field definitions configured in incident.io.",
		List: &plugin.ListConfig{
			Hydrate: listCustomFields,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the custom field."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the custom field."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the custom field."},
			{Name: "field_type", Type: proto.ColumnType_STRING, Description: "Type of the field (single_select, multi_select, text, link, numeric, catalog_entry)."},
			{Name: "required", Type: proto.ColumnType_STRING, Description: "Whether the field is required (never, before_closure, always)."},
			{Name: "required_v2", Type: proto.ColumnType_STRING, Description: "Updated required setting."},
			{Name: "catalog_type_id", Type: proto.ColumnType_STRING, Description: "If the field is a catalog entry type, the catalog type ID."},
			{Name: "show_before_creation", Type: proto.ColumnType_BOOL, Description: "Whether to show the field before incident creation."},
			{Name: "show_before_closure", Type: proto.ColumnType_BOOL, Description: "Whether to show the field before incident closure."},
			{Name: "show_before_update", Type: proto.ColumnType_BOOL, Description: "Whether to show the field before incident updates."},
			{Name: "show_in_announcement_post", Type: proto.ColumnType_BOOL, Description: "Whether to show the field in announcement posts."},
			{Name: "options", Type: proto.ColumnType_JSON, Description: "Available options for select-type fields."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the custom field was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time the custom field was last updated."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listCustomFields(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result customFieldsResponse
	if err := client.get(ctx, "/v1/custom_fields", nil, &result); err != nil {
		return nil, fmt.Errorf("listing custom fields: %w", err)
	}

	for _, f := range result.CustomFields {
		d.StreamListItem(ctx, f)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
