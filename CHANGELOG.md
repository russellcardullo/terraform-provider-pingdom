## Unreleased

## 1.1.2 (September 13, 2020)

NEW FEATURES:

  * Add support for managing teams, contacts, users (#36)
  * Allow adding users to teams (#61)

IMPROVEMENTS:

  * Add responsetimethreshold to checks (#36)
  * CI improvements (#48)
  * Documentation improvements (#45)
  * Uses latest patch version of go in CI builds (#50)
  * Update to terraform 0.12.18 (#51, #54)
  * Migrate to terraform-plugin-sdk (#65)
  * Sort tags on write to prevent unnecessary diffs (#58)
  * Documentation improvements (#53, #66, #67)
  * Use GitHub actions for builds and releases (#72)

BUG FIXES:
  * Include existing probefilter values on reads (#47)
  * Fix issue importing contacts (#60)

## 1.1.1 (October 5, 2019)

IMPROVEMENTS:

  * Add sensitive flags for secret values (#44)
  * Publish releases from Travis CI (#41)

## 1.0.0 (July 3, 2019)

NEW FEATURES:

  * Add TCP Checks (#21)
  * Add support for Public Reports (#21)
  * Add support for integrations/webhooks in checks (#14)
  * Add support for probe filters (#13)
  * Checks support paused parameter (#22)
  * Add support for tags on checks (#8)
  * Add support for importing existing checks (#9)
  * Add contact IDs to schema (#3)
  * Add option to use legacy notifications
  * Add spport for multi-user accounts (#1)

IMPROVEMENTS:

  * Use go modules to manage dependencies (#30)
  * Update to go-pingdom v1.0.0
  * Update to terraform 0.12.3 (#38)
  * Documentation updates (#12)

BUG FIXES:

  * Cannot use imported resource without forced re-creation (#31)
  * Fix teams response (#27)
  * Stop using `id` check schema in (#11)
  * Add default value for check URLs (#4)

## 0.2.0 (October 17, 2014)

IMPROVEMENTS:

  * Add support for Terraform 0.3.0

## 0.1.1 (September 16, 2014)

FEATURES:

  * Add support for managing ping type checks

BUG FIXES:

  * Don't override check resource values on create

## 0.1.0 (September 7, 2014)

  * Initial release
