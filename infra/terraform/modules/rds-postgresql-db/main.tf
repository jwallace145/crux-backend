resource "aws_db_instance" "this" {
  # Database Configuration
  identifier     = "${var.db_name}-${var.environment}"
  engine         = "postgres"
  engine_version = var.postgres_version
  instance_class = var.instance_class

  # Database Credentials
  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  # Storage Configuration (cost-optimized)
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = "gp3" # Cheaper than gp2, better performance
  storage_encrypted     = true  # Free encryption at rest

  # Network Configuration
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.this.id]
  publicly_accessible    = var.publicly_accessible
  port                   = 5432

  # Backup Configuration
  backup_retention_period = var.backup_retention_days
  backup_window           = "03:00-04:00" # 3-4 AM UTC
  maintenance_window      = "Mon:04:00-Mon:05:00"

  # High Availability (DISABLED for cost savings)
  multi_az = false # Enabling would double cost

  # Performance Insights (DISABLED for cost savings)
  performance_insights_enabled = false # Would add ~$7/month

  # Enhanced Monitoring (DISABLED for cost savings)
  enabled_cloudwatch_logs_exports = [] # Would add CloudWatch costs

  # Deletion Configuration
  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "${var.db_name}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  deletion_protection       = var.deletion_protection

  # Auto Minor Version Updates
  auto_minor_version_upgrade = true

  # Parameter Group (use default to avoid extra costs)
  parameter_group_name = "default.postgres16"
}

resource "aws_db_subnet_group" "this" {
  name       = "${var.db_name}-subnet-group-${var.environment}"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "${var.db_name}-subnet-group-${var.environment}"
  }
}

resource "aws_security_group" "this" {
  name        = "${var.db_name}-rds-sg-${var.environment}"
  description = "Security group for ${var.db_name} RDS instance."
  vpc_id      = var.vpc_id

  ingress {
    description = "CruxBackend PostgreSQL DB access"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
  }

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.db_name}-rds-sg"
  }
}
