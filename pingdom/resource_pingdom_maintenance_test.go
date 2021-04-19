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

func TestAccResourcePingdomMaintenance_basic(t *testing.T) {
	resourceName := "pingdom_maintenance.test"
	checkResourceName := "pingdom_check.test"

	description := acctest.RandomWithPrefix("tf-acc-test")
	updatedDescription := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomMaintenanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomMaintenanceConfig(description),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "from", "2066-01-02T22:00:00+08:00"),
					resource.TestCheckResourceAttr(resourceName, "to", "2066-01-02T23:00:00+08:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomMaintenanceConfigUpdate(updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "from", "2066-01-05T22:10:00+08:00"),
					resource.TestCheckResourceAttr(resourceName, "to", "2066-01-05T23:59:00+08:00"),
					resource.TestCheckResourceAttr(resourceName, "effectiveto", "2088-10-22T06:07:08+08:00"),
					resource.TestCheckResourceAttr(resourceName, "recurrencetype", "week"),
					resource.TestCheckResourceAttr(resourceName, "repeatevery", "4"),
					resource.TestCheckResourceAttr(resourceName, "uptimeids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "uptimeids.*", checkResourceName, "id"),
				),
			},
		},
	})
}

func testAccCheckPingdomMaintenanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).Pingdom

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_maintenance" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Maintenance ID is not valid: %s", rs.Primary.ID)
		}

		resp, err := client.Maintenances.Read(id)
		if err == nil {
			if strconv.Itoa(resp.ID) == rs.Primary.ID {
				return fmt.Errorf("Maintenance (%s) still exists.", rs.Primary.ID)
			}
		}

		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	return nil
}

func testAccResourcePingdomMaintenanceConfig(name string) string {
	return fmt.Sprintf(`
resource "pingdom_maintenance" "test" {
	description = "%s"
	from        = "2066-01-02T22:00:00+08:00"
	to          = "2066-01-02T23:00:00+08:00"
}
`, name)
}

func testAccResourcePingdomMaintenanceConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "test" {
	name = "%s"
	host = "www.example.com"
	type = "http"
}

resource "pingdom_maintenance" "test" {
	description    = "%s"
	from           = "2066-01-05T22:10:00+08:00"
	to             = "2066-01-05T23:59:00+08:00"
	effectiveto    = "2088-10-22T06:07:08+08:00"
	recurrencetype = "week"
	repeatevery    = 4
	uptimeids      = [pingdom_check.test.id]
}
`, name, name)
}
