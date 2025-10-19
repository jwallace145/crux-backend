# Application Load Balancer for ECS Module

This Terraform module creates an Application Load Balancer (ALB) that routes HTTP/HTTPS traffic to ECS tasks running the Crux Backend API. The ALB provides a stable DNS name that can be used as a CNAME target for custom domains.

## Features

- **Internet-facing or internal ALB** with configurable accessibility
- **HTTP (port 80) and HTTPS (port 443)** listeners with optional redirect
- **SSL/TLS termination** with customizable security policies
- **Security groups** for ALB and ECS tasks with least-privilege access
- **Target group** configured for ECS Fargate tasks (awsvpc network mode)
- **Health checks** with configurable parameters
- **Sticky sessions** support (optional)
- **Access logs** to S3 (optional)
- **Cross-zone load balancing** enabled by default
- **HTTP/2 support** enabled by default

## Architecture

```
Internet/VPC
    |
    v
┌─────────────────────────────────────────┐
│  Application Load Balancer              │
│  - HTTP Listener (Port 80)              │
│  - HTTPS Listener (Port 443) [Optional] │
│  - Security Group (Allow 80/443)        │
└─────────────────────────────────────────┘
    |
    | Forward to Target Group
    v
┌─────────────────────────────────────────┐
│  Target Group (IP Target Type)          │
│  - Health Check: GET /health            │
│  - Deregistration Delay: 30s            │
└─────────────────────────────────────────┘
    |
    | Route to ECS Tasks
    v
┌─────────────────────────────────────────┐
│  ECS Tasks (Fargate)                    │
│  - Container Port: 3000                 │
│  - Security Group (Allow from ALB only) │
└─────────────────────────────────────────┘
```

## Prerequisites

- Existing VPC with at least 2 public subnets in different AZs
- SSL certificate in ACM (if HTTPS is enabled)
- S3 bucket for access logs (if access logging is enabled)

## Usage

### Basic HTTP-only Configuration

```hcl
module "alb" {
  source = "./modules/alb-ecs"

  service_name       = "crux-backend"
  environment        = "prod"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  container_port     = 3000

  # Health check configuration
  health_check_path     = "/health"
  health_check_interval = 30
  health_check_timeout  = 5

  tags = {
    Project = "Crux"
    Owner   = "DevOps Team"
  }
}
```

### HTTPS with SSL Certificate

```hcl
module "alb" {
  source = "./modules/alb-ecs"

  service_name       = "crux-backend"
  environment        = "prod"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  container_port     = 3000

  # Enable HTTPS
  enable_https           = true
  certificate_arn        = "arn:aws:acm:us-east-1:123456789012:certificate/abc-123"
  ssl_policy             = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  redirect_http_to_https = true

  # Health check
  health_check_path = "/health"

  tags = {
    Project = "Crux"
  }
}
```

### Internal ALB with Access Logs

```hcl
module "alb" {
  source = "./modules/alb-ecs"

  service_name       = "crux-backend"
  environment        = "staging"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.private_subnet_ids
  container_port     = 3000

  # Internal ALB (not internet-facing)
  internal = true

  # Restrict access to specific CIDR blocks
  allowed_cidr_blocks = ["10.0.0.0/8"]

  # Enable access logs
  enable_access_logs  = true
  access_logs_bucket  = "my-alb-logs-bucket"
  access_logs_prefix  = "crux-backend-staging"

  # Health check
  health_check_path = "/health"

  tags = {
    Project     = "Crux"
    Environment = "Staging"
  }
}
```

### Production Configuration with Sticky Sessions

