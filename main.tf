data "aws_route53_zone" "selected" {
  count = var.ses_mode == "domain" && var.ses_domain.is_verify_domain ? 1 : 0
  name  = var.ses_domain.route53_zone_name
}

resource "aws_ses_email_identity" "this" {
  count = var.ses_mode == "email" && var.ses_email.email != null ? 1 : 0
  email = var.ses_email.email
}

resource "aws_ses_domain_identity" "this" {
  count  = var.ses_mode == "domain" && var.ses_domain.domain != null ? 1 : 0
  domain = var.ses_domain.domain
}

resource "aws_route53_record" "this_domain_verification" {
  count   = var.ses_mode == "domain" && var.ses_domain.is_verify_domain ? 1 : 0
  zone_id = join("", data.aws_route53_zone.selected.*.id)
  name    = "_amazonses.${var.ses_domain.domain}"
  type    = "TXT"
  ttl     = "600"
  records = [join("", aws_ses_domain_identity.this.*.verification_token)]
}

resource "aws_ses_domain_dkim" "this" {
  count  = var.ses_mode == "domain" && var.ses_domain.is_verify_dkim ? 1 : 0
  domain = join("", aws_ses_domain_identity.this.*.domain)
}

resource "aws_route53_record" "this_dkim_verification" {
  count   = var.ses_mode == "domain" && var.ses_domain.is_verify_domain && var.ses_domain.is_verify_dkim ? 3 : 0
  zone_id = join("", data.aws_route53_zone.selected.*.id)
  name    = "${element(aws_ses_domain_dkim.this[0].dkim_tokens, count.index)}._domainkey.${var.ses_domain.domain}"
  type    = "CNAME"
  ttl     = "600"
  records = ["${element(aws_ses_domain_dkim.this[0].dkim_tokens, count.index)}.dkim.amazonses.com"]
}
