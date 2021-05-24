package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePingdomTeams_basic(t *testing.T) {
	datasourceName := "data.pingdom_teams.test"
	resourceName := "pingdom_team.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomTeamsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(datasourceName),
					resource.TestCheckTypeSetElemAttrPair(datasourceName, "ids.*", resourceName, "id"),
					resource.TestCheckTypeSetElemAttrPair(datasourceName, "names.*", resourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourcePingdomTeamsConfig(name string) string {
	return fmt.Sprintf(`
resource "pingdom_contact" "test" {
	name = "%s"
	sms_notification {
		number   = "66666666"
		severity = "HIGH"
	}
	sms_notification {
		number   = "88888888"
		severity = "LOW"
	}
}

resource "pingdom_team" "test" {
	name = "%s"
	member_ids = [pingdom_contact.test.id]
}

data "pingdom_teams" "test" {
  depends_on = [
    pingdom_team.test,
  ]
}
`, name, name)
}
