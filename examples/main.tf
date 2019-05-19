# Copyright (C) 2018-2019 Nicolas Lamirault <nicolas.lamirault@gmail.com>

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

terraform {
  required_version = ">= 0.11.0"
}

provider "pingdom" {
  user     = "${var.pingdom_user}"
  password = "${var.pingdom_password}"
  api_key  = "${var.pingdom_api_key}"
}

resource "pingdom_team" "test" {
  name = "Team for testing"
}

resource "pingdom_user" "first_user" {
  username = "johndoe"
}

resource "pingdom_contact" "first_user_contact_email_1" {
  user_id        = "${pingdom_user.first_user.id}"
  email          = "john@doe.com"
  severity_level = "HIGH"
}

resource "pingdom_contact" "first_user_contact_email_2" {
  user_id        = "${pingdom_user.first_user.id}"
  email          = "john.doe@doe.com"
  severity_level = "LOW"
}

resource "pingdom_contact" "first_user_contact_sms_1" {
  user_id        = "${pingdom_user.first_user.id}"
  number         = "700000000"
  country_code   = "33"
  phone_provider = "nexmo"
  severity_level = "HIGH"
}

resource "pingdom_user" "second_user" {
  username = "janedoe"
}

resource "pingdom_contact" "second_user_contact_email_1" {
  user_id        = "${pingdom_user.second_user.id}"
  email          = "jane@doe.com"
  severity_level = "high"
}
