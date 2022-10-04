module "ses" {
  source = "../.."

  prefix      = var.prefix
  environment = var.environment
  tags        = var.tags

  ses_mode = "email"
  ses_email = {
    email = "test@example.com"
  }
}
