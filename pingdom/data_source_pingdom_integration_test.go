package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePingdomIntegration_basic(t *testing.T) {
	resourceName := "pingdom_integration.test"
	datasourceName := "data.pingdom_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomIntegrationConfig("webhook", false, "test-3", "https://www.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(datasourceName),
					resource.TestCheckResourceAttrPair(datasourceName, "provider_name", resourceName, "provider_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "active", resourceName, "active"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "url", resourceName, "url"),
				),
			},
		},
	})
}

func testAccDataSourcePingdomIntegrationConfig(providerName string, active bool, name string, url string) string {
	return fmt.Sprintf(`resource "pingdom_integration" "test" {
		provider_name = "%s"
		active = %t
		name="%s"
		url="%s"
	  }

	  data "pingdom_integration" "test" {
		name = pingdom_integration.test.name
	  }
	  `, providerName, active, name, url)
}
