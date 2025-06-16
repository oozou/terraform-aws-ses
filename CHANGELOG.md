# Change Log

All notable changes to this module will be documented in this file.

## [1.2.0] - 2025-06-11

### Updated

- Upgraded AWS provider from `>= 4.0` to `~> 6.0` to use the latest 6.x.x version
- Updated all example configurations to use AWS provider `~> 6.0`
- Updated documentation to reflect new provider requirements

## [1.1.1] - 2025-06-11

### Updated

- Switch to data.aws_route53_zone.selected[0].id for easier to read and maintain

## [1.1.0] - 2025-06-10

### Added

- resource `aws_ses_domain_mail_from`
- resource `aws_route53_record.ses_domain_mail_from_mx` spf record
- resource `aws_route53_record.ses_domain_mail_from_txt` spf record
- resource `aws_route53_record.dmarc` dmarc record
- var `dmarc_record`

### Updated

- var `ses_domain`, add is_verify_dmarc property


## [1.0.2] - 2023-12-25

### Updated

- Update aws_route53_record condition

## [1.0.1] - 2022-10-04

### Added

- `is_create_consumer_policy` to enabled create the policy to grant permission for sending email
- output `cosumer_policy_arn`

## [1.0.0] - 2022-07-22

### Added

- init terraform-aws-ses module

### Noted

- ses configuration sets should be create along with ses identity, but terraform doesn't support this yet. (ref: https://github.com/hashicorp/terraform-provider-aws/issues/21129)
