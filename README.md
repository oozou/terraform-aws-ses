# terraform-aws-ses

Terraform module use to create AWS SES.

## Usage

- **Domain verification**

```terraform
module "ses" {
  source = "git::ssh://git@github.com:oozou/terraform-aws-ses.git?ref=main"
  ses_mode = "domain"
  ses_domain = {
    domain            = "<domain>"
    is_verify_dkim    = true
    is_verify_domain  = true
    route53_zone_name = "<route53-zone-name>"
  }
}
```

- **Email verification**

```terraform
module "ses" {
  source = "git::ssh://git@github.com:oozou/terraform-aws-ses.git?ref=main"
  ses_mode = "email"
  ses_email = {
    email = "<email>"
  }
}
```

<!-- BEGIN_TF_DOCS -->

## Requirements

| Name                                                                     | Version  |
| ------------------------------------------------------------------------ | -------- |
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 1.0.0 |
| <a name="requirement_aws"></a> [aws](#requirement_aws)                   | >= 4.0   |

## Providers

| Name                                             | Version |
| ------------------------------------------------ | ------- |
| <a name="provider_aws"></a> [aws](#provider_aws) | 4.23.0  |

## Modules

No modules.

## Resources

| Name                                                                                                                                      | Type        |
| ----------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| [aws_route53_record.this_dkim_verification](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record)   | resource    |
| [aws_route53_record.this_domain_verification](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource    |
| [aws_ses_domain_dkim.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_domain_dkim)                   | resource    |
| [aws_ses_domain_identity.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_domain_identity)           | resource    |
| [aws_ses_email_identity.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_email_identity)             | resource    |
| [aws_route53_zone.selected](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone)                  | data source |

## Inputs

| Name                                                            | Description                                                                       | Type                                                                                                                                | Default                                                                                                                         | Required |
| --------------------------------------------------------------- | --------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | :------: |
| <a name="input_ses_domain"></a> [ses_domain](#input_ses_domain) | Domain that will be use as SES identity.                                          | <pre>object({<br> domain = string<br> route53_zone_name = string<br> is_verify_domain = bool<br> is_verify_dkim = bool<br> })</pre> | <pre>{<br> "domain": null,<br> "is_verify_dkim": false,<br> "is_verify_domain": false,<br> "route53_zone_name": null<br>}</pre> |    no    |
| <a name="input_ses_email"></a> [ses_email](#input_ses_email)    | Email that will be use as SES identity.                                           | <pre>object({<br> email = string<br> })</pre>                                                                                       | <pre>{<br> "email": null<br>}</pre>                                                                                             |    no    |
| <a name="input_ses_mode"></a> [ses_mode](#input_ses_mode)       | Mode defines which method to verify identity for SES, which are email and domain. | `string`                                                                                                                            | `"email"`                                                                                                                       |    no    |

## Outputs

| Name                                                                                                                                                  | Description                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| <a name="output_ses_dkim_tokens"></a> [ses_dkim_tokens](#output_ses_dkim_tokens)                                                                      | A list of DKIM Tokens which, when added to the DNS Domain as CNAME records, allows for receivers to verify that emails were indeed authorized by the domain owner.                                                                                                                                                                                                                                                               |
| <a name="output_ses_domain_identity_arn"></a> [ses_domain_identity_arn](#output_ses_domain_identity_arn)                                              | The ARN of the SES domain identity                                                                                                                                                                                                                                                                                                                                                                                               |
| <a name="output_ses_domain_identity_verification_token"></a> [ses_domain_identity_verification_token](#output_ses_domain_identity_verification_token) | A code which when added to the domain as a TXT record will signal to SES that the owner of the domain has authorised SES to act on their behalf. The domain identity will be in state 'verification pending' until this is done. See below for an example of how this might be achieved when the domain is hosted in Route 53 and managed by Terraform. Find out more about verifying domains in Amazon SES in the AWS SES docs. |
| <a name="output_ses_email_identity_arn"></a> [ses_email_identity_arn](#output_ses_email_identity_arn)                                                 | The ARN of the SES email identity                                                                                                                                                                                                                                                                                                                                                                                                |

<!-- END_TF_DOCS -->
