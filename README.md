# terraform-provider-pingdom #

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
      source = "nordcloud/pingdom"
      version = "1.1.4"
    }
  }
}

variable "pingdom_api_token" {}
variable "solarwinds_user" {}
variable "solarwinds_passwd" {}

provider "pingdom" {
    api_token = "${var.pingdom_api_token}"
    solarwinds_user  = "${var.solarwinds_user}"
    solarwinds_password  = "${var.solarwinds_password}"
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
    -var 'solarwinds_user=YOUR_SOLARWINDS_USER'
    -var 'solarwinds_user=YOUR_SOLARWINDS_PASSWD'
```

**Using attributes from other resources**

```hcl
variable "heroku_email" {}
variable "heroku_api_key" {}

variable "pingdom_api_token" {}
variable "solarwinds_user" {}
variable "solarwinds_passwd" {}



provider "heroku" {
    email = var.heroku_email
    api_key = var.heroku_api_key
}

provider "pingdom" {
    api_token = "${var.pingdom_api_token}"
    solarwinds_user  = "${var.solarwinds_user}"
    solarwinds_password  = "${var.solarwinds_password}"
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

**Maintenance**

The maintenance resource is used to define a one time or repetitive maintenance window and bond the maintenance window with one or more uptime and tms checks.

```hcl
resource "pingdom_check" "test" {
  name = "test-check"
  host = "www.example.com"
  type = "http"
}

resource "pingdom_maintenance" "test" {
  description    = "test-maintenance"
  from           = 2717878693
  to             = 2718878693
  effectiveto    = 2718978693
  recurrencetype = "week"
  repeatevery    = 4
  uptimeids      = [pingdom_check.test.id]
}
```

**User**

An user resource is either an user invitation or an active user on the Solarwinds user portal.
These users can be configured to access a range of applications, including Pingdom. An user becomes active once he
manually accepts the invitation sent to his email as specified at the time the invitation is created.

```hcl

resource "pingdom_user" "user" {
  email = "foo@nordcloud.com"
  role = "MEMBER"
  products {
    name = "APPOPTICS"
    role = "MEMBER"
  }
  products {
    name = "PINGDOM"
    role = "MEMBER"
  }
}
```

**Occurrence**

An occurrence resource usually represents a group of maintenance occurrences, as determined by the triple
(maintenanceid, effective_from, effective_to). This triple is effectively a query against all existing
maintenance occurrences. Please note that 'effective_from' and 'effective_to' are different from the attributes
pair 'from' and 'to' of maintenance/occurrence. The latter are used to specify the start/end within a maintenance
cycle, while the former are purely query conditions used to retrieve occurrence objects.

It is not possible to import occurrences as they are queries, which does not exist on Pingdom.

'from' and 'to' attributes can be updated, which results in all occurrences matched by the query being
updated. In other words, all occurrences matched by a single resource will share the same values for `from` and `to`.

```hcl

resource "pingdom_occurrence" "test" {
    maintenance_id = pingdom_maintenance.test.id
    effective_from = pingdom_maintenance.test.from
    effective_to = pingdom_maintenance.test.effectiveto
    from = pingdom_maintenance.test.from
    to = "2021-04-10T22:00:00+08:00"
}
```

## Resources ##

### Pingdom Check ###

#### Common Attributes ####

The following common attributes for all check types can be set:

  * **name** - (Required) The name of the check

  * **host** - (Required) The hostname to check.  Should be in the format `example.com`.

  * **type** - (Required) The check type.  Allowed values: (http, ping, tcp, dns).

  * **resolution** - The time in minutes between each check. Allowed values: (1,5,15,30,60). Default is `5`

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

  * **verify_certificate** - Treat target site as down if an invalid/unverifiable certificate is found. Allowed values (bool): `true`, `false`

  * **ssl_down_days_before** - Treat the target site as down if a certificate expires within the given number of days. This parameter will be ignored if `verify_certificate` is set to `false`. Default value is 0.

#### TCP specific attributes ####

For the TCP checks, you can set these attributes:

  * **port** - Target port for TCP checks.

  * **stringtosend** - (optional) This string will be sent to the port

  * **stringtoexpect** - (optional) This string must be returned by the remote host for the check to pass

#### DNS specific attributes ####

For the DNS checks, you can set these attributes:

  * **expectedip** - The expected IP address of the `host`.

  * **nameserver** - The DNS server used to resolve the host to IP address.

The following attributes are exported:

  * **id** The ID of the Pingdom check

### Pingdom TMS Check ###
 * **name** - (Required) The name of the TMS check
 * **steps** - (Required) At least one steps block is needed to describe the actions of the transaction
    * **args** - The arguments for the function of the step
    * **fn** - The function to perform for the step
 * **active** - Whether the TMS check is enabled
 * **contact_ids** - The id of the contacts to be notified
 * **custom_message** - The message to send in alerts
 * **integration_ids** - The id of integrations
 * **interval** - The interval in which the check is performed, can only be one of [5 10 20 60 720 1440]. Default value is 10
 * **metadata** - The metadata is for recording transactions only
    * **authentication** - Authentication information
    * **disable_websecurity** - 
    * **height** -
    * **width** -
 * **region** - The region within which the check is performed, default is 'us-east'  
 * **send_notification_when_down** - the number of times for the check to fail before the site is considered down, default is 1.
 * **security_level** - how important are the alerts when the check fails. Allowed values: low, high. Default is 'high'
 * **tags** - List of tags for a check. The tag name may contain the characters 'A-Z', 'a-z', '0-9', '_' and '-'. The maximum length of a tag is 64 characters.
 * **team_ids** - Teams to alert

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
    
### Pingdom Maintenance ###

  * **description** - (Required) The name of the team

  * **from** - (Required) Initial maintenance window start. RFC3339 format time like `2066-01-02T22:00:00+08:00`

  * **to** - (Required) Initial maintenance window end. RFC3339 format time like `2066-01-02T22:00:00+08:00`

  * **effectiveto** - Recurrence end. RFC3339 format time like `2066-01-02T22:00:00+08:00` Default: equal to `to`.

  * **recurrencetype** - Type of recurrence. Allowed values: `none` `day` `week` `month`. Default is `none`
  
  * **repeatevery** - Repeat every n-th day/week/month. Default is `0`

  * **tmsids** - Identifiers of transaction checks to assign to the maintenance window - Comma separated Integers

  * **uptimeids** - Identifiers of uptime checks to assign to the maintenance window - Comma separated Integers

### Pingdom Maintenance Occurrence ###

* **maintenance_id** - (Required) The id of the maintenance which the occurrence belongs to, 
  please use references to maintenance resources.

* **effective_from** - (Required) The start time of the occurrence query, RFC3339 format time like `2066-01-02T22:00:00+08:00`. 
  If not specified, the default value is the current time, which means only future occurrences will be returned. This is
  usually desired because it only makes sense to manipulate future occurrences.
  NOTE: this is for query only, not related to the actual start time of the maintenance window.

* **effective_to** - (Required) The end time of the occurrence query, RFC3339 format time like `2066-01-02T22:00:00+08:00`. 
  NOTE: this is for query only, not related to the actual end time of the maintenance window.

* **from** - The start time of an occurrence, RFC3339 format time like `2066-01-02T22:00:00+08:00`.

* **to** - The end time of an occurrence, RFC3339 format time like `2066-01-02T22:00:00+08:00`.

* **size** - The total number of occurrences match the current query. This just serves as read-only information.

### Pingdom User ###

* **email**: (Required) Email of the contact, an invitation will be sent to this address
* **role**: The role in the Solarwinds adminpanel, possible values: "MEMBER", "ADMIN"
* **products**: Permission to each application, the list must be comprehensive. The user will have access to application listed only
    * **name**: The name of the application, possible values: "APPOPTICS", "PINGDOM", "LOGGLY", "PAPERTRAIL"
    * **role**: The permission, allowed values may vary for each service. For Pingdom, possible values: "ADMIN", "OWNER", "VIEWER", "EDITOR", "NO_ACCESS"
    
### Pingdom Integration ###

  * **provider_name** - (Required) The name of the integration provider,One of webhook|librato. 'librato' not presently operational

  * **active** -  (Required) The status of the integration

  * **name**:  (Required) The integration name

  * **url**:  (Optional)  The integration url, only required while provider is webhook

      

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
