# terraform-provider-pingdom #

This project is a [terraform](http://www.terraform.io/) provider for [pingdom](https://www.pingdom.com/).

This currently only supports working with basic HTTP and ping checks.

## Build and install ##

### Dependencies ###

You should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This use [dep](https://github.com/golang/dep) for dependency management.  To fetch all dependencies run `dep ensure`.

### Build ###

```
go get github.com/aylien/terraform-provider-pingdom
```

The binary will then be available at `~/go/bin/terraform-provider-pingdom`, unless you've set `$GOPATH`, in which case you'll find it at `$GOPATH/bin`

### Install ###

You will need to install the binary as a [terraform third party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).  Terraform will then pick up the binary from the local filesystem when you run `terraform init`.

```
ln -s ~/go/bin/terraform-provider-pingdom ~/.terraform.d/plugins/$(uname | tr '[:upper:]' '[:lower:]')_amd64/terraform-provider-pingdom_v$(date +%Y.%m.%d)
```

Updates can now be applied by re-running `go get github.com/aylien/terraform-provider-pingdom`.

## Usage ##

**Basic Check**

```
variable "pingdom_user" {}
variable "pingdom_password" {}
variable "pingdom_api_key" {}
variable "pingdom_account_email" {} # Optional: only required for multi-user accounts

provider "pingdom" {
    user = "${var.pingdom_user}"
    password = "${var.pingdom_password}"
    api_key = "${var.pingdom_api_key}"
    account_email = "${var.pingdom_account_email}" # Optional: only required for multi-user accounts
}

resource "pingdom_check" "example" {
    type = "http"
    name = "my http check"
    host = "example.com"
    resolution = 5
}

resource "pingdom_check" "example_with_alert" {
    type = "http"
    name = "my http check"
    host = "example.com"
    resolution = 5
    sendnotificationwhendown = 2
    integrationids = [
      12345678,
      23456789
    ]
}

resource "pingdom_check" "ping_example" {
    type = "ping"
    name = "my ping check"
    host = "example.com"
    resolution = 1
}
```

Apply with:
```
 terraform apply \
    -var 'pingdom_user=YOUR_USERNAME' \
    -var 'pingdom_password=YOUR_PASSWORD' \
    -var 'pingdom_api_key=YOUR_API_KEY'
```

**Using attributes from other resources**

```
variable "heroku_email" {}
variable "heroku_api_key" {}

variable "pingdom_user" {}
variable "pingdom_password" {}
variable "pingdom_api_key" {}

provider "heroku" {
    email = "${var.heroku_email}"
    api_key = "${var.heroku_api_key}"
}

provider "pingdom" {
    user = "${var.pingdom_user}"
    password = "${var.pingdom_password}"
    api_key = "${var.pingdom_api_key}"
}

resource "heroku_app" "example" {
    name = "my-app"
    region = "us"
}

resource "pingdom_check" "example" {
    name = "my check"
    host = "${heroku_app.example.heroku_hostname}"
    resolution = 5
}
```
## Resources ##

### Pingdom Check ###

#### Common Attibutes ####

The following common attributes for all check types can be set:

**name** - (Required) The name of the check

**host** - (Required) The hostname to check.  Should be in the format `example.com`.

**resolution** - (Required) The time in minutes between each check.  Allowed values: (1,5,15,30,60).

**type** - (Required) The check type.  Allowed values: (http, ping).

**sendnotificationwhendown** - The number of consecutive failed checks required to trigger an alert. Values of 0 are ignored.

**notifyagainevery** - Notify again after n results.  A value of 0 means no additional notifications will be sent.

**notifywhenbackup** - Notify when backup.

**integrationids** - List of integer integration IDs (defined by webhook URL) that will be triggered by the alerts. The ID can be extracted from the integrations page URL on the pingdom website.

**userids** - List of integer user IDs that will be notified when the check is down.

**teamids** - List of integer team IDs that will be notified when the check is down.

#### HTTP specific attibutes ####

For the HTTP checks, you can set these attributes:

**url** - Target path on server.

**encryption** - Enable encryption in the HTTP check (aka HTTPS).

**port** - Target port for HTTP checks.

**username** - Username for target HTTP authentication.

**password** - Password for target HTTP authentication.

**shouldcontain** - Target site should contain this string.

**shouldnotcontain** - Target site should NOT contain this string. Not allowed defined together with `shouldcontain`.

**postdata** - Data that should be posted to the web page, for example submission data for a sign-up or login form. The data needs to be formatted in the same way as a web browser would send it to the web server.

**requestheaders** - Custom HTTP headers. It should be a hash with pairs, like `{ "header_name" = "header_content" }`

**tags** - List of tags the check should contain. Should be in the format "tagA,tagB"

**probefilters** - Region from which the check should originate. One of NA, EU, APAC, or LATAM. Should be in the format "region:NA"

The following attributes are exported:

**id** The ID of the Pingdom check