```hcl
module "alb" {
  source = "./modules/alb-ecs"

  service_name       = "crux-backend"
  environment        = "prod"
  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  container_port     = 3000

  # HTTPS configuration
  enable_https           = true
  certificate_arn        = var.ssl_certificate_arn
  redirect_http_to_https = true

  # Enable deletion protection for production
  enable_deletion_protection = true

  # Sticky sessions for session affinity
  enable_stickiness   = true
  stickiness_duration = 3600 # 1 hour

  # Fine-tuned health checks
  health_check_path               = "/health"
  health_check_interval           = 15
  health_check_timeout            = 5
  health_check_healthy_threshold  = 2
  health_check_unhealthy_threshold = 3
  health_check_matcher            = "200-299"

  # Longer deregistration delay for graceful shutdown
  deregistration_delay = 60

  # Access logs
  enable_access_logs  = true
  access_logs_bucket  = var.access_logs_bucket
  access_logs_prefix  = "crux-backend-prod"

  tags = {
    Project     = "Crux"
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}
```

## Integrating with ECS Service

Use the outputs from this module when configuring your ECS service:

```hcl
resource "aws_ecs_service" "main" {
  name            = "crux-backend"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.main.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = var.private_subnet_ids
    security_groups = [module.alb.ecs_tasks_security_group_id]
  }

  load_balancer {
    target_group_arn = module.alb.target_group_arn
    container_name   = "crux-backend"
    container_port   = 3000
  }

  depends_on = [module.alb]
}
```

## Setting Up Custom Domain with Route53

Create a CNAME record pointing to the ALB DNS name:

```hcl
resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.crux.com"
  type    = "A"

  alias {
    name                   = module.alb.alb_dns_name
    zone_id                = module.alb.alb_zone_id
    evaluate_target_health = true
  }
}
```

Or using CNAME (less recommended):

```hcl
resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.crux.com"
  type    = "CNAME"
  ttl     = 300
  records = [module.alb.alb_dns_name]
}
```

## Input Variables

### Required Variables

| Name | Description | Type | Example |
|------|-------------|------|---------|
| `service_name` | Name of the service | `string` | `"crux-backend"` |
| `environment` | Environment name (dev/staging/prod) | `string` | `"prod"` |
| `vpc_id` | VPC ID | `string` | `"vpc-abc123"` |
| `public_subnet_ids` | List of public subnet IDs (min 2) | `list(string)` | `["subnet-1", "subnet-2"]` |
| `container_port` | Container port number | `number` | `3000` |

### Optional Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `health_check_path` | Health check path | `string` | `"/health"` |
| `health_check_interval` | Health check interval (seconds) | `number` | `30` |
| `health_check_timeout` | Health check timeout (seconds) | `number` | `5` |
| `health_check_healthy_threshold` | Healthy threshold count | `number` | `2` |
| `health_check_unhealthy_threshold` | Unhealthy threshold count | `number` | `3` |
| `health_check_matcher` | HTTP status codes for healthy | `string` | `"200"` |
| `enable_https` | Enable HTTPS listener | `bool` | `false` |
| `certificate_arn` | ACM certificate ARN | `string` | `""` |
| `ssl_policy` | SSL security policy | `string` | `"ELBSecurityPolicy-TLS13-1-2-2021-06"` |
| `redirect_http_to_https` | Redirect HTTP to HTTPS | `bool` | `true` |
| `internal` | Internal ALB (true) or internet-facing (false) | `bool` | `false` |
| `enable_deletion_protection` | Enable deletion protection | `bool` | `false` |
| `enable_http2` | Enable HTTP/2 | `bool` | `true` |
| `idle_timeout` | Connection idle timeout (seconds) | `number` | `60` |
| `allowed_cidr_blocks` | CIDR blocks allowed to access ALB | `list(string)` | `["0.0.0.0/0"]` |
| `enable_access_logs` | Enable access logging | `bool` | `false` |
| `access_logs_bucket` | S3 bucket for access logs | `string` | `""` |
| `access_logs_prefix` | Access logs prefix in S3 | `string` | `"alb-logs"` |
| `deregistration_delay` | Target deregistration delay (seconds) | `number` | `30` |
| `enable_stickiness` | Enable sticky sessions | `bool` | `false` |
| `stickiness_duration` | Sticky session duration (seconds) | `number` | `86400` |
| `tags` | Additional tags | `map(string)` | `{}` |

