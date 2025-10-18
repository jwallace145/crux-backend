module "crux_api_ecr" {
  source = "./modules/ecr-repository"

  repository_name      = "crux-api"
  image_tag_mutability = "MUTABLE"
  scan_on_push         = true
  encryption_type      = "AES256"

  # Lifecycle policy settings
  lifecycle_policy_enabled = true
  image_count_to_keep      = 10
  tag_prefixes_to_keep     = ["dev", "prod", "staging", "latest"]
  untagged_days_to_keep    = 7
}
