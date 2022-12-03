package pingdom

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/DrFaust92/go-pingdom/pingdom"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePingdomContact_basic(t *testing.T) {
	resourceName := "pingdom_contact.test"
	name := acctest.RandomWithPrefix("tf-acc-test")
	sms := pingdom.SMSNotification{
		CountryCode: "1",
		Number:      "66666666",
		Provider:    "nexmo",
		Severity:    "HIGH",
	}
	email := pingdom.EmailNotification{
		Address:  "test@example.com",
		Severity: "LOW",
	}

	updatedName := acctest.RandomWithPrefix("tf-acc-test")
	updatedSMS := pingdom.SMSNotification{
		CountryCode: "86",
		Number:      "88888888",
		Provider:    "esendex",
		Severity:    "LOW",
	}
	updatedEmail := pingdom.EmailNotification{
		Address:  "testupdate@example.com",
		Severity: "HIGH",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomContactDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomContactConfig(name, false, sms, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.number", "66666666"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.provider", "nexmo"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.country_code", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_notification.0.address", "test@example.com"),
					resource.TestCheckResourceAttr(resourceName, "email_notification.0.severity", "LOW"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomContactConfig(updatedName, true, updatedSMS, updatedEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email_notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.number", updatedSMS.Number),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.severity", updatedSMS.Severity),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.provider", updatedSMS.Provider),
					resource.TestCheckResourceAttr(resourceName, "sms_notification.0.country_code", updatedSMS.CountryCode),
					resource.TestCheckResourceAttr(resourceName, "email_notification.0.address", updatedEmail.Address),
					resource.TestCheckResourceAttr(resourceName, "email_notification.0.severity", updatedEmail.Severity),
				),
			},
		},
	})
}

func testAccCheckPingdomContactDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).Pingdom

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_contact" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Contact ID is not valid: %s", rs.Primary.ID)
		}

		resp, err := client.Contacts.Read(id)
		if err == nil {
			if strconv.Itoa(resp.ID) == rs.Primary.ID {
				return fmt.Errorf("Contact (%s) still exists.", rs.Primary.ID)
			}
		}

		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}

	return nil
}

func testAccResourcePingdomContactConfig(name string, pause bool, sms pingdom.SMSNotification, email pingdom.EmailNotification) string {
	return fmt.Sprintf(`
resource "pingdom_contact" "test" {
	name = "%s"
	paused = %v
	sms_notification {
		number   = "%s"
		severity = "%s"
		provider = "%s"
		country_code = "%s"
	}
	email_notification {
		address   = "%s"
		severity = "%s"
	}
}
`, name, pause, sms.Number, sms.Severity, sms.Provider, sms.CountryCode, email.Address, email.Severity)
}
