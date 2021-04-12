package pingdom

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"pingdom": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	var expectedToken string
	var expectedUser string
	var expectedPassword string

	if v := os.Getenv("PINGDOM_API_TOKEN"); v != "" {
		expectedToken = v
	} else {
		expectedToken = "foo"
	}

	if v := os.Getenv("SOLARWINDS_USER"); v != "" {
		expectedUser = v
	} else {
		expectedUser = "foo"
	}

	if v := os.Getenv("SOLARWINDS_PASSWD"); v != "" {
		expectedPassword = v
	} else {
		expectedPassword = "foo"
	}

	raw := map[string]interface{}{
		"api_token":         expectedToken,
		"solarwinds_user":   expectedUser,
		"solarwinds_passwd": expectedPassword,
	}
	var isAccTestEnabled bool
	if v := os.Getenv("TF_ACC"); v != "" {
		isAccTestEnabled = true
	}

	if isAccTestEnabled {
		rp := Provider()
		err := rp.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
		if err != nil {
			t.Fatal(err)
		}

		config := rp.Meta().(*Clients).Pingdom

		if config.APIToken != expectedToken {
			t.Fatalf("bad: %#v", config)
		}
	}

}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PINGDOM_API_TOKEN"); v == "" {
		t.Fatal("PINGDOM_API_TOKEN environment variable must be set for acceptance tests")
	}
	if v := os.Getenv("SOLARWINDS_USER"); v == "" {
		t.Fatal("SOLARWINDS_USER environment variable must be set for acceptance tests")
	}
	if v := os.Getenv("SOLARWINDS_PASSWD"); v == "" {
		t.Fatal("SOLARWINDS_PASSWD environment variable must be set for acceptance tests")
	}
}

func testAccCheckPingdomResourceID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID not set: %s", name)
		}
		return nil
	}
}
