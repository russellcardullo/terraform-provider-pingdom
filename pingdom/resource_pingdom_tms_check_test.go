package pingdom

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePingdomTmsCheck(t *testing.T) {
	resourceName := "pingdom_tms_check.test"
	contactResourceName := "pingdom_contact.test"
	teamResourceName := "pingdom_team.test"
	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomTmsCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomTmsCheckConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "steps.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "contact_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_message", ""),
					resource.TestCheckResourceAttr(resourceName, "integrationids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "metadata.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east"),
					resource.TestCheckResourceAttr(resourceName, "tags", ""),
					resource.TestCheckResourceAttr(resourceName, "teamids.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomTmsCheckConfig_update(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "steps.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "contact_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_message", "custom alert message"),
					resource.TestCheckResourceAttr(resourceName, "integrationids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "interval", "20"),
					resource.TestCheckResourceAttr(resourceName, "metadata.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east"),
					resource.TestCheckResourceAttr(resourceName, "tags", "bar,foo"),
					resource.TestCheckResourceAttr(resourceName, "team_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "contact_ids.0", contactResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_ids.0", teamResourceName, "id"),
				),
			},
		},
	})
}

func testAccCheckPingdomTmsCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).Pingdom

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_tms_check" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Check ID is not valid: %s", rs.Primary.ID)
		}

		resp, err := client.TMSCheck.List()
		if err != nil {
			return err
		}
		for _, ck := range resp {
			if ck.ID == id {
				return fmt.Errorf("TMS Check (%s) still exists.", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccResourcePingdomTmsCheckConfig(name string) string {
	return fmt.Sprintf(`
resource "pingdom_tms_check" "test" {
  name = "%s"
  active = true
  steps {
    args = {
      url = "www.ibm.com"
    }
    fn = "go_to"
  }
}
`, name)
}

func testAccResourcePingdomTmsCheckConfig_update(name string) string {
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

resource "pingdom_tms_check" "test" {
  name = "%s"
  active = true
  contact_ids = [
    pingdom_contact.test.id
  ]
  custom_message = "custom alert message"
  interval = 20
  region = "us-east"
  send_notification_when_down = 1
  security_level = "high"
  tags = "foo,bar"
  team_ids = [
	pingdom_team.test.id
  ]

  steps {
    args = {
      url = "www.ibm.com"
    }
    fn = "go_to"
  }

  steps {
    args = {
      element = "Red Hat OpenShift"
    }
    fn = "click"
  }

  steps {
    args = {
      element = "Try it at no charge Arrow Right"
    }
    fn = "click"
  }
}
`, name, name, name)
}
