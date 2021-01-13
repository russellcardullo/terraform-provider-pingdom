package pingdom

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"pingdom": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	var expectedToken string

	if v := os.Getenv("PINGDOM_API_TOKEN"); v != "" {
		expectedToken = v
	} else {
		expectedToken = "foo"
	}

	raw := map[string]interface{}{
		"api_token": expectedToken,
	}

	rp := Provider().(*schema.Provider)
	err := rp.Configure(terraform.NewResourceConfigRaw(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config := rp.Meta().(*pingdom.Client)

	if config.APIToken != expectedToken {
		t.Fatalf("bad: %#v", config)
	}
}
