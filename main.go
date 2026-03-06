package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/vthiery/steampipe-plugin-incidentio/incidentio"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: incidentio.Plugin})
}
