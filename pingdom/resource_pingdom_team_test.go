package pingdom

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePingdomTeam_basic(t *testing.T) {
	resourceName := "pingdom_team.test"
	contact1ResourceName := "pingdom_contact.test"
	contact2ResourceName := "pingdom_contact.test_update"

	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomTeamConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "member_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "member_ids.0", contact1ResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomTeamConfigUpdate(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "member_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "member_ids.*", contact1ResourceName, "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "member_ids.*", contact2ResourceName, "id"),
				),
			},
		},
	})
}

func testAccCheckPingdomTeamDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).Pingdom

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_team" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Team ID is not valid: %s", rs.Primary.ID)
		}

		resp, err := client.Teams.Read(id)
		if err == nil {
			if strconv.Itoa(resp.ID) == rs.Primary.ID {
				return fmt.Errorf("Team (%s) still exists.", rs.Primary.ID)
			}
		}

		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	return nil
}

func testAccResourcePingdomTeamConfig(name string) string {
	return fmt.Sprintf(`
resource "pingdom_contact" "test" {
	name = "tf-acc-test-pingdom-contact1"
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
`, name)
}

func testAccResourcePingdomTeamConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "pingdom_contact" "test" {
	name = "tf-acc-test-pingdom-contact1"
	sms_notification {
		number   = "66666666"
		severity = "HIGH"
	}
	sms_notification {
		number   = "88888888"
		severity = "LOW"
	}
}

resource "pingdom_contact" "test_update" {
	name = "tf-acc-test-pingdom-contact2"
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
	member_ids = [pingdom_contact.test.id, pingdom_contact.test_update.id]
}
`, name)
}
