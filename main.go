package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/russellcardullo/terraform-provider-pingdom/pingdom"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pingdom.Provider,
	})
}
