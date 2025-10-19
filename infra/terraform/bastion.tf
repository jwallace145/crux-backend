# ============
# Bastion Host
# ============

module "bastion" {
  count  = var.bastion.enabled ? 1 : 0
  source = "./modules/ec2-bastion-host"

  # Service details
  service_name = var.service_name
  environment  = var.environment

  # Network configs
  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  # Allow SSH from anywhere (use with caution - consider restricting to your IP)
  allowed_ssh_cidr_blocks = ["0.0.0.0/0"]

  # Static IP for consistent access
  allocate_elastic_ip = true

  # SSH public key
  ssh_public_key_path = "./id_rsa.pub"

  tags = {
    Project = "Crux"
  }
}
