package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/jayhding/terraform-provider-pingdom/pingdom"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pingdom.Provider,
	})
}
