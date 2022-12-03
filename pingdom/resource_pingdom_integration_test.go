package pingdom

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/DrFaust92/go-pingdom/pingdomext"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPingdomIntegration_basic(t *testing.T) {
	var integration pingdomext.IntegrationGetResponse
	resourceName := "pingdom_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPingdomIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookIntegration_basicConfig("webhook", false, "test-3", "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingdomIntegrationExists(resourceName, &integration),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-3"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://www.example.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPingdomIntegration_update(t *testing.T) {
	var integration pingdomext.IntegrationGetResponse
	resourceName := "pingdom_integration.test"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckPingdomIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWebhookIntegration_basicConfig("webhook", false, "test-3", "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingdomIntegrationExists(resourceName, &integration),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-3"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://www.example.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWebhookIntegration_basicConfig("webhook", true, "test-4", "https://www.example2.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPingdomIntegrationExists(resourceName, &integration),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "webhook"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", "test-4"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://www.example2.com"),
				),
			},
		},
	})
}

func testAccCheckPingdomIntegrationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Clients).PingdomExt

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pingdom_integration" {
			continue
		}

		resID, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return err
		}

		integration, err := client.Integrations.Read(resID)
		if err != nil {
			return err
		}
		if integration != nil {
			return fmt.Errorf("Error deleting integration %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPingdomIntegrationExists(n string, obj *pingdomext.IntegrationGetResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testAccProvider.Meta().(*Clients).PingdomExt

		resID, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return err
		}
		integration, err := client.Integrations.Read(resID)
		if err != nil {
			return err
		}
		if integration == nil {
			return fmt.Errorf("Error finding integration %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccWebhookIntegration_basicConfig(providerName string, active bool, name string, url string) string {
	return fmt.Sprintf(`resource "pingdom_integration" "test" {
		provider_name = "%s"
		active = %t
		name="%s"
		url="%s"
	  }
	  `, providerName, active, name, url)
}
