variable "region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  type    = list(string)
  default = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  type    = list(string)
  default = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "redis_engine_version" {
  type    = string
  default = "6.x"
}

variable "redis_cluster_id" {
  type    = string
  default = "my-redis-cache"
}

variable "lambda_function_name" {
  type    = string
  default = "my-cache-lambda"
}

variable "lambda_handler" {
  type    = string
  default = "main"
}

variable "lambda_runtime" {
  type    = string
  default = "provided.al2"
}

variable "lambda_source_path" {
  type    = string
  default = ""
}


variable "api_name" {
  type    = string
  default = "my-api"
}

variable "api_stage_name" {
  type    = string
  default = "dev"
}

# ECS Variables
variable "ecs_cluster_name" {
  type    = string
  default = "my-ecs-cluster"
}

variable "ecs_service_name" {
  type    = string
  default = "my-ecs-service"
}

variable "ecs_task_family" {
  type    = string
  default = "my-task-family"
}

variable "ecs_container_name" {
  type    = string
  description = "Name of the container"
  default = "my-app-container"
}

variable "ecs_container_image" {
  type = string
  description = "The Docker image to use for the container"
  default = "quote-repo/quote:latest"
}

variable "ecs_container_port" {
  description = "The port the container is listening on"
  type = number

}

variable "run_local_mode" {
  description = "Run the ECS service in local mode"
  type = bool
}

variable "server_address" {
  type        = string
  description = "Server address for the application"
  default     = null
}

variable "base_url" {
  type        = string
  description = "Base URL for the scraper"
  default     = null
}

variable "repository_name" {
  type        = string
  description = "Name of the ECR repository"
}

variable "environment" {
  type        = string
  description = "Environment for tagging or organizing resources"
}

variable "secrets_name" {
  type        = string
  description = "Name of the secret in Secrets Manager"
}

variable "cache_provider" {
  type        = string
  description = "Cache provider type (e.g., elasticache)"
}

variable "default_count" {
  type        = number
  description = "Default count for the scraper"
}

variable "ecs_task_role_name" {
  type = string

}

variable "source_path" {
  type    = string
  default = "./infra/tf-live/production/function.zip"
}

variable "filename" {
  type    = string
  default = "./infra/tf-live/production/function.zip"
}

variable "bucket_name" {
  description = "The name of the S3 bucket"
  type        = string
}

variable "s3_bucket_name" {
  description = "Descriptive name for S3 bucket tags"
  type        = string
}