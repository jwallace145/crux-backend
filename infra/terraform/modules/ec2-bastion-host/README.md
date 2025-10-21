# EC2 Bastion Host Module

This Terraform module creates a lightweight EC2 bastion host that provides SSH access into your VPC for debugging and investigating network issues. The bastion host is configured with useful networking tools and can be used to access private resources like RDS databases, ECS tasks, and other internal services.

## Features

- **Small instance size** (t3.micro by default) to minimize costs
- **SSH access** from specific IP addresses only
- **Elastic IP** for consistent access (optional)
- **AWS Session Manager** support as SSH alternative
- **Pre-installed debugging tools** (PostgreSQL client, network utilities, AWS CLI)
- **IAM role** with optional ReadOnly access for AWS CLI commands
- **Encrypted EBS volume** for security
- **IMDSv2** enforced for metadata access
- **Automatic security updates** configured
- **Latest Amazon Linux 2023** AMI (or custom AMI)

## Architecture

```
Internet
    |
    | SSH (Port 22)
    v
┌─────────────────────────────────────────┐
│  Bastion Host (EC2)                     │
│  - Instance Type: t3.micro              │
│  - OS: Amazon Linux 2023                │
│  - Security Group: SSH from your IP     │
│  - Elastic IP: Static public IP         │
└─────────────────────────────────────────┘
    |
    | Access to VPC resources
    v
┌─────────────────────────────────────────┐
│  Private Resources in VPC               │
│  - RDS databases                        │
│  - ECS tasks                            │
│  - Private ALBs                         │
│  - ElastiCache clusters                 │
└─────────────────────────────────────────┘
```

## Prerequisites

- VPC with at least one public subnet
- SSH key pair (either existing or provide public key material)
- Your public IP address for SSH access

## Getting Your Public IP

```bash
# Get your current public IP
curl -4 ifconfig.me

# Or use dig
dig +short myip.opendns.com @resolver1.opendns.com
```

## Usage

### Basic Configuration (Recommended)

```hcl
module "bastion" {
  source = "./modules/ec2-bastion-host"

  service_name = "crux"
  environment  = "prod"

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  # Your public IP address (update this!)
  allowed_ssh_cidr_blocks = ["203.0.113.42/32"]

  # SSH key (choose one method)
  ssh_public_key = file("~/.ssh/id_rsa.pub")  # Option 1: Use existing key
  # key_name     = "my-existing-key"          # Option 2: Use AWS key pair

  # Allocate static IP for consistent access
  allocate_elastic_ip = true

  tags = {
    Project = "Crux"
    Owner   = "DevOps Team"
  }
}
```

### Production Configuration

```hcl
module "bastion" {
  source = "./modules/ec2-bastion-host"

  service_name = "crux"
  environment  = "prod"

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  # Allow SSH from office and VPN IPs
  allowed_ssh_cidr_blocks = [
    "203.0.113.42/32",  # Office IP
    "198.51.100.0/24"   # VPN subnet
  ]

  # Use existing key pair
  key_name = "crux-bastion-prod"

  # Production settings
  instance_type               = "t3.micro"
  allocate_elastic_ip         = true
  enable_termination_protection = true
  enable_readonly_access      = true  # Allow AWS CLI read operations

  tags = {
    Project     = "Crux"
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}
```

### Development Configuration (Cost-Optimized)

```hcl
module "bastion" {
  source = "./modules/ec2-bastion-host"

  service_name = "crux"
  environment  = "dev"

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  # Your local IP
  allowed_ssh_cidr_blocks = ["203.0.113.42/32"]

  # Use smallest instance
  instance_type = "t3.nano"

  # Don't need static IP in dev
  allocate_elastic_ip = false

  # SSH key
  ssh_public_key = file("~/.ssh/id_rsa.pub")

  tags = {
    Project     = "Crux"
    Environment = "Development"
  }
}
```

## Connecting to the Bastion Host

### Method 1: SSH (Recommended for debugging)

```bash
# Get the public IP from Terraform outputs
terraform output -json | jq -r '.bastion_public_ip.value'

# Connect via SSH
ssh -i ~/.ssh/your-key.pem ec2-user@<bastion-public-ip>

# Or use the auto-generated command
terraform output bastion_ssh_command
```

### Method 2: AWS Session Manager (No SSH key required)

