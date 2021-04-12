package pingdom

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nordcloud/go-pingdom/pingdom"
)

func TestAccResourcePingdomCheck_http(t *testing.T) {
	resourceName := "pingdom_check.http"
	contactResourceName := "pingdom_contact.test"
	teamResourceName := "pingdom_team.test"
	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomCheckConfig_http(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					resource.TestCheckResourceAttr(resourceName, "resolution", "5"),
					resource.TestCheckResourceAttr(resourceName, "sendnotificationwhendown", "2"),
					resource.TestCheckResourceAttr(resourceName, "notifyagainevery", "0"),
					resource.TestCheckResourceAttr(resourceName, "notifywhenbackup", "false"),
					resource.TestCheckResourceAttr(resourceName, "verify_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssl_down_days_before", "10"),
					resource.TestCheckResourceAttr(resourceName, "url", "/"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "responsetime_threshold", "30000"),
					resource.TestCheckResourceAttr(resourceName, "postdata", ""),
					resource.TestCheckResourceAttr(resourceName, "integrationids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags", ""),
					resource.TestCheckResourceAttr(resourceName, "probefilters.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "userids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "teamids.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomCheckConfig_http_update(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.org"),
					resource.TestCheckResourceAttr(resourceName, "type", "http"),
					resource.TestCheckResourceAttr(resourceName, "resolution", "15"),
					resource.TestCheckResourceAttr(resourceName, "sendnotificationwhendown", "5"),
					resource.TestCheckResourceAttr(resourceName, "notifyagainevery", "3"),
					resource.TestCheckResourceAttr(resourceName, "notifywhenbackup", "true"),
					resource.TestCheckResourceAttr(resourceName, "verify_certificate", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_down_days_before", "0"),
					resource.TestCheckResourceAttr(resourceName, "url", "/test"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttr(resourceName, "requestheaders.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "requestheaders.X-Test-Data", "test"),
					resource.TestCheckResourceAttr(resourceName, "postdata", "test message"),
					resource.TestCheckResourceAttr(resourceName, "integrationids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags", "a,b"),
					resource.TestCheckResourceAttr(resourceName, "probefilters", "region:APAC"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "userids.0", contactResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "teamids.0", teamResourceName, "id"),
				),
			},
		},
	})
}

func TestAccResourcePingdomCheck_tcp(t *testing.T) {
	resourceName := "pingdom_check.tcp"
	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomCheckConfig_tcp(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomCheckConfig_tcp_update(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.org"),
					resource.TestCheckResourceAttr(resourceName, "type", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
				),
			},
		},
	})
}

func TestAccResourcePingdomCheck_ping(t *testing.T) {
	resourceName := "pingdom_check.ping"
	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomCheckConfig_ping(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "type", "ping"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomCheckConfig_ping_update(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "host", "www.example.org"),
					resource.TestCheckResourceAttr(resourceName, "type", "ping"),
				),
			},
		},
	})
}

func TestAccResourcePingdomCheck_dns(t *testing.T) {
	resourceName := "pingdom_check.dns"
	name := acctest.RandomWithPrefix("tf-acc-test")
	updatedName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPingdomCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePingdomCheckConfig_dns(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "host", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "nameserver", "a.iana-servers.net"),
					resource.TestCheckResourceAttr(resourceName, "expectedip", "2606:2800:220:1:248:1893:25c8:1946"),
					resource.TestCheckResourceAttr(resourceName, "type", "dns"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourcePingdomCheckConfig_dns_update(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPingdomResourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "host", "example.org"),
					resource.TestCheckResourceAttr(resourceName, "nameserver", "b.iana-servers.net"),
					resource.TestCheckResourceAttr(resourceName, "expectedip", "93.184.216.34"),
					resource.TestCheckResourceAttr(resourceName, "type", "dns"),
				),
			},
		},
	})
}

func testAccCheckPingdomCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pingdom.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_check" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Check ID is not valid: %s", rs.Primary.ID)
		}

		resp, err := client.Checks.Read(id)
		if err == nil {
			if strconv.Itoa(resp.ID) == rs.Primary.ID {
				return fmt.Errorf("Check (%s) still exists.", rs.Primary.ID)
			}
		}

		if !strings.Contains(err.Error(), "403") {
			return err
		}
	}

	return nil
}

func testAccResourcePingdomCheckConfig_dns(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "dns" {
	name       = "%s"
	host       = "example.com"
	nameserver = "a.iana-servers.net"
	expectedip = "2606:2800:220:1:248:1893:25c8:1946"
	type       = "dns"
}
`, name)
}

func testAccResourcePingdomCheckConfig_dns_update(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "dns" {
	name       = "%s"
	host       = "example.org"
	nameserver = "b.iana-servers.net"
	expectedip = "93.184.216.34"
	type       = "dns"
}
`, name)
}

func testAccResourcePingdomCheckConfig_ping(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "ping" {
	name = "%s"
	host = "www.example.com"
	type = "ping"
}
`, name)
}

func testAccResourcePingdomCheckConfig_ping_update(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "ping" {
	name = "%s"
	host = "www.example.org"
	type = "ping"
}
`, name)
}

func testAccResourcePingdomCheckConfig_tcp(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "tcp" {
	name = "%s"
	host = "www.example.com"
	port = 80
	type = "tcp"
}
`, name)
}

func testAccResourcePingdomCheckConfig_tcp_update(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "tcp" {
	name = "%s"
	host = "www.example.org"
	port = 443
	type = "tcp"
}
`, name)
}

func testAccResourcePingdomCheckConfig_http(name string) string {
	return fmt.Sprintf(`
resource "pingdom_check" "http" {
	name                 = "%s"
	host                 = "www.example.com"
	type                 = "http"
	ssl_down_days_before = 10
}
`, name)
}

func testAccResourcePingdomCheckConfig_http_update(name string) string {
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

resource "pingdom_check" "http" {
	name                     = "%s"
	host                     = "www.example.org"
	type                     = "http"
	resolution               = 15
	sendnotificationwhendown = 5
	notifyagainevery         = 3
	notifywhenbackup         = true
	verify_certificate       = false
	ssl_down_days_before     = 0
	url                      = "/test"
	port                     = 443
	postdata                 = "test message"
	tags                     = "a,b"
	probefilters             = "region:APAC"
	shouldnotcontain         = "shouldnotcontain"
	username                 = "user"
	password                 = "password"
	encryption               = true
	paused                   = true
	userids                  = [pingdom_contact.test.id]
	teamids                  = [pingdom_team.test.id]
	requestheaders = {
		X-Test-Data = "test"
	}
}
`, name, name, name)
}
