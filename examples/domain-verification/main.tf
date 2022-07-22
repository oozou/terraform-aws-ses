module "ses" {
  source   = "../.."
  ses_mode = "domain"
  ses_domain = {
    domain            = "example.com"
    is_verify_dkim    = true
    is_verify_domain  = true
    route53_zone_name = "example.com"
  }
}
