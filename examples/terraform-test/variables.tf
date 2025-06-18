/* -------------------------------------------------------------------------- */
/*                                   Generic                                  */
/* -------------------------------------------------------------------------- */

variable "prefix" {
  description = "The prefix name of customer to be displayed in AWS console and resource"
  type        = string
}

variable "environment" {
  description = "To manage a resources with tags"
  type        = string
}

variable "tags" {
  description = "Custom tags which can be passed on to the AWS resources. They should be key value pairs having distinct keys."
  type        = map(string)
  default     = {}
}

/* -------------------------------------------------------------------------- */
/*                                    SES                                     */
/* -------------------------------------------------------------------------- */

variable "ses_domain" {
  description = "SES domain configuration"
  type = object({
    domain            = string
    is_verify_dkim    = bool
    is_verify_domain  = bool
    route53_zone_name = string
    is_verify_dmarc   = bool
  })
  default = {
    domain            = "mail.domain.com"
    is_verify_dkim    = true
    is_verify_domain  = true
    route53_zone_name = "domain.com"
    is_verify_dmarc   = true
  }
}

variable "route53_zone_name" {
  description = "Route53 zone name to create hosted zone for domain verification"
  type        = string
  default     = "domain.com"
}

variable "ses_email" {
  description = "SES email configuration"
  type = object({
    email = string
  })
  default = {
    email = "test@example.com"
  }
}