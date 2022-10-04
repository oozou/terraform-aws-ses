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
/*                                     SES                                    */
/* -------------------------------------------------------------------------- */

variable "ses_mode" {
  description = "Mode defines which method to verify identity for SES, which are email and domain."
  type        = string
  default     = "domain"
  validation {
    condition     = contains(["email", "domain"], var.ses_mode)
    error_message = "Allowed values: `email`, `domain`."
  }
}

variable "ses_email" {
  description = "Email that will be use as SES identity."
  type = object({
    email = string
  })
  default = {
    email = null
  }
}

variable "ses_domain" {
  description = "Domain that will be use as SES identity."
  type = object({
    domain            = string
    route53_zone_name = string
    is_verify_domain  = bool
    is_verify_dkim    = bool
  })
  default = {
    domain            = null
    is_verify_dkim    = false
    is_verify_domain  = false
    route53_zone_name = null
  }
}

variable "is_create_consumer_policy" {
  description = "Whether to create consumer readonly policy"
  type        = bool
  default     = false
}
