# Create Route53 hosted zone for domain verification
resource "aws_route53_zone" "main" {
  name = var.route53_zone_name

  tags = merge(var.tags, {
    Name = "${var.prefix}-${var.environment}-${var.route53_zone_name}"
  })
}

module "ses" {
  source = "../.."

  prefix      = var.prefix
  environment = var.environment
  tags        = var.tags

  ses_mode   = "domain"
  ses_domain = var.ses_domain

  is_create_consumer_policy = true
}
