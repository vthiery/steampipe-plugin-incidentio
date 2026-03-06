package incidentio

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// incidentioConfig stores the connection configuration for the plugin.
type incidentioConfig struct {
	APIKey *string `hcl:"api_key"`
}

// ConfigInstance returns a new instance of incidentioConfig (used by the SDK).
func ConfigInstance() interface{} {
	return &incidentioConfig{}
}

// GetConfig retrieves and casts the connection config from the plugin query data.
func GetConfig(connection *plugin.Connection) incidentioConfig {
	if connection == nil || connection.Config == nil {
		return incidentioConfig{}
	}
	config, _ := connection.Config.(incidentioConfig)
	return config
}
