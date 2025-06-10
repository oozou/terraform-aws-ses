module "ses" {
  source = "../.."

  prefix      = var.prefix
  environment = var.environment
  tags        = var.tags

  ses_mode = "domain"
  ses_domain = {
    domain            = "mail.domain.com"
    is_verify_dkim    = true
    is_verify_domain  = true
    route53_zone_name = "domain.com"
    is_verify_dmarc   = true
  }

  is_create_consumer_policy = true
}
