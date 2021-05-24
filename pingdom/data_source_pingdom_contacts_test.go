package pingdom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePingdomContacts_basic(t *testing.T) {
	datasourceName := "data.pingdom_contacts.test"
	resourceName := "pingdom_contact.test"
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePingdomContactsConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(datasourceName),
					resource.TestCheckTypeSetElemAttrPair(datasourceName, "ids.*", resourceName, "id"),
					resource.TestCheckTypeSetElemAttrPair(datasourceName, "names.*", resourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourcePingdomContactsConfig(name string) string {
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

data "pingdom_contacts" "test" {

}
`, name)
}
