provider "aws" {
  access_key                  = "mock_access_key"
  secret_key                  = "mock_secret_key"
  region                      = "us-east-1"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    sns = "http://localstack:4566"
    sqs = "http://localstack:4566"
    ses = "http://localstack:4566"
  }
}

resource "aws_sns_topic" "notifications" {
  name = "notifications"
}

resource "aws_sqs_queue" "sms_notifications" {
  name = "sms_notifications"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.sms_notifications_dlq.arn
    maxReceiveCount     = 5
  })
}
resource "aws_sqs_queue" "sms_notifications_dlq" {
  name = "sms_notifications_dlq"
}

resource "aws_sns_topic_subscription" "sms_notifications" {
  topic_arn = aws_sns_topic.notifications.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.sms_notifications.arn
}

resource "aws_sqs_queue" "email_notifications" {
  name = "email_notifications"

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.email_notifications_dlq.arn
    maxReceiveCount     = 5
  })
}
resource "aws_sqs_queue" "email_notifications_dlq" {
  name = "email_notifications_dlq"
}

resource "aws_sns_topic_subscription" "email_notifications" {
  topic_arn = aws_sns_topic.notifications.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.email_notifications.arn
}

resource "aws_ses_email_identity" "sender_email" {
  email = "sender@notifications.com"
}
