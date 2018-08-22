package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/aylien/terraform-provider-pingdom/pingdom"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pingdom.Provider,
	})
}
