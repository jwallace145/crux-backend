module "db_secrets" {
  source = "./modules/secrets-manager-secret"

  secret_name = "${var.service_name}-api-${var.environment}/db-user-secrets"

  secrets = {
    DB_USER     = var.db_user
    DB_PASSWORD = var.db_password
  }
}

module "jwt_secrets" {
  source = "./modules/secrets-manager-secret"

  secret_name = "${var.service_name}-api-${var.environment}/jwt-secrets"

  secrets = {
    ACCESS_TOKEN_SECRET_KEY  = var.access_token_secret_key
    REFRESH_TOKEN_SECRET_KEY = var.refresh_token_secret_key
  }
}
