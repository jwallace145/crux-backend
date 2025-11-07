# =====================
# CruxProject S3 Bucket
# =====================

module "media_bucket" {
  source = "./modules/private-s3-bucket"

  project_name = var.project_name
  environment  = var.environment

  read_access = {
    api_read = {
      principal = module.api.task_role_arn
      prefixes  = ["users/"]
    }
  }

  write_access = {
    api_write = {
      principal = module.api.task_role_arn
      prefixes  = ["users/"]
    }
  }

  full_access_principals = [
    module.cicd.user_arn
  ]
}
