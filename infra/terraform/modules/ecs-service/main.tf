resource "aws_ecs_cluster" "cluster" {
  name = "${var.service_name}-cluster-${var.environment}"

  setting {
    name  = "containerInsights"
    value = var.enable_container_insights ? "enabled" : "disabled"
  }

  tags = {
    Name = "${var.service_name}-cluster-${var.environment}"
  }
}

resource "aws_ecs_cluster_capacity_providers" "cluster_capacity_providers" {
  cluster_name = aws_ecs_cluster.cluster.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

resource "aws_cloudwatch_log_group" "task_log_group" {
  name              = "/ecs/${var.service_name}-logs-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = {
    Name = "${var.service_name}-logs-${var.environment}"
  }
}

resource "aws_iam_role" "task_execution_role" {
  name = "${var.service_name}-ecs-task-execution-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "${var.service_name}-ecs-task-execution-role-${var.environment}"
  }
}

resource "aws_iam_role_policy_attachment" "task_execution_policy_attachment" {
  role       = aws_iam_role.task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "task_role" {
  name = "${var.service_name}-ecs-task-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name = "${var.service_name}-ecs-task-role-${var.environment}"
  }
}

resource "aws_security_group" "task_security_group" {
  name        = "${var.service_name}-ecs-task-sg-${var.environment}"
  description = "Security group for ECS tasks"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow HTTP from anywhere"
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.service_name}-ecs-task-${var.environment}"
  }
}

resource "aws_ecs_task_definition" "task_definition" {
  family                   = "${var.service_name}-${var.environment}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.task_execution_role.arn
  task_role_arn            = aws_iam_role.task_role.arn

  container_definitions = jsonencode([
    {
      name      = "${var.service_name}-${var.environment}"
      image     = var.container_image
      essential = true

      portMappings = [
        {
          containerPort = var.container_port
          hostPort      = var.container_port
          protocol      = "tcp"
        }
      ]

      environment = var.environment_variables

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.task_log_group.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "${var.service_name}-${var.environment}"
        }
      }

      healthCheck = {
        command     = ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:${var.container_port}/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }
    }
  ])

  tags = {
    Name = "${var.service_name}-task-definition-${var.environment}"
  }
}

resource "aws_ecs_service" "service" {
  name            = "${var.service_name}-service-${var.environment}"
  cluster         = aws_ecs_cluster.cluster.id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.public_subnet_ids
    security_groups  = [aws_security_group.task_security_group.id]
    assign_public_ip = true
  }

  enable_execute_command = var.enable_ecs_exec

  tags = {
    Name = "${var.service_name}-service-${var.environment}"
  }

  depends_on = [
    aws_iam_role_policy_attachment.task_execution_policy_attachment
  ]
}

################################################################################
# Allow ECS to access RDS
################################################################################

# resource "aws_security_group_rule" "rds_from_ecs" {
#   type                     = "ingress"
#   from_port                = var.db_port
#   to_port                  = var.db_port
#   protocol                 = "tcp"
#   security_group_id        = var.db_security_group_id
#   source_security_group_id = aws_security_group.ecs_task.id
#   description              = "Allow ECS tasks to access RDS"
# }
