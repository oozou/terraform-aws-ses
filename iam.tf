data "aws_iam_policy_document" "consumers_send" {
  count = var.is_create_consumer_policy ? 1 : 0
  statement {
    sid    = ""
    effect = "Allow"
    actions = [
      "ses:SendRawEmail",
      "ses:SendEmail",
      "ses:SendBounce"

    ]
    resources = [
      try(aws_ses_domain_identity.this[0].arn, aws_ses_email_identity.this[0].arn)
    ]
  }
}

resource "aws_iam_policy" "consumers_send" {
  count  = var.is_create_consumer_policy ? 1 : 0
  name   = "${var.prefix}-AllowSESSend-policy"
  policy = data.aws_iam_policy_document.consumers_send[0].json

  tags = merge({ Name = "${var.prefix}-AllowSESSend-${data.aws_region.active.name}-policy" }, local.tags)
}
