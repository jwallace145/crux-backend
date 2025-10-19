#!/bin/bash
# ============================================================================
# Bastion Host User Data Script
# ============================================================================
# This script runs on first boot to configure the bastion host
# ============================================================================

set -e

# Update all packages
echo "Updating system packages..."
yum update -y

# Install useful tools for debugging
echo "Installing debugging tools..."
yum install -y \
  bind-utils \
  curl \
  git \
  htop \
  jq \
  nc \
  net-tools \
  telnet \
  traceroute \
  vim \
  wget

# Install PostgreSQL client for database debugging
echo "Installing PostgreSQL client..."
yum install -y postgresql15

# Install Session Manager plugin (for AWS SSM)
echo "Installing AWS Session Manager plugin..."
yum install -y amazon-ssm-agent
systemctl enable amazon-ssm-agent
systemctl start amazon-ssm-agent

# Set hostname
echo "Setting hostname..."
hostnamectl set-hostname ${hostname}

# Configure timezone
echo "Configuring timezone..."
timedatectl set-timezone UTC

# Enable automatic security updates
echo "Enabling automatic security updates..."
yum install -y yum-cron
sed -i 's/apply_updates = no/apply_updates = yes/' /etc/yum/yum-cron.conf
systemctl enable yum-cron
systemctl start yum-cron

# Create a message of the day
cat > /etc/motd << 'EOF'
================================================================================
  BASTION HOST - ${hostname}
================================================================================
  This is a bastion host for accessing private resources in the VPC.

  Available tools:
  - PostgreSQL client (psql)
  - Network utilities (dig, nslookup, netcat, telnet, traceroute)
  - AWS CLI
  - jq (JSON processor)

  Usage:
  - Connect to RDS: psql -h <rds-endpoint> -U <username> -d <database>
  - Test connectivity: nc -zv <host> <port>
  - DNS lookup: dig <domain>

  IMPORTANT: This instance should be stopped when not in use!
================================================================================
EOF

echo "Bastion host setup complete!"
