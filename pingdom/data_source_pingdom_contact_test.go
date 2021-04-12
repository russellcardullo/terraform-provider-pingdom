package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePingdomContact_basic(t *testing.T) {
	contactResourceName := "pingdom_contact.test"
	teamResourceName := "pingdom_team.test"
	datasourceName := "data.pingdom_contact.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomContactConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "name", name),
					resource.TestCheckResourceAttrPair(datasourceName, "paused", contactResourceName, "paused"),
					resource.TestCheckResourceAttrPair(datasourceName, "sms_notification.#", contactResourceName, "sms_notification.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "email_notification.#", contactResourceName, "email_notification.#"),
					resource.TestCheckResourceAttr(datasourceName, "teams.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceName, "teams.0.id", teamResourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "teams.0.name", teamResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourcePingdomContactConfig(name string) string {
	return fmt.Sprintf(`
resource "pingdom_contact" "test" {
	name = "%s"
	sms_notification {
		number   = "66666666"
		severity = "HIGH"
	}
	email_notification {
		address  = "test@test.com"
		severity = "LOW"
	}
}

resource "pingdom_team" "test" {
	name = "%s"
	member_ids = [pingdom_contact.test.id]
}

data "pingdom_contact" "test" {
  name = pingdom_contact.test.name
}
`, name, name)
}
