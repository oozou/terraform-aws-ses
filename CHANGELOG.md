# Change Log

All notable changes to this module will be documented in this file.

## [1.0.1] - 2022-10-04

### Added

- `is_create_consumer_policy` to enabled create the policy to grant permission for sending email
- output `cosumer_policy_arn`

## [1.0.0] - 2022-07-22

### Added

- init terraform-aws-ses module

### Noted

- ses configuration sets should be create along with ses identity, but terraform doesn't support this yet. (ref: https://github.com/hashicorp/terraform-provider-aws/issues/21129)
