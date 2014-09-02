package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/russellcardullo/terraform-provider-pingdom/pingdom"
)

func main() {
	plugin.Serve(pingdom.Provider())
}
