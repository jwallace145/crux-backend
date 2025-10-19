module "jwt_secrets" {
  source = "./modules/secrets-manager-secret"

  secret_name = "${var.service_name}-api-${var.environment}/jwt-secrets"

  secrets = {
    ACCESS_TOKEN_SECRET_KEY  = var.access_token_secret_key
    REFRESH_TOKEN_SECRET_KEY = var.refresh_token_secret_key
  }
}
