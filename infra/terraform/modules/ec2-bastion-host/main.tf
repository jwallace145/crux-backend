# ============================================================================
# EC2 Bastion Host Module
# ============================================================================
# This module creates a small EC2 instance that serves as a bastion host
# for SSH access into the VPC to investigate network issues and access
# private resources.
# ============================================================================

# ============================================================================
# Data Sources
# ============================================================================

# Get latest Amazon Linux 2023 AMI
data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

# ============================================================================
# Security Group for Bastion Host
# ============================================================================

resource "aws_security_group" "bastion" {
  name        = "${var.service_name}-bastion-sg-${var.environment}"
  description = "Security group for bastion host - allows SSH from specified IPs"
  vpc_id      = var.vpc_id

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-sg-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# Allow SSH from specified CIDR blocks
resource "aws_security_group_rule" "bastion_ssh_ingress" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = var.allowed_ssh_cidr_blocks
  security_group_id = aws_security_group.bastion.id
  description       = "Allow SSH from specified IP addresses"
}

# Allow all outbound traffic (for package updates, AWS API calls, etc.)
resource "aws_security_group_rule" "bastion_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.bastion.id
  description       = "Allow all outbound traffic"
}

# ============================================================================
# IAM Role for Bastion Host
# ============================================================================

resource "aws_iam_role" "bastion" {
  name = "${var.service_name}-bastion-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-role-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# Attach SSM policy for Session Manager access (alternative to SSH)
resource "aws_iam_role_policy_attachment" "bastion_ssm" {
  role       = aws_iam_role.bastion.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

# Optional: Attach read-only policy for AWS CLI debugging
resource "aws_iam_role_policy_attachment" "bastion_readonly" {
  count      = var.enable_readonly_access ? 1 : 0
  role       = aws_iam_role.bastion.name
  policy_arn = "arn:aws:iam::aws:policy/ReadOnlyAccess"
}

# Create instance profile
resource "aws_iam_instance_profile" "bastion" {
  name = "${var.service_name}-bastion-profile-${var.environment}"
  role = aws_iam_role.bastion.name

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-profile-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# EC2 Key Pair
# ============================================================================

resource "aws_key_pair" "bastion" {
  key_name   = "${var.service_name}-bastion-key-${var.environment}"
  public_key = file(var.ssh_public_key_path)

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-key-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )
}

# ============================================================================
# EC2 Bastion Instance
# ============================================================================

resource "aws_instance" "bastion" {
  ami                    = data.aws_ami.amazon_linux_2023.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.bastion.key_name
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [aws_security_group.bastion.id]
  iam_instance_profile   = aws_iam_instance_profile.bastion.name

  # Enable detailed monitoring
  monitoring = var.enable_detailed_monitoring

  # Instance metadata options (IMDSv2)
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
    instance_metadata_tags      = "enabled"
  }

  # Root volume configuration
  root_block_device {
    volume_type           = "gp3"
    volume_size           = var.root_volume_size
    delete_on_termination = true
    encrypted             = true

    tags = merge(
      var.tags,
      {
        Name        = "${var.service_name}-bastion-root-${var.environment}"
        Environment = var.environment
        ManagedBy   = "terraform"
      }
    )
  }

  # User data script for initial setup
  user_data = var.user_data != "" ? var.user_data : templatefile("${path.module}/user_data.sh", {
    hostname = "${var.service_name}-bastion-${var.environment}"
  })

  # Enable termination protection in production
  disable_api_termination = var.enable_termination_protection

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-${var.environment}"
      Environment = var.environment
      Role        = "bastion"
      ManagedBy   = "terraform"
    }
  )

  lifecycle {
    ignore_changes = [
      ami,
      user_data
    ]
  }
}

# ============================================================================
# Elastic IP (if enabled)
# ============================================================================

resource "aws_eip" "bastion" {
  count    = var.allocate_elastic_ip ? 1 : 0
  domain   = "vpc"
  instance = aws_instance.bastion.id

  tags = merge(
    var.tags,
    {
      Name        = "${var.service_name}-bastion-eip-${var.environment}"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  )

  depends_on = [aws_instance.bastion]
}
