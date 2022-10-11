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

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 4.33.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_policy.consumers_send](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_route53_record.this_dkim_verification](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_record.this_domain_verification](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_ses_domain_dkim.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_domain_dkim) | resource |
| [aws_ses_domain_identity.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_domain_identity) | resource |
| [aws_ses_email_identity.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ses_email_identity) | resource |
| [aws_caller_identity.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_iam_policy_document.consumers_send](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_region.active](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |
| [aws_route53_zone.selected](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/route53_zone) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_environment"></a> [environment](#input\_environment) | To manage a resources with tags | `string` | n/a | yes |
| <a name="input_prefix"></a> [prefix](#input\_prefix) | The prefix name of customer to be displayed in AWS console and resource | `string` | n/a | yes |
| <a name="input_is_create_consumer_policy"></a> [is\_create\_consumer\_policy](#input\_is\_create\_consumer\_policy) | Whether to create consumer readonly policy | `bool` | `false` | no |
| <a name="input_ses_domain"></a> [ses\_domain](#input\_ses\_domain) | Domain that will be use as SES identity. | <pre>object({<br>    domain            = string<br>    route53_zone_name = string<br>    is_verify_domain  = bool<br>    is_verify_dkim    = bool<br>  })</pre> | <pre>{<br>  "domain": null,<br>  "is_verify_dkim": false,<br>  "is_verify_domain": false,<br>  "route53_zone_name": null<br>}</pre> | no |
| <a name="input_ses_email"></a> [ses\_email](#input\_ses\_email) | Email that will be use as SES identity. | <pre>object({<br>    email = string<br>  })</pre> | <pre>{<br>  "email": null<br>}</pre> | no |
| <a name="input_ses_mode"></a> [ses\_mode](#input\_ses\_mode) | Mode defines which method to verify identity for SES, which are email and domain. | `string` | `"domain"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Custom tags which can be passed on to the AWS resources. They should be key value pairs having distinct keys. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cosumer_policy_arn"></a> [cosumer\_policy\_arn](#output\_cosumer\_policy\_arn) | ARN of Consumer Policy |
| <a name="output_ses_dkim_tokens"></a> [ses\_dkim\_tokens](#output\_ses\_dkim\_tokens) | A list of DKIM Tokens which, when added to the DNS Domain as CNAME records, allows for receivers to verify that emails were indeed authorized by the domain owner. |
| <a name="output_ses_domain_identity_arn"></a> [ses\_domain\_identity\_arn](#output\_ses\_domain\_identity\_arn) | The ARN of the SES domain identity |
| <a name="output_ses_domain_identity_verification_token"></a> [ses\_domain\_identity\_verification\_token](#output\_ses\_domain\_identity\_verification\_token) | A code which when added to the domain as a TXT record will signal to SES that the owner of the domain has authorised SES to act on their behalf. The domain identity will be in state 'verification pending' until this is done. See below for an example of how this might be achieved when the domain is hosted in Route 53 and managed by Terraform. Find out more about verifying domains in Amazon SES in the AWS SES docs. |
| <a name="output_ses_email_identity_arn"></a> [ses\_email\_identity\_arn](#output\_ses\_email\_identity\_arn) | The ARN of the SES email identity |
<!-- END_TF_DOCS -->
