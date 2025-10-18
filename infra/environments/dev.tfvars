/***************************************
 * CruxBackend Terraform Variables (dev)
 ***************************************/

service_name = "crux"
environment  = "dev"

network = {
  region         = "us-east-1"
  vpc_cidr_block = "10.0.0.0/22" # 2^10 = 1024 IPs

  # At least two subnets are required for database subnet group
  # Multi-AZ database disabled though to save costs
  public_subnets = {
    public_subnet_az1 = {
      availability_zone = "us-east-1a"
      subnet_cidr_block = "10.0.0.0/25" # 2^7 = 128 IPs
    }
    public_subnet_az2 = {
      availability_zone = "us-east-1b"
      subnet_cidr_block = "10.0.0.128/25" # 2^7 = 128 IPs
    }
  }
}

api = {
  task = {
    cpu    = 512
    memory = 1024
  }

  container = {
    cpu    = 512
    memory = 1024
    image  = "650503560686.dkr.ecr.us-east-1.amazonaws.com/crux-api:latest"
    port   = 3000
  }
}

database = {
  name                  = "cruxdb"
  username              = "cruxdbadmin"
  password              = "cruxdbpassword"
  instance_class        = "db.t4g.micro"
  postgres_version      = "16.3"
  allocated_storage     = 20
  max_allocated_storage = 50
}
