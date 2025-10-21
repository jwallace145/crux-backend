#!/bin/bash

################################################################################
# Terraform Deployment Script for CruxBackend
################################################################################
# Usage:
#   ./deploy.sh <action> <environment>
#
# Actions:
#   plan     - Show what changes will be made
#   apply    - Apply infrastructure changes
#   destroy  - Destroy all infrastructure (DANGEROUS)
#   output   - Show output values
#   validate - Validate Terraform configuration
#
# Environments:
#   dev, stg, prod
#
# Examples:
#   ./deploy.sh plan dev
#   ./deploy.sh apply dev
#   ./deploy.sh destroy dev
################################################################################

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$INFRA_DIR")"
TERRAFORM_DIR="$INFRA_DIR/terraform"
ENVIRONMENTS_DIR="$INFRA_DIR/environments"
ENV_FILE="$PROJECT_ROOT/.env"

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "${BOLD}${BLUE}===================================================${NC}"
    echo -e "${BOLD}${BLUE}$1${NC}"
    echo -e "${BOLD}${BLUE}===================================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ Error: $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ Warning: $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

usage() {
    cat << EOF
${BOLD}Usage:${NC}
    $0 <action> <environment>

${BOLD}Actions:${NC}
    plan        Show what infrastructure changes will be made
    apply       Apply infrastructure changes
    destroy     Destroy all infrastructure (requires confirmation)
    output      Show Terraform output values
    validate    Validate Terraform configuration
    refresh     Refresh Terraform state

${BOLD}Environments:${NC}
    dev         Development environment
    stg         Staging environment
    prod        Production environment

${BOLD}Examples:${NC}
    $0 plan dev             # Plan changes for dev environment
    $0 apply dev            # Apply changes to dev environment
    $0 destroy dev          # Destroy dev infrastructure
    $0 output dev           # Show dev environment outputs

${BOLD}Environment Variables:${NC}
    ACCESS_TOKEN_SECRET_KEY           The secret key used for the API access token (required)
    REFRESH_TOKEN_SECRET_KEY          The secret key used for the API refresh token (required)
    DB_PASSWORD                       Database master password (required)
    AWS_PROFILE                       AWS profile to use (optional)
    AWS_REGION                        AWS region (optional, defaults to tfvars)

${BOLD}.env File:${NC}
    The script will automatically load environment variables from:
    $ENV_FILE

    And convert them to TF_VAR_* format for Terraform.

${BOLD}Backend Configuration:${NC}
    The script will look for backend configuration files at:
    $ENVIRONMENTS_DIR/<environment>.tfbackend

    These files configure remote state storage (S3 bucket, DynamoDB table, etc.)

EOF
    exit 1
}

load_env_file() {
    if [ -f "$ENV_FILE" ]; then
        print_info "Loading environment variables from .env file..."

        # Read .env file and export as TF_VAR_* variables
        while IFS='=' read -r key value || [ -n "$key" ]; do
            # Skip empty lines and comments
            [[ -z "$key" || "$key" =~ ^#.*$ ]] && continue

            # Remove leading/trailing whitespace
            key=$(echo "$key" | xargs)
            value=$(echo "$value" | xargs)

            # Remove quotes if present
            value="${value%\"}"
            value="${value#\"}"
            value="${value%\'}"
            value="${value#\'}"

            # Convert to lowercase and export as TF_VAR_*
            tf_var_name="TF_VAR_$(echo "$key" | tr '[:upper:]' '[:lower:]')"
            export "$tf_var_name=$value"
            print_success "Loaded: $tf_var_name"
        done < "$ENV_FILE"

        echo ""
    else
        print_warning ".env file not found at: $ENV_FILE"
        print_info "Continuing without .env file..."
        echo ""
    fi
}

check_dependencies() {
    local missing_deps=()

    if ! command -v terraform &> /dev/null; then
        missing_deps+=("terraform")
    fi

    if ! command -v aws &> /dev/null; then
        missing_deps+=("aws-cli")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        echo ""
        echo "Install instructions:"
        echo "  terraform: https://www.terraform.io/downloads"
        echo "  aws-cli: https://aws.amazon.com/cli/"
        exit 1
    fi

    print_success "All dependencies found"
}

check_aws_credentials() {
    print_info "Checking AWS credentials..."

    if ! aws sts get-caller-identity &> /dev/null; then
        print_error "AWS credentials not configured or invalid"
        echo ""
        echo "Configure AWS credentials using one of:"
        echo "  1. aws configure"
        echo "  2. Set AWS_PROFILE environment variable"
        echo "  3. Set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY"
        exit 1
    fi

    local caller_identity=$(aws sts get-caller-identity --output json)
    local account_id=$(echo "$caller_identity" | grep -o '"Account": "[^"]*' | cut -d'"' -f4)
    local user_arn=$(echo "$caller_identity" | grep -o '"Arn": "[^"]*' | cut -d'"' -f4)

    print_success "AWS credentials valid"
    print_info "Account: $account_id"
    print_info "Identity: $user_arn"
}

check_environment_file() {
    local env=$1
    local tfvars_file="$ENVIRONMENTS_DIR/${env}.tfvars"

    if [ ! -f "$tfvars_file" ]; then
        print_error "Environment file not found: $tfvars_file"
        echo ""
        echo "Available environments:"
        ls -1 "$ENVIRONMENTS_DIR"/*.tfvars 2>/dev/null | xargs -n1 basename | sed 's/.tfvars$//' || echo "  (none found)"
        exit 1
    fi

    print_success "Environment file found: ${env}.tfvars"
}

check_backend_file() {
    local env=$1
    local backend_file="$ENVIRONMENTS_DIR/${env}.tfbackend"

    if [ -f "$backend_file" ]; then
        print_success "Backend configuration found: ${env}.tfbackend"
        print_info "Will use remote state backend"
        return 0
    else
        print_warning "Backend configuration not found: ${env}.tfbackend"
        print_info "Will use local state (not recommended for production)"
        return 1
    fi
}

check_jwt_secrets() {
    local missing_secrets=()

    # Check for ACCESS_TOKEN_SECRET_KEY (becomes TF_VAR_access_token_secret_key via load_env_file)
    if [ -z "$TF_VAR_access_token_secret_key" ]; then
        missing_secrets+=("ACCESS_TOKEN_SECRET_KEY")
    fi

    # Check for REFRESH_TOKEN_SECRET_KEY (becomes TF_VAR_refresh_token_secret_key via load_env_file)
    if [ -z "$TF_VAR_refresh_token_secret_key" ]; then
        missing_secrets+=("REFRESH_TOKEN_SECRET_KEY")
    fi

    if [ ${#missing_secrets[@]} -ne 0 ]; then
        print_error "Missing required JWT secret environment variables: ${missing_secrets[*]}"
        echo ""
        echo "These variables must be set in your .env file:"
        for secret in "${missing_secrets[@]}"; do
            echo "  $secret=<your-secret-value>"
        done
        exit 1
    fi

    print_success "JWT secrets configured"
    print_info "  ACCESS_TOKEN_SECRET_KEY: ****${TF_VAR_access_token_secret_key: -4}"
    print_info "  REFRESH_TOKEN_SECRET_KEY: ****${TF_VAR_refresh_token_secret_key: -4}"
}

check_database_password() {
    # Check if db password was loaded from .env file first
    if [ -n "$TF_VAR_db_password" ]; then
        print_success "Database password loaded from .env file"
        # Also set the legacy variable name for backward compatibility
        export TF_VAR_database_password="$TF_VAR_db_password"
        return
    fi

    if [ -z "$TF_VAR_database_password" ]; then
        print_warning "Database password not set"
        echo ""
        read -sp "Enter database password: " db_password
        echo ""
        export TF_VAR_database_password="$db_password"
    fi
    print_success "Database password configured"
}

terraform_init() {
    local env=$1
    local backend_file="$ENVIRONMENTS_DIR/${env}.tfbackend"

    print_header "Initializing Terraform"
    cd "$TERRAFORM_DIR"

    # Check if backend configuration exists
    if [ -f "$backend_file" ]; then
        print_info "Initializing with backend configuration: ${env}.tfbackend"
        terraform init -upgrade -backend-config="$backend_file" -reconfigure
    else
        print_info "Initializing with local backend (no backend configuration found)"
        terraform init -upgrade
    fi

    print_success "Terraform initialized"
}

terraform_validate() {
    print_header "Validating Terraform Configuration"
    cd "$TERRAFORM_DIR"
    terraform validate
    print_success "Terraform configuration is valid"
}

terraform_plan() {
    local env=$1
    local tfvars_file="$ENVIRONMENTS_DIR/${env}.tfvars"

    print_header "Planning Infrastructure Changes ($env)"
    cd "$TERRAFORM_DIR"

    terraform plan \
        -var-file="$tfvars_file" \
        -out="$TERRAFORM_DIR/tfplan-${env}"

    print_success "Plan complete"
    print_info "Plan saved to: tfplan-${env}"
}

terraform_apply() {
    local env=$1
    local tfvars_file="$ENVIRONMENTS_DIR/${env}.tfvars"

    print_header "Applying Infrastructure Changes ($env)"

    # Check if plan file exists
    if [ -f "$TERRAFORM_DIR/tfplan-${env}" ]; then
        print_info "Using existing plan file: tfplan-${env}"
        echo ""
        read -p "Apply this plan? (yes/no): " confirm

        if [ "$confirm" != "yes" ]; then
            print_warning "Apply cancelled"
            exit 0
        fi

        cd "$TERRAFORM_DIR"
        terraform apply "tfplan-${env}"
    else
        print_warning "No plan file found, creating new plan..."
        echo ""
        read -p "Continue with apply? (yes/no): " confirm

        if [ "$confirm" != "yes" ]; then
            print_warning "Apply cancelled"
            exit 0
        fi

        cd "$TERRAFORM_DIR"
        terraform apply -var-file="$tfvars_file"
    fi

    print_success "Infrastructure applied successfully"
}

terraform_destroy() {
    local env=$1
    local tfvars_file="$ENVIRONMENTS_DIR/${env}.tfvars"

    print_header "Destroying Infrastructure ($env)"

    echo -e "${RED}${BOLD}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║                    ⚠️  WARNING  ⚠️                          ║"
    echo "║                                                            ║"
    echo "║  This will DESTROY all infrastructure in the $env       ║"
    echo "║  environment including:                                    ║"
    echo "║    • RDS Database (ALL DATA WILL BE LOST)                 ║"
    echo "║    • VPC and all networking resources                     ║"
    echo "║    • Security groups and network ACLs                     ║"
    echo "║                                                            ║"
    echo "║  This action CANNOT be undone!                            ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo ""

    read -p "Type 'destroy-${env}' to confirm: " confirm

    if [ "$confirm" != "destroy-${env}" ]; then
        print_warning "Destroy cancelled"
        exit 0
    fi

    echo ""
    read -p "Are you absolutely sure? (yes/no): " final_confirm

    if [ "$final_confirm" != "yes" ]; then
        print_warning "Destroy cancelled"
        exit 0
    fi

    cd "$TERRAFORM_DIR"
    terraform destroy -var-file="$tfvars_file"

    print_success "Infrastructure destroyed"
}

terraform_output() {
    local env=$1

    print_header "Terraform Outputs ($env)"
    cd "$TERRAFORM_DIR"
    terraform output
}

terraform_refresh() {
    local env=$1
    local tfvars_file="$ENVIRONMENTS_DIR/${env}.tfvars"

    print_header "Refreshing Terraform State ($env)"
    cd "$TERRAFORM_DIR"
    terraform refresh -var-file="$tfvars_file"
    print_success "State refreshed"
}

################################################################################
# Main Script
################################################################################

main() {
    # Check arguments
    if [ $# -lt 2 ]; then
        usage
    fi

    local action=$1
    local environment=$2

    # Validate action
    case $action in
        plan|apply|destroy|output|validate|refresh)
            ;;
        *)
            print_error "Invalid action: $action"
            echo ""
            usage
            ;;
    esac

    # Validate environment
    case $environment in
        dev|stg|prod)
            ;;
        *)
            print_error "Invalid environment: $environment"
            echo ""
            usage
            ;;
    esac

    print_header "CruxBackend Terraform Deployment"
    echo ""
    echo "Action:      $action"
    echo "Environment: $environment"
    echo ""

    # Load .env file early
    load_env_file

    # Pre-flight checks
    check_dependencies
    check_aws_credentials
    check_environment_file "$environment"
    check_backend_file "$environment"

    # Only check secrets for actions that need them
    if [ "$action" != "output" ] && [ "$action" != "validate" ]; then
        check_jwt_secrets
        check_database_password
    fi

    # Initialize with environment-specific backend
    terraform_init "$environment"

    # Execute action
    case $action in
        validate)
            terraform_validate
            ;;
        plan)
            terraform_validate
            terraform_plan "$environment"
            ;;
        apply)
            terraform_apply "$environment"
            terraform_output "$environment"
            ;;
        destroy)
            terraform_destroy "$environment"
            ;;
        output)
            terraform_output "$environment"
            ;;
        refresh)
            terraform_refresh "$environment"
            ;;
    esac

    echo ""
    print_success "Operation completed successfully"
}

# Run main function
main "$@"
