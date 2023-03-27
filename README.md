# terraform-provider-pingdom #

**I no longer use Pingdom and no longer maintain this project.**

This project is a [terraform](http://www.terraform.io/) provider for [pingdom](https://www.pingdom.com/).

This currently only supports working with basic HTTP and ping checks.

This supports Pingdom API v3.1: [API reference docs](https://docs.pingdom.com/api/)

## Requirements ##
* Terraform 0.12.x
* Go 1.14 (to build the provider plugin)

## Usage ##

**Use provider**
```hcl
terraform {
  required_providers {
    pingdom = {
      source = "russellcardullo/pingdom"
      version = "1.1.3"
    }
  }
}

variable "pingdom_api_token" {}

provider "pingdom" {
    api_token = "${var.pingdom_api_token}"
}
```

**Basic Check**
```hcl
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
    sendnotificationwhendown = 2 # alert after 5 mins, with resolution 5*(2-1)
    integrationids = [
      12345678,
      23456789
    ]
    userids = [
      24680,
      13579
    ]
}

resource "pingdom_check" "ping_example" {
    type = "ping"
    name = "my ping check"
    host = "example.com"
    resolution = 1
    userids = [
      24680
    ]
}
```

Apply with:
```sh
 terraform apply \
    -var 'pingdom_api_token=YOUR_API_TOKEN'
```

**Using attributes from other resources**

```hcl
variable "heroku_email" {}
variable "heroku_api_key" {}

variable "pingdom_api_token" {}

provider "heroku" {
    email = var.heroku_email
    api_key = var.heroku_api_key
}

provider "pingdom" {
    api_token = var.pingdom_api_token
}

resource "heroku_app" "example" {
    name = "my-app"
    region = "us"
}

resource "pingdom_check" "example" {
    name = "my check"
    host = heroku_app.example.heroku_hostname
    resolution = 5
}
```

**Teams**

```hcl
resource "pingdom_team" "test" {
  name = "The Test team"
  member_ids = [
    pingdom_contact.first_contact.id,
  ]
}
```

**Contacts**

Note that all contacts _must_ have both a high and low severity notification

```hcl

resource "pingdom_contact" "first_contact" {
  name = "johndoe"

  sms_notification {
    number   = "5555555555"
    severity = "HIGH"
  }

  sms_notification {
    number       = "3333333333"
    country_code = "91"
    severity     = "LOW"
    provider     = "esendex"
  }

  email_notification {
    address  = "test@test.com"
    severity = "LOW"
  }
}

resource "pingdom_contact" "second_contact" {
  name   = "janedoe"
  paused = true

  email_notification {
    address  = "test@test.com"
    severity = "LOW"
  }

  email_notification {
    address  = "test@test.com"
    severity = "HIGH"
  }
}
```

## Resources ##

### Pingdom Check ###

#### Common Attributes ####

The following common attributes for all check types can be set:

  * **name** - (Required) The name of the check

  * **host** - (Required) The hostname to check.  Should be in the format `example.com`.

  * **resolution** - (Required) The time in minutes between each check.  Allowed values: (1,5,15,30,60).

  * **type** - (Required) The check type.  Allowed values: (http, ping).

  * **paused** - Whether the check is active or not (defaults to `false`, if not provided). Allowed values (bool): `true`, `false`

  * **responsetime_threshold** = How long (int: milliseconds) pingdom should wait before marking a probe as failed (defaults to 30000 ms)

  * **sendnotificationwhendown** - The consecutive failed checks required to trigger an alert. Values of 1 imply notification instantly. Values of 2 mean pingdom will wait for a second check to fail, i.e. `resolution` minutes, to trigger an alert. For example `sendnotificationwhendown: 2` and `resolution: 1`, will trigger an alert after 1 minute. Further, values of N will trigger an alert after `(N - 1) * resolution` minutes, e.g. `sendnotificationwhendown: 6` and `resolution: 1` will trigger an alert after 5 minutes. Values of 0 are ignored. See note about interaction with `integrationids` below.

  * **notifyagainevery** - Notify again after n results.  A value of 0 means no additional notifications will be sent.

  * **notifywhenbackup** - Notify when back up.

  * **integrationids** - List of integer integration IDs (defined by webhook URL) that will be triggered by the alerts. The ID can be extracted from the integrations page URL on the pingdom website. See note about interaction with `sendnotificationwhendown` below.

  * **userids** - List of integer user IDs that will be notified when the check is down.

  * **teamids** - List of integer team IDs that will be notified when the check is down.

Note that when using `integrationids`, the `sendnotificationwhendown` value will be ignored when sending webhook notifications.  You may need to contact Pingdom support for more details.  See #52.

#### HTTP specific attributes ####

For the HTTP checks, you can set these attributes:

  * **url** - Target path on server.

  * **encryption** - Enable encryption in the HTTP check (aka HTTPS).

  * **port** - Target port for HTTP checks.

  * **username** - Username for target HTTP authentication.

  * **password** - Password for target HTTP authentication.

  * **shouldcontain** - Target site should contain this string.

  * **shouldnotcontain** - Target site should NOT contain this string. Not allowed defined together with `shouldcontain`.

  * **postdata** - Data that should be posted to the web page, for example submission data for a sign-up or login form. The data needs to be formatted in the same way as a web browser would send it to the web server.

  * **requestheaders** - Custom HTTP headers. It should be a hash with pairs, like `{ "header_name" = "header_content" }`

  * **tags** - List of tags the check should contain. Should be in the format "tagA,tagB"

  * **probefilters** - Region from which the check should originate. One of NA, EU, APAC, or LATAM. Should be in the format "region:NA"

  * **verifycertificate** - Enable monitoring of SSL/TLS certificate.

  * **ssldowndaysbefore** - Days prior to certificate expiring to consider down.

#### TCP specific attributes ####

For the TCP checks, you can set these attributes:

  * **port** - Target port for TCP checks.

  * **stringtosend** - (optional) This string will be sent to the port

  * **stringtoexpect** - (optional) This string must be returned by the remote host for the check to pass

The following attributes are exported:

  * **id** The ID of the Pingdom check


### Pingdom Team ###

  * **name** - (Required) The name of the team

  * **member_ids** - List of integer contact IDs that will be notified when the check is down.


### Pingdom Contact ###

  * **name**: (Required) Name of the contact

  * **paused**: Whether alerts for this contact should be disabled

  * **sms_notification**: Block resource describing an SMS notification

      * **country_code**: The country code, defaults to "1"

      * **number**: The phone number

      * **provider**: Provider for SMS messaging. One of nexmo|bulksms|esendex|cellsynt. 'bulksms' not presently operational

      * **severity**: Severity of this notification. One of HIGH|LOW

  * **email_notification**: Block resource describing an Email notification

      * **address**: Email address to notify

      * **severity**: Severity of this notification. One of HIGH|LOW

## Develop The Provider ##

### Dependencies for building from source ###

If you need to build from source, you should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.  To fetch all dependencies run `go get` inside this repository.

### Build ###

```sh
make build
```

The binary will then be available at `_build/terraform-provider-pingdom_VERSION`.

### Install ###

```sh
make install
```

This will place the binary under `$HOME/.terraform.d/plugins/OS_ARCH/terraform-provider-pingdom_VERSION`.  After installing you will need to run `terraform init` in any project using the plugin.
