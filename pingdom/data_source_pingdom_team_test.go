package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePingdomTeam_basic(t *testing.T) {
	resourceName := "pingdom_team.test"
	datasourceName := "data.pingdom_team.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomTeamConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingdomTeamDataSourceID(datasourceName),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "member_ids.#", resourceName, "member_ids.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "member_ids.0", resourceName, "member_ids.0"),
				),
			},
		},
	})
}

func testAccCheckPingdomTeamDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Team data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Team data source ID not set")
		}
		return nil
	}
}

func testAccDataSourcePingdomTeamConfig(name string) string {
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

data "pingdom_team" "test" {
  name = pingdom_team.test.name
}
`, name, name)
}
