output "api_invoke_url" {
  value = module.api_gateway.invoke_url
  description = "API Invoke URL"
}

output "public_subnet_ids" {
  value = module.vpc.public_subnet_ids
  description = "Public Subnet IDs"
}

output "elasticache_endpoint" {
  value       = module.elasticache.primary_endpoint_address
  description = "Endpoint for the Redis ElastiCache cluster"
  sensitive   = true
}

output "lambda_arn" {
  value       = module.lambda.lambda_arn
  description = "ARN of the Lambda function"
  sensitive   = true
}

output "secrets_arn" {
  value       = module.secrets_manager.secrets_arn
  description = "ARN of the secret stored in AWS Secrets Manager"
}

output "ecs_alb_dns_name" {
  value = module.ecs.alb_dns_name
  description = "DNS name of the Application Load Balancer"
}

output "primary_endpoint_address" {
  value       = module.elasticache.primary_endpoint_address
  description = "Primary endpoint address of the ElastiCache cluster"
  sensitive   = true
}

output "s3_bucket_id" {
  description = "The ID of the created S3 bucket"
  value       = module.s3.bucket_id
}

output "s3_bucket_arn" {
  description = "The ARN of the created S3 bucket"
  value       = module.s3.bucket_arn
}

output "s3_bucket_name" {
  description = "The name of the created S3 bucket"
  value       = module.s3.bucket_name
}