```bash
# Start session (no SSH key needed!)
aws ssm start-session --target <instance-id>

# Or use the output command
terraform output bastion_ssm_command
```

## Common Debugging Tasks

### Connect to RDS Database

```bash
# SSH to bastion
ssh -i ~/.ssh/key.pem ec2-user@<bastion-ip>

# Connect to PostgreSQL
psql -h crux-db-prod.abc123.us-east-1.rds.amazonaws.com -U cruxadmin -d cruxdb

# Test connectivity
nc -zv crux-db-prod.abc123.us-east-1.rds.amazonaws.com 5432
```

### Test Network Connectivity

```bash
# Test port connectivity
nc -zv <hostname> <port>

# DNS lookup
dig api.crux.com
nslookup rds-endpoint.amazonaws.com

# Trace route
traceroute api.crux.com

# Check if service is listening
telnet <hostname> <port>
```

### Debug ECS Tasks

```bash
# Install ECS CLI
curl -Lo /usr/local/bin/ecs-cli https://amazon-ecs-cli.s3.amazonaws.com/ecs-cli-linux-amd64-latest
chmod +x /usr/local/bin/ecs-cli

# List tasks
aws ecs list-tasks --cluster crux-api-cluster-prod

# Describe task
aws ecs describe-tasks --cluster crux-api-cluster-prod --tasks <task-id>

# Check task ENI and IP
aws ecs describe-tasks --cluster crux-api-cluster-prod --tasks <task-id> | jq '.tasks[0].attachments[0].details'
```

### Test HTTP Endpoints

```bash
# Test internal ALB
curl http://internal-alb.internal:8080/health

# Test with headers
curl -H "Host: api.crux.com" http://10.0.1.100/health

# Verbose output
curl -v http://api.internal/health
```

## Pre-installed Tools

The bastion host comes with these tools pre-installed:

- **Network Tools**: `dig`, `nslookup`, `netcat`, `telnet`, `traceroute`, `net-tools`
- **Database Clients**: `psql` (PostgreSQL 15)
- **AWS Tools**: AWS CLI, SSM Agent
- **Utilities**: `jq`, `vim`, `git`, `htop`, `curl`, `wget`

## Input Variables

### Required Variables

| Name | Description | Type |
|------|-------------|------|
| `service_name` | Service name | `string` |
| `environment` | Environment (dev/stg/prod) | `string` |
| `vpc_id` | VPC ID | `string` |
| `subnet_id` | Public subnet ID | `string` |
| `allowed_ssh_cidr_blocks` | CIDR blocks allowed to SSH | `list(string)` |

### Optional Variables

| Name | Description | Type | Default |
|------|-------------|------|---------|
| `instance_type` | EC2 instance type | `string` | `"t3.micro"` |
| `ami_id` | AMI ID (defaults to Amazon Linux 2023) | `string` | `""` |
| `key_name` | Existing EC2 key pair name | `string` | `""` |
| `ssh_public_key` | SSH public key content | `string` | `""` |
| `root_volume_size` | Root volume size in GB | `number` | `8` |
| `allocate_elastic_ip` | Allocate Elastic IP | `bool` | `true` |
| `enable_readonly_access` | Attach ReadOnlyAccess IAM policy | `bool` | `false` |
| `enable_termination_protection` | Enable termination protection | `bool` | `false` |
| `enable_detailed_monitoring` | Enable detailed monitoring | `bool` | `false` |
| `user_data` | Custom user data script | `string` | `""` |
| `tags` | Additional tags | `map(string)` | `{}` |

## Outputs

| Name | Description |
|------|-------------|
| `instance_id` | Instance ID |
| `instance_arn` | Instance ARN |
| `public_ip` | Public IP address |
| `private_ip` | Private IP address |
| `elastic_ip_id` | Elastic IP allocation ID |
| `public_dns` | Public DNS name |
| `security_group_id` | Security group ID |
| `iam_role_arn` | IAM role ARN |
| `ssh_command` | SSH connection command |
| `ssm_command` | SSM Session Manager command |

## Security Best Practices

### SSH Key Management

**Option 1: Use existing AWS key pair**
```hcl
key_name = "my-existing-key"
```

**Option 2: Provide public key (Terraform creates key pair)**
```hcl
ssh_public_key = file("~/.ssh/id_rsa.pub")
```

