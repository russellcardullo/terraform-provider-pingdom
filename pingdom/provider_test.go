package pingdom

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/russellcardullo/go-pingdom"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"pingdom": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	var expectedUser string
	var expectedPassword string
	var expectedKey string

	if v := os.Getenv("PINGDOM_USER"); v != "" {
		expectedUser = v
	} else {
		expectedUser = "foo"
	}

	if v := os.Getenv("PINGDOM_PASSWORD"); v != "" {
		expectedPassword = v
	} else {
		expectedPassword = "foo"
	}

	if v := os.Getenv("PINGDOM_API_KEY"); v != "" {
		expectedKey = v
	} else {
		expectedKey = "foo"
	}

	raw := map[string]interface{}{
		"user":     expectedUser,
		"password": expectedPassword,
		"api_key":  expectedKey,
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	rp := Provider()
	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config := rp.Meta().(*pingdom.Client)
	if config.User != expectedUser {
		t.Fatalf("bad: %#v", config)
	}

	if config.Password != expectedPassword {
		t.Fatalf("bad: %#v", config)
	}

	if config.APIKey != expectedKey {
		t.Fatalf("bad: %#v", config)
	}
}
