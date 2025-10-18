/*********************
 * CruxBackend Network
 *********************/

module "network" {
  source         = "./modules/vpc-network"
  service_name   = var.service_name
  environment    = var.environment
  region         = var.network.region
  vpc_cidr_block = var.network.vpc_cidr_block
  public_subnets = var.network.public_subnets
}
