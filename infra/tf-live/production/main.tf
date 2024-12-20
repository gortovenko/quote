terraform {
  required_version = ">= 1.0.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}



module "vpc" {
  source               = "git::https://github.com/gortovenko/tf-modules.git//vpc"
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}

module "security_groups" {
  source         = "git::https://github.com/gortovenko/tf-modules.git//security_groups"
  vpc_id         = module.vpc.vpc_id
  redis_port     = 6379
  lambda_sg_name = "lambda_sg"
  redis_sg_name  = "redis_sg"
  ecs_sg_name    = "ecs_sg"
}

module "elasticache" {
  source         = "git::https://github.com/gortovenko/tf-modules.git//elasticache"
  engine_version = var.redis_engine_version
  cluster_id     = var.redis_cluster_id
  subnet_ids     = module.vpc.private_subnets
  sg_ids         = [module.security_groups.redis_sg_id]
}

module "iam" {
  source             = "git::https://github.com/gortovenko/tf-modules.git//iam"
  lambda_name        = var.lambda_function_name
  secrets_arn        = module.secrets_manager.secrets_arn
  ecs_task_role_name = var.ecs_task_role_name
}

module "lambda" {
  source                 = "git::https://github.com/gortovenko/tf-modules.git//lambda"
  function_name          = var.lambda_function_name
  handler                = var.lambda_handler
  runtime                = var.lambda_runtime
  role_arn               = module.iam.lambda_role_arn
  vpc_subnet_ids         = module.vpc.private_subnets
  vpc_security_group_ids = [module.security_groups.lambda_sg_id]
  source_path            = var.source_path
  filename               = var.filename

  environment_variables = {
    REDIS_ENDPOINT = module.elasticache.primary_endpoint_address
    REDIS_PORT     = tostring(module.elasticache.primary_endpoint_port)
    SECRETS_ARN    = module.secrets_manager.secrets_arn
  }
  s3_bucket = module.s3.bucket_name
  s3_key    = "example.zip"

  depends_on = [module.secrets_manager]


}


module "api_gateway" {
  source     = "git::https://github.com/gortovenko/tf-modules.git//api_gateway"
  api_name   = var.api_name
  lambda_arn = module.lambda.lambda_arn
  stage_name = var.api_stage_name
}

module "ecs" {
  source                = "git::https://github.com/gortovenko/tf-modules.git//ecs"
  cluster_name          = var.ecs_cluster_name
  service_name          = var.ecs_service_name
  task_family           = var.ecs_task_family
  container_name        = var.ecs_container_name
  container_image       = var.ecs_container_image
  container_port        = var.ecs_container_port
  vpc_id                = module.vpc.vpc_id
  public_subnets        = module.vpc.public_subnet_ids
  ecs_security_group_id = module.security_groups.ecs_sg_id
  secrets_arn           = module.secrets_manager.secrets_arn


  depends_on = [module.secrets_manager]
}

module "ecr" {
  source          = "git::https://github.com/gortovenko/tf-modules.git//ecr"
  repository_name = var.repository_name
  environment     = var.environment
}


module "secrets_manager" {
  source          = "git::https://github.com/gortovenko/tf-modules.git//secrets_manager"
  secrets_name    = var.secrets_name
  cache_provider  = var.cache_provider
  run_local_mode  = var.run_local_mode
  elasticache_url = module.elasticache.primary_endpoint_address
  server_address  = var.server_address
  base_url        = var.base_url
  default_count   = var.default_count
  environment     = var.environment
}

module "s3" {
  source         = "git::https://github.com/gortovenko/tf-modules.git//s3"
  bucket_name    = var.bucket_name
  s3_bucket_name = var.s3_bucket_name
  environment    = var.environment

}