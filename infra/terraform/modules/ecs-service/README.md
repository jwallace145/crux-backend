# ECS Service Module

This Terraform module creates an ECS Fargate cluster and service for running the Crux Backend API containers. It includes task definitions, IAM roles, security groups, and CloudWatch logging.

## Features

- **ECS Fargate Cluster** with Container Insights support
- **ECS Service** with configurable desired count
- **Task Definition** with health checks and environment variables
- **IAM Roles** for task execution and task runtime
- **Security Groups** configured to accept traffic from ALB only
- **CloudWatch Logs** with configurable retention
- **ECS Exec** support for debugging (optional)
- **Load Balancer Integration** with target group attachment

## Architecture

```
┌─────────────────────────────────────────┐
│  Application Load Balancer              │
│  (Created by alb-ecs module)            │
└─────────────────────────────────────────┘
    |
    | Traffic forwarded to target group
    v
┌─────────────────────────────────────────┐
│  ECS Service (Fargate)                  │
│  - Desired Count: Configurable          │
│  - Network Mode: awsvpc                 │
│  - Security Group: ALB traffic only     │
└─────────────────────────────────────────┘
    |
    | Runs task definitions
    v
┌─────────────────────────────────────────┐
│  ECS Tasks                              │
│  - Container: Go Fiber API              │
│  - Health Check: /health endpoint       │
│  - CloudWatch Logs                      │
│  - Environment Variables                │
└─────────────────────────────────────────┘
```

## Prerequisites

- VPC with public subnets (created by vpc-network module)
- Application Load Balancer (created by alb-ecs module)
- ECR repository with container image (created by ecr-repository module)
- Target group for load balancer integration

## Usage

### Basic Configuration

```hcl
module "api" {
  source = "./modules/ecs-service"

  # Service configuration
  service_name = "crux-api"
  environment  = "prod"

  # Networking
  region            = "us-east-1"
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids

  # Load balancer integration
  target_group_arn      = module.alb.target_group_arn
  alb_security_group_id = module.alb.alb_security_group_id

  # Task resources
  task_cpu    = 256
  task_memory = 512

  # Container configuration
  container_image = "${aws_ecr_repository.api.repository_url}:latest"
  container_port  = 3000
  enable_ecs_exec = true

  # Environment variables
  environment_variables = [
    { name = "PORT", value = "3000" },
    { name = "ENVIRONMENT", value = "prod" }
  ]
}
```

### Complete Example with Database Integration

```hcl
module "api" {
  source = "./modules/ecs-service"

  service_name = "crux-backend-api"
  environment  = "prod"

  # Networking
  region            = "us-east-1"
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids

  # Load balancer (from alb-ecs module)
  target_group_arn      = module.alb.target_group_arn
  alb_security_group_id = module.alb.alb_security_group_id

  # Task configuration
  task_cpu    = 512
  task_memory = 1024
  desired_count = 2

  # Container
  container_image = "${module.ecr.repository_url}:latest"
  container_port  = 3000
  enable_ecs_exec = true

  # Application configuration
  environment_variables = [
    { name = "ENVIRONMENT", value = "prod" },
    { name = "PORT", value = "3000" },
    { name = "DB_HOST", value = module.db.address },
    { name = "DB_PORT", value = module.db.port },
    { name = "DB_NAME", value = module.db.database_name },
    { name = "DB_USER", value = module.db.username },
    { name = "DB_PASSWORD", value = "secure-password-from-secrets-manager" },
    { name = "DB_SSLMODE", value = "require" }
  ]

  # Monitoring
  enable_container_insights = true
  log_retention_days        = 14
}
```

## Input Variables

### Required Variables

| Name | Description | Type |
|------|-------------|------|
| `service_name` | Name of the service | `string` |
| `environment` | Environment (dev/stg/prod) | `string` |
| `region` | AWS region | `string` |
| `vpc_id` | VPC ID | `string` |
| `public_subnet_ids` | List of public subnet IDs | `list(string)` |
| `container_image` | Container image URI | `string` |
| `container_port` | Container port number | `number` |
| `target_group_arn` | ARN of ALB target group | `string` |
| `alb_security_group_id` | Security group ID of the ALB | `string` |
| `task_cpu` | CPU units (256, 512, 1024, etc.) | `number` |
| `task_memory` | Memory in MB (512, 1024, 2048, etc.) | `number` |
| `environment_variables` | List of environment variables | `list(object)` |
| `enable_ecs_exec` | Enable ECS Exec for debugging | `string` |

### Optional Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `desired_count` | Number of tasks to run | `number` | `1` |
| `enable_container_insights` | Enable Container Insights | `bool` | `true` |
| `log_retention_days` | CloudWatch log retention | `number` | `7` |

## Outputs

| Name | Description |
|------|-------------|
| `cluster_id` | ECS cluster ID |
| `cluster_name` | ECS cluster name |
| `cluster_arn` | ECS cluster ARN |
| `service_id` | ECS service ID |
| `service_name` | ECS service name |
| `task_definition_arn` | Task definition ARN |
| `task_security_group_id` | Security group ID for ECS tasks |
| `task_execution_role_arn` | Task execution role ARN |
| `task_role_arn` | Task role ARN |
| `log_group_name` | CloudWatch log group name |

## Integration with ALB Module

This module is designed to work with the `alb-ecs` module:

```hcl
# Create ALB first
module "alb" {
  source = "./modules/alb-ecs"

  service_name      = "crux-api"
  environment       = "prod"
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
  container_port    = 3000
  health_check_path = "/health"
}

# Create ECS service
module "api" {
  source = "./modules/ecs-service"

  service_name = "crux-api"
  environment  = "prod"

  # Pass ALB outputs to ECS service
  target_group_arn      = module.alb.target_group_arn
  alb_security_group_id = module.alb.alb_security_group_id

  # ... other configuration
}
```

## Security

### Security Group Configuration

The module creates a security group for ECS tasks with:

**Ingress Rules:**
- Allow traffic from ALB security group only on container port

**Egress Rules:**
- Allow all outbound traffic (for ECR pulls, RDS access, external APIs)

This ensures that ECS tasks can only receive traffic from the ALB, not directly from the internet.

### IAM Roles

Two IAM roles are created:

1. **Task Execution Role** - Used by ECS to pull images and write logs
   - Policy: `AmazonECSTaskExecutionRolePolicy`
   - Permissions: ECR pulls, CloudWatch logs

2. **Task Role** - Used by the running container
   - Can be customized with additional policies for AWS service access

## Container Health Checks

The task definition includes a health check:

```json
{
  "command": ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1"],
  "interval": 30,
  "timeout": 5,
  "retries": 3,
  "startPeriod": 60
}
```

Ensure your Go Fiber application implements the `/health` endpoint:

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
    })
})
```

## CloudWatch Logs

Logs are automatically sent to CloudWatch under the log group:
- Log Group: `/ecs/{service_name}-logs-{environment}`
- Stream Prefix: `{service_name}-{environment}`

View logs:
```bash
aws logs tail /ecs/crux-api-logs-prod --follow
```

## ECS Exec (Debugging)

When `enable_ecs_exec = true`, you can execute commands in running containers:

```bash
# List running tasks
aws ecs list-tasks --cluster crux-api-cluster-prod

# Execute shell in container
aws ecs execute-command \
  --cluster crux-api-cluster-prod \
  --task <task-id> \
  --container crux-api-prod \
  --interactive \
  --command "/bin/sh"
```

## Fargate CPU and Memory Combinations

Valid CPU and memory combinations:

| CPU (units) | Memory (MB) |
|-------------|-------------|
| 256 | 512, 1024, 2048 |
| 512 | 1024, 2048, 3072, 4096 |
| 1024 | 2048, 3072, 4096, 5120, 6144, 7168, 8192 |
| 2048 | 4096 - 16384 (1024 increments) |
| 4096 | 8192 - 30720 (1024 increments) |

## Troubleshooting

### Tasks Not Starting

Check:
- Container image exists in ECR
- Task execution role has ECR pull permissions
- Security group allows outbound traffic for ECR pulls
- CloudWatch log group exists

### Health Check Failures

Common causes:
- Application not listening on `0.0.0.0` (must not use `localhost`)
- `/health` endpoint returns non-200 status
- Container takes longer than 60 seconds to start

### ALB Connection Issues

Verify:
- Security group allows traffic from ALB SG
- Target group health check passes
- Container port matches ALB target group port
- ECS service is registered with correct target group

## Updates and Deployments

To deploy a new container version:

1. Build and push new image to ECR
2. Update the `container_image` variable
3. Run `terraform apply`
4. ECS will perform a rolling deployment

Or force new deployment without Terraform:
```bash
aws ecs update-service \
  --cluster crux-api-cluster-prod \
  --service crux-api-service-prod \
  --force-new-deployment
```

## Module Structure

```
ecs-service/
├── main.tf          # ECS cluster, service, task definition, security groups
├── variables.tf     # Input variable definitions
├── outputs.tf       # Output value definitions
└── README.md        # This file
```

## License

This module is part of the Crux Backend infrastructure.
