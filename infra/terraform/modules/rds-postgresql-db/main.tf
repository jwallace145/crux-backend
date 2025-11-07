locals {
  db_identifier     = replace(var.db_name, "_", "-")
  rds_db_identifier = "${local.db_identifier}-${var.environment}"
}

# =======================
# PostgreSQL RDS Database
# =======================

resource "aws_db_instance" "this" {
  # Database Configuration
  identifier     = local.rds_db_identifier
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
  storage_type          = "gp3"
  storage_encrypted     = true

  # Network Configuration
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.this.id]
  publicly_accessible    = var.publicly_accessible
  port                   = 5432

  # Backup Configuration
  backup_retention_period = var.backup_retention_days
  backup_window           = "03:00-04:00" # 3-4 AM UTC
  maintenance_window      = "Mon:04:00-Mon:05:00"

  # High Availability
  multi_az = var.multi_az

  # Performance Insights (DISABLED for cost savings)
  performance_insights_enabled = false

  # Enhanced Monitoring (DISABLED for cost savings)
  enabled_cloudwatch_logs_exports = [] # Would add CloudWatch costs

  # Deletion Configuration
  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "${local.rds_db_identifier}-final-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  deletion_protection       = var.deletion_protection

  # Auto Minor Version Updates
  auto_minor_version_upgrade = true

  # Parameter Group (use default to avoid extra costs)
  parameter_group_name = "default.postgres16"
}

# =====================
# Database Subnet Group
# =====================

resource "aws_db_subnet_group" "this" {
  name       = "${local.db_identifier}-subnet-group-${var.environment}"
  subnet_ids = var.subnet_ids

  tags = {
    Name = "${local.db_identifier}-subnet-group-${var.environment}"
  }
}

# =======================
# Database Security Group
# =======================

resource "aws_security_group" "this" {
  name        = "${local.db_identifier}-sg-${var.environment}"
  description = "Security group for ${local.rds_db_identifier} RDS database."
  vpc_id      = var.vpc_id

  egress {
    description = "Allow all outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${local.db_identifier}-sg-${var.environment}"
  }
}

# Security group rule for CIDR-based access (backward compatibility)
resource "aws_security_group_rule" "cidr_ingress" {
  count = length(var.allowed_cidr_blocks) > 0 ? 1 : 0

  security_group_id = aws_security_group.this.id
  description       = "PostgreSQL DB access from CIDR block(s) (${var.environment})"
  type              = "ingress"
  from_port         = local.POSTGRESQL_DB_PORT
  to_port           = local.POSTGRESQL_DB_PORT
  protocol          = "tcp"
  cidr_blocks       = var.allowed_cidr_blocks
}

# Security group rules for security group-based access (recommended)
resource "aws_security_group_rule" "sg_ingress" {
  count = length(var.allowed_security_group_ids)

  security_group_id        = aws_security_group.this.id
  description              = "PostgreSQL DB access from security group ${var.allowed_security_group_ids[count.index]} (${var.environment})"
  type                     = "ingress"
  from_port                = local.POSTGRESQL_DB_PORT
  to_port                  = local.POSTGRESQL_DB_PORT
  protocol                 = "tcp"
  source_security_group_id = var.allowed_security_group_ids[count.index]
}
