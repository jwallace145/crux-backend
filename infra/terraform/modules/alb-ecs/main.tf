# ============================================================================
# Application Load Balancer Module for ECS
# ============================================================================
# This module creates an Application Load Balancer (ALB) that routes HTTP/HTTPS
# traffic to ECS tasks running the Crux Backend API. The ALB provides a stable
# DNS name that can be used as a CNAME target for custom domains.
# ============================================================================

# ============================================================================
# Security Group for ALB
# ============================================================================

resource "aws_security_group" "alb" {
  name        = "${var.service_name}-${var.environment}-alb-sg"
  description = "Security group for Application Load Balancer"
  vpc_id      = var.vpc_id

  # HTTP ingress from anywhere
  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
  }

  # HTTPS ingress from anywhere (if HTTPS is enabled)
  dynamic "ingress" {
    for_each = var.enable_https ? [1] : []
    content {
      description = "HTTPS from anywhere"
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = var.allowed_cidr_blocks
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-${var.environment}-alb-sg"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# Security Group for ECS Tasks
# ============================================================================

resource "aws_security_group" "ecs_tasks" {
  name        = "${var.service_name}-${var.environment}-ecs-tasks-sg"
  description = "Security group for ECS tasks to receive traffic from ALB"
  vpc_id      = var.vpc_id

  # Egress to internet (for external API calls, ECR pulls, etc.)
  egress {
    description = "To internet"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-${var.environment}-ecs-tasks-sg"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# Security Group Rules (to avoid circular dependency)
# ============================================================================

# Allow ALB to communicate with ECS tasks
resource "aws_security_group_rule" "alb_to_ecs" {
  type                     = "egress"
  from_port                = var.container_port
  to_port                  = var.container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.ecs_tasks.id
  security_group_id        = aws_security_group.alb.id
  description              = "Allow ALB to communicate with ECS tasks"
}

# Allow ECS tasks to receive traffic from ALB
resource "aws_security_group_rule" "ecs_from_alb" {
  type                     = "ingress"
  from_port                = var.container_port
  to_port                  = var.container_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  security_group_id        = aws_security_group.ecs_tasks.id
  description              = "Allow ECS tasks to receive traffic from ALB"
}

# ============================================================================
# Application Load Balancer
# ============================================================================

resource "aws_lb" "main" {
  name               = "${var.service_name}-alb-${var.environment}"
  internal           = var.internal
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection       = var.enable_deletion_protection
  enable_http2                     = var.enable_http2
  enable_cross_zone_load_balancing = true

  idle_timeout = var.idle_timeout

  access_logs {
    enabled = var.enable_access_logs
    bucket  = var.access_logs_bucket
    prefix  = var.access_logs_prefix
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-alb-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# Target Group for ECS Tasks
# ============================================================================

resource "aws_lb_target_group" "ecs" {
  name        = "${var.service_name}-${var.environment}-tg"
  port        = var.container_port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip" # Required for Fargate with awsvpc network mode

  # Health check configuration
  health_check {
    enabled             = true
    healthy_threshold   = var.health_check_healthy_threshold
    unhealthy_threshold = var.health_check_unhealthy_threshold
    timeout             = var.health_check_timeout
    interval            = var.health_check_interval
    path                = var.health_check_path
    matcher             = var.health_check_matcher
    protocol            = "HTTP"
  }

  # Deregistration delay
  deregistration_delay = var.deregistration_delay

  # Stickiness configuration (optional)
  dynamic "stickiness" {
    for_each = var.enable_stickiness ? [1] : []
    content {
      type            = "lb_cookie"
      cookie_duration = var.stickiness_duration
      enabled         = true
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-${var.environment}-tg"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )

  # Ensure target group is created before ECS service
  lifecycle {
    create_before_destroy = true
  }
}

# ============================================================================
# HTTP Listener (Port 80)
# ============================================================================

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  # Default action - forward to ECS target group or redirect to HTTPS
  dynamic "default_action" {
    for_each = var.enable_https && var.redirect_http_to_https ? [1] : []
    content {
      type = "redirect"

      redirect {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
  }

  dynamic "default_action" {
    for_each = var.enable_https && var.redirect_http_to_https ? [] : [1]
    content {
      type             = "forward"
      target_group_arn = aws_lb_target_group.ecs.arn
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-${var.environment}-http-listener"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# HTTPS Listener (Port 443) - Optional
# ============================================================================

resource "aws_lb_listener" "https" {
  count = var.enable_https ? 1 : 0

  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = var.ssl_policy
  certificate_arn   = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ecs.arn
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-${var.environment}-https-listener"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}


# ============================================================================
# ALB DNS Record
# ============================================================================

data "aws_route53_zone" "hosted_zone" {
  name         = var.domain
  private_zone = false
}

resource "aws_route53_record" "api" {
  zone_id = data.aws_route53_zone.hosted_zone.zone_id
  name    = "${var.environment}-api.${var.domain}"
  type    = "A"

  alias {
    name                   = aws_lb.main.dns_name
    zone_id                = aws_lb.main.zone_id
    evaluate_target_health = true
  }
}
