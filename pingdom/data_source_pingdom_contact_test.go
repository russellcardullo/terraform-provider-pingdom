package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourcePingdomContact_basic(t *testing.T) {
	resourceName := "pingdom_contact.test"
	datasourceName := "data.pingdom_contact.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomContactConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingdomContactDataSourceID(datasourceName),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "paused", resourceName, "paused"),
					resource.TestCheckResourceAttrPair(datasourceName, "sms_notification.#", resourceName, "sms_notification.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "email_notification.#", resourceName, "email_notification.#"),
					resource.TestCheckResourceAttr(datasourceName, "teams.#", "0"),
				),
			},
		},
	})
}

func testAccCheckPingdomContactDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Contact data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Contact data source ID not set")
		}
		return nil
	}
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

data "pingdom_contact" "test" {
  name = pingdom_contact.test.name
}
`, name)
}
