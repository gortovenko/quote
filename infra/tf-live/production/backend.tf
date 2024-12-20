terraform {
  backend "s3" {
    bucket         = "quote-my"
    key            = "infra/tf-live/production/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "quote"
    acl            = "private"
  }
}
