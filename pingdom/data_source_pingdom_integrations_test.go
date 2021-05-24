package pingdom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePingdomIntegrations_basic(t *testing.T) {
	datasourceName := "data.pingdom_integrations.all"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomIntegrationsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "ids.#", "7"),
					resource.TestCheckResourceAttr(datasourceName, "names.#", "7"),
				),
			},
		},
	})
}

func testAccDataSourcePingdomIntegrationsConfig() string {
	return `data "pingdom_integrations" "all" {}`
}
