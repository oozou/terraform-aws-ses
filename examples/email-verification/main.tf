module "ses" {
  source   = "../.."
  ses_mode = "email"
  ses_email = {
    email = "test@example.com"
  }
}
