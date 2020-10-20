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

resource "pingdom_team" "test_one" {
  name = "Team 1 (updated by Terraform)"
}

provider "pingdom" {
  api_token  = "${var.pingdom_api_token}"
}

resource "pingdom_team" "test" {
  name = "Team 2 (updated by Terraform) with contacts"

  member_ids = [pingdom_contact.first_contact.id, pingdom_contact.second_contact.id]
}

resource "pingdom_contact" "first_contact" {
  name = "johndoe"

  sms_notification {
    number   = "5555555555"
    severity = "HIGH"
  }
  sms_notification {
    number   = "2222222222"
    severity = "LOW"
  }
}

resource "pingdom_contact" "second_contact" {
  name   = "janedoe"
  paused = true

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

data "pingdom_contact" "data_contact" {
  name = "janedoe"
}

output "test" {
  value = data.pingdom_contact.data_contact
}
