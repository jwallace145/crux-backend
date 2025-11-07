# ======================================================
# Continuous Integration & Continuous Deployment (CI/CD)
# ======================================================

module "cicd" {
  source = "./modules/cicd-user"

  project_name = var.project_name
  environment  = var.environment
}
