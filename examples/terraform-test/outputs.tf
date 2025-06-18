output "route53_zone_id" {
  description = "The Route53 hosted zone ID"
  value       = aws_route53_zone.main.zone_id
}

output "name_servers" {
  description = "The name servers for the hosted zone"
  value       = aws_route53_zone.main.name_servers
}

# SES Domain Module Outputs
output "ses_domain_identity_arn" {
  description = "The ARN of the SES domain identity"
  value       = module.ses.ses_domain_identity_arn
}

output "ses_domain_identity_verification_token" {
  description = "A code which when added to the domain as a TXT record will signal to SES that the owner of the domain has authorised SES to act on their behalf"
  value       = module.ses.ses_domain_identity_verification_token
}

output "ses_dkim_tokens" {
  description = "A list of DKIM Tokens which, when added to the DNS Domain as CNAME records, allows for receivers to verify that emails were indeed authorized by the domain owner"
  value       = module.ses.ses_dkim_tokens
}

output "cosumer_policy_arn" {
  description = "ARN of Consumer Policy"
  value       = module.ses.cosumer_policy_arn
}

# SES Email Module Outputs
output "ses_email_identity_arn" {
  description = "The ARN of the SES email identity"
  value       = module.ses_email.ses_email_identity_arn
}
