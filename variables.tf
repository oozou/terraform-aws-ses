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
