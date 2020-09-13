# terraform-provider-pingdom #

This project is a [terraform](http://www.terraform.io/) provider for [pingdom](https://www.pingdom.com/).

This currently only supports working with basic HTTP and ping checks.

## Build and install ##

### Using released versions ###

Prebuild releases for most platforms are available [here](https://github.com/russellcardullo/terraform-provider-pingdom/releases).
Download the release corresponding to your particular platform and place in `$HOME/.terraform.d/plugins/[os]_[arch]`.  For instance
on Linux AMD64 the path would be `$HOME/.terraform.d/plugins/linux_amd64`.

After copying the plugin run `terraform init` in your projects that use this provider.

### Dependencies for building from source ###

If you need to build from source, you should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) for dependency management.  To fetch all dependencies run `go get` inside this repository.

### Build ###

```
make build
```

The binary will then be available at `_build/terraform-provider-pingdom_VERSION`.

### Install ###

```
make install
```

This will place the binary under `$HOME/.terraform.d/plugins/OS_ARCH/terraform-provider-pingdom_VERSION`.  After installing you will need to run `terraform init` in any project using the plugin.

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
    email = var.heroku_email
    api_key = var.heroku_api_key
}

provider "pingdom" {
    user = var.pingdom_user
    password = var.pingdom_password
    api_key = var.pingdom_api_key
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

```
resource "pingdom_team" "test" {
  name = "The Test team"
  userids = [
    pingdom_user.first_user.id,
  ]
}
```

**Users**

```
resource "pingdom_user" "first_user" {
  username = "johndoe"
}

resource "pingdom_user" "second_user" {
  username = "janedoe"
}
```

**Contacts**

```

resource "pingdom_contact" "first_user_contact_email_2" {
  user_id        = pingdom_user.first_user.id
  email          = "john.doe@doe.com"
  severity_level = "LOW"
}

resource "pingdom_contact" "first_user_contact_sms_1" {
  user_id        = pingdom_user.first_user.id
  number         = "700000000"
  country_code   = "33"
  phone_provider = "nexmo"
  severity_level = "HIGH"
}

resource "pingdom_user" "second_user" {
  username = "janedoe"
}

resource "pingdom_contact" "second_user_contact_email_1" {
  user_id        = pingdom_user.second_user.id
  email          = "jane@doe.com"
  severity_level = "high"
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

  * **sendnotificationwhendown** - The number of consecutive failed checks required to trigger an alert. Values of 0 are ignored. See note about interaction with `integrationids` below.

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

  * **publicreport** - If `true`, this check will be included in the public report (default: `false`)

#### TCP specific attributes ####

For the TCP checks, you can set these attributes:

  * **port** - Target port for TCP checks.

  * **stringtosend** - (optional) This string will be sent to the port

  * **stringtoexpect** - (optional) This string must be returned by the remote host for the check to pass

The following attributes are exported:

  * **id** The ID of the Pingdom check


### Pingdom Team ###

  * **name** - (Required) The name of the team

  * **userids** - List of integer user IDs that will be notified when the check is down.


### Pingdom User ###

  * **username** - (Required) The name of the user


### Pingdom Contact ###

  * **user_id**: (Required) ID of the user linked to this contact

  * **severity_level**: (Required) Severity level for target

  * **email**: Email

  * **number**: Cellphone number, without the country code part. (Requires countrycode)

  * **country_code**: Cellphone country code (Requires number)

  * **phone_provider**: SMS provider (Requires number and countrycode)