Generate a new SSH key if needed:
```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/crux-bastion -C "bastion-access"
```

### IP Whitelisting

Always restrict SSH access to specific IPs:

```hcl
# ✅ Good: Specific IPs only
allowed_ssh_cidr_blocks = ["203.0.113.42/32"]

# ❌ Bad: Open to the internet
allowed_ssh_cidr_blocks = ["0.0.0.0/0"]
```

### Instance Management

- **Stop when not in use**: Bastion hosts should be stopped when not actively debugging
- **Enable termination protection** for production environments
- **Use Session Manager** when possible (no SSH key exposure)
- **Review CloudTrail logs** for bastion access audit

### Cost Optimization

```bash
# Stop instance when not in use
aws ec2 stop-instances --instance-ids <instance-id>

# Start instance when needed
aws ec2 start-instances --instance-ids <instance-id>

# Get new public IP after start (if not using Elastic IP)
aws ec2 describe-instances --instance-ids <instance-id> --query 'Reservations[0].Instances[0].PublicIpAddress'
```

## SSH Tunneling

Use the bastion for SSH tunneling to access private resources:

### Port Forwarding to RDS

```bash
# Forward local port 5432 to RDS through bastion
ssh -i ~/.ssh/key.pem -L 5432:rds-endpoint.amazonaws.com:5432 ec2-user@<bastion-ip>

# Now connect to localhost
psql -h localhost -U username -d db
```

### SOCKS Proxy

```bash
# Create SOCKS proxy through bastion
ssh -i ~/.ssh/key.pem -D 8080 -N ec2-user@<bastion-ip>

# Configure browser to use localhost:8080 as SOCKS5 proxy
```

## Troubleshooting

### Cannot SSH to Bastion

1. **Check security group**: Verify your IP is in `allowed_ssh_cidr_blocks`
2. **Verify your current IP**: `curl ifconfig.me`
3. **Check instance state**: `aws ec2 describe-instances --instance-ids <id>`
4. **Try Session Manager**: `aws ssm start-session --target <id>`

### SSH Connection Times Out

- Verify bastion is in a **public subnet**
- Check subnet has an **Internet Gateway** attached
- Verify **route table** has route to 0.0.0.0/0 via IGW
- Confirm instance has a **public IP** assigned

### Cannot Access Private Resources

1. **Check security groups**: Ensure private resource SG allows traffic from bastion SG
2. **Verify routing**: Check route tables for proper VPC routing
3. **Test connectivity**: Use `nc -zv <host> <port>` from bastion

## Example: Adding Bastion to Existing Infrastructure

```hcl
# In your main.tf or bastion.tf

module "bastion" {
  source = "./modules/ec2-bastion-host"

  service_name = var.service_name
  environment  = var.environment

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  allowed_ssh_cidr_blocks = var.bastion_allowed_ips

  ssh_public_key = var.ssh_public_key

  tags = merge(
    var.tags,
    {
      Purpose = "Network debugging and RDS access"
    }
  )
}

# Allow bastion to access RDS
resource "aws_security_group_rule" "rds_from_bastion" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = module.bastion.security_group_id
  security_group_id        = module.db.security_group_id
  description              = "Allow PostgreSQL access from bastion host"
}

# Outputs
output "bastion_ip" {
  description = "Bastion host public IP"
  value       = module.bastion.public_ip
}

output "bastion_ssh_command" {
  description = "Command to SSH to bastion"
  value       = module.bastion.ssh_command
}
```

## Module Structure

```
ec2-bastion-host/
├── main.tf          # EC2 instance, security group, IAM role
├── variables.tf     # Input variable definitions
├── outputs.tf       # Output value definitions
├── user_data.sh     # Instance initialization script
└── README.md        # This file
```

## Cost Estimate

| Component | Type | Monthly Cost (us-east-1) |
|-----------|------|--------------------------|
| EC2 Instance | t3.micro (stopped 20h/day) | ~$1.50 |
| EC2 Instance | t3.micro (24/7) | ~$7.50 |
| Elastic IP (attached) | N/A | $0 |
| Elastic IP (unattached) | Per hour | ~$3.60 |
| EBS Volume | 8 GB gp3 | ~$0.64 |

**Recommendation**: Stop the bastion when not in use to save ~80% on EC2 costs.

## License

This module is part of the Crux Backend infrastructure.
