# AWS Secrets Manager Secret Module

This Terraform module creates an AWS Secrets Manager secret that stores multiple key-value pairs as a JSON object.

## Features

- ✅ Store multiple secrets as a single JSON object
- ✅ Optional KMS encryption with custom keys
- ✅ Configurable recovery window for deleted secrets
- ✅ Optional automatic rotation support
- ✅ Comprehensive validation and tagging
- ✅ Sensitive output handling

## Usage

### Basic Example

```hcl
module "jwt_secrets" {
  source = "./modules/secrets-manager-secret"

  secret_name        = "crux-api/jwt-secrets"
  secret_description = "JWT token secrets for Crux API"

  secrets = {
    ACCESS_TOKEN_SECRET_KEY  = var.access_token_secret_key
    REFRESH_TOKEN_SECRET_KEY = var.refresh_token_secret_key
  }

  tags = {
    Environment = "dev"
    Application = "crux-api"
  }
}
```

### Advanced Example with Rotation

```hcl
module "database_credentials" {
  source = "./modules/secrets-manager-secret"

  secret_name        = "crux-api/database-credentials"
  secret_description = "Database credentials for Crux API"

  secrets = {
    username = "cruxadmin"
    password = var.database_password
    host     = "cruxdb-dev.us-east-1.rds.amazonaws.com"
    port     = "5432"
    database = "cruxdb"
  }

  # Custom KMS key for encryption
  kms_key_id = aws_kms_key.secrets.id

  # Enable automatic rotation
  enable_rotation     = true
  rotation_lambda_arn = aws_lambda_function.rotate_db_secret.arn
  rotation_days       = 30

  # Longer recovery window for production
  recovery_window_in_days = 30

  tags = {
    Environment = "prod"
    Application = "crux-api"
    Sensitive   = "true"
  }
}
```

### Using the Secret in ECS Task Definitions

```hcl
# Reference the secret ARN in ECS task definition
resource "aws_ecs_task_definition" "api" {
  # ... other configuration ...

  container_definitions = jsonencode([{
    name  = "crux-api"
    image = var.container_image

    secrets = [
      {
        name      = "ACCESS_TOKEN_SECRET_KEY"
        valueFrom = "${module.jwt_secrets.secret_arn}:ACCESS_TOKEN_SECRET_KEY::"
      },
      {
        name      = "REFRESH_TOKEN_SECRET_KEY"
        valueFrom = "${module.jwt_secrets.secret_arn}:REFRESH_TOKEN_SECRET_KEY::"
      }
    ]
  }])
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| secret_name | The name of the secret in AWS Secrets Manager | `string` | n/a | yes |
| secrets | Map of secret key-value pairs to store as JSON | `map(string)` | n/a | yes |
| secret_description | Description of the secret | `string` | `""` | no |
| recovery_window_in_days | Number of days to retain the secret before permanent deletion | `number` | `30` | no |
| tags | Tags to apply to the secret | `map(string)` | `{}` | no |
| kms_key_id | ARN or ID of the AWS KMS key to encrypt the secret | `string` | `null` | no |
| enable_rotation | Enable automatic rotation for this secret | `bool` | `false` | no |
| rotation_lambda_arn | ARN of the Lambda function that can rotate the secret | `string` | `null` | no |
| rotation_days | Number of days between automatic rotations | `number` | `30` | no |

## Outputs

| Name | Description | Sensitive |
|------|-------------|:---------:|
| secret_id | The ID of the secret (same as ARN) | no |
| secret_arn | The ARN of the secret | no |
| secret_name | The name of the secret | no |
| secret_version_id | The version ID of the secret value | no |
| secret_json | The secret stored as JSON | yes |

## JSON Format

The module automatically converts the `secrets` map to JSON. For example:

**Input:**
```hcl
secrets = {
  ACCESS_TOKEN_SECRET_KEY  = "my-access-secret"
  REFRESH_TOKEN_SECRET_KEY = "my-refresh-secret"
}
```

**Stored in AWS Secrets Manager:**
```json
{
  "ACCESS_TOKEN_SECRET_KEY": "my-access-secret",
  "REFRESH_TOKEN_SECRET_KEY": "my-refresh-secret"
}
```

## Security Best Practices

1. **Never commit secrets to version control**: Use variables and pass secrets via environment variables or secure parameter stores
2. **Enable KMS encryption**: Use customer-managed KMS keys for sensitive secrets
3. **Configure recovery window**: Set appropriate recovery window for production (30 days recommended)
4. **Use rotation**: Enable automatic rotation for database credentials and API keys
5. **Apply least privilege**: Ensure only necessary IAM roles/users can access the secret
6. **Tag appropriately**: Use tags to identify sensitive secrets and their purpose

## IAM Permissions Required

### For ECS Tasks to Read Secrets

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": [
        "arn:aws:secretsmanager:region:account-id:secret:crux-api/*"
      ]
    }
  ]
}
```

### For Terraform to Manage Secrets

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:CreateSecret",
        "secretsmanager:UpdateSecret",
        "secretsmanager:DeleteSecret",
        "secretsmanager:DescribeSecret",
        "secretsmanager:PutSecretValue",
        "secretsmanager:GetSecretValue",
        "secretsmanager:TagResource",
        "secretsmanager:UntagResource"
      ],
      "Resource": "*"
    }
  ]
}
```

## Retrieving Secrets

### AWS CLI

```bash
# Get the entire JSON object
aws secretsmanager get-secret-value \
  --secret-id crux-api/jwt-secrets \
  --query SecretString \
  --output text | jq .

# Get a specific key
aws secretsmanager get-secret-value \
  --secret-id crux-api/jwt-secrets \
  --query SecretString \
  --output text | jq -r .ACCESS_TOKEN_SECRET_KEY
```

### In Application Code (Go)

```go
import (
    "encoding/json"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type JWTSecrets struct {
    AccessTokenSecretKey  string `json:"ACCESS_TOKEN_SECRET_KEY"`
    RefreshTokenSecretKey string `json:"REFRESH_TOKEN_SECRET_KEY"`
}

func getJWTSecrets() (*JWTSecrets, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)

    result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
        SecretId: aws.String("crux-api/jwt-secrets"),
    })
    if err != nil {
        return nil, err
    }

    var secrets JWTSecrets
    err = json.Unmarshal([]byte(*result.SecretString), &secrets)
    return &secrets, err
}
```

## Notes

- **Recovery Window**: Setting `recovery_window_in_days = 0` deletes the secret immediately without recovery option
- **Cost**: AWS Secrets Manager charges $0.40 per secret per month + $0.05 per 10,000 API calls
- **Rotation**: Automatic rotation requires a Lambda function that implements the rotation logic
- **Versioning**: AWS Secrets Manager automatically versions secrets; this module creates a new version on each update