## Outputs

| Name | Description |
|------|-------------|
| `alb_dns_name` | DNS name of the ALB - use as CNAME target |
| `alb_arn` | ARN of the ALB |
| `alb_zone_id` | Canonical hosted zone ID (for Route53) |
| `alb_id` | ID of the ALB |
| `target_group_arn` | ARN of the target group for ECS |
| `target_group_name` | Name of the target group |
| `alb_security_group_id` | Security group ID of the ALB |
| `ecs_tasks_security_group_id` | Security group ID for ECS tasks |
| `http_listener_arn` | ARN of the HTTP listener |
| `https_listener_arn` | ARN of the HTTPS listener |
| `alb_url` | Full HTTP URL of the ALB |
| `alb_https_url` | Full HTTPS URL of the ALB |

## Security Considerations

### Security Groups

The module creates two security groups:

1. **ALB Security Group**
   - Ingress: Allows HTTP (80) and HTTPS (443) from `allowed_cidr_blocks`
   - Egress: Only allows traffic to ECS tasks security group on `container_port`

2. **ECS Tasks Security Group**
   - Ingress: Only allows traffic from ALB security group on `container_port`
   - Egress: Allows all outbound traffic (for ECR pulls, external API calls, etc.)

### Best Practices

- **Production**: Set `enable_deletion_protection = true`
- **HTTPS**: Always enable HTTPS for production with `enable_https = true`
- **SSL Policy**: Use the latest TLS policy: `ELBSecurityPolicy-TLS13-1-2-2021-06`
- **Access Logs**: Enable for production to track access patterns and debug issues
- **CIDR Blocks**: Restrict `allowed_cidr_blocks` for internal ALBs
- **Health Checks**: Ensure your application responds quickly to health check requests

## Health Check Configuration

The default health check configuration is:
- Path: `/health`
- Interval: 30 seconds
- Timeout: 5 seconds
- Healthy threshold: 2 consecutive successes
- Unhealthy threshold: 3 consecutive failures
- Expected response: HTTP 200

Ensure your Go Fiber application has a health endpoint:

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
    })
})
```

## Sticky Sessions

When `enable_stickiness = true`, the ALB uses load balancer cookies to route requests from the same client to the same target. This is useful for:
- Session affinity requirements
- WebSocket connections
- Applications that cache data in memory

Note: Sticky sessions may impact load distribution.

## Access Logs

To enable access logs:

1. Create an S3 bucket with appropriate permissions
2. Set `enable_access_logs = true`
3. Provide `access_logs_bucket` name
4. Optionally set `access_logs_prefix`

Access logs include:
- Request time
- Client IP
- Request method and URL
- Response status
- User agent
- SSL cipher and protocol

## Troubleshooting

### Targets Not Registering

If ECS tasks aren't registering with the target group:
- Verify the target type is `ip` (required for Fargate with awsvpc mode)
- Check security group rules allow ALB → ECS communication
- Ensure health check path returns 200 status
- Review CloudWatch logs for health check failures

### 502 Bad Gateway Errors

Common causes:
- Container isn't listening on `0.0.0.0` (must not use `localhost`)
- Security group blocks ALB → ECS traffic
- Health check is failing
- Container hasn't started yet

### HTTPS Not Working

Check:
- Certificate ARN is correct and in same region
- Certificate is validated in ACM
- `enable_https = true` is set
- DNS points to ALB

## Module Structure

```
alb-ecs/
├── main.tf          # Main ALB, target group, and listener resources
├── variables.tf     # Input variable definitions with validation
├── outputs.tf       # Output value definitions
└── README.md        # This file
```

## License

This module is part of the Crux Backend infrastructure.

## Support

For issues or questions, please contact the DevOps team or open an issue in the repository.
