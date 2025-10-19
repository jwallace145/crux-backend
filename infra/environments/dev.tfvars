/***************************************
 * CruxProject Terraform Variables (dev)
 ***************************************/

service_name = "crux"
environment  = "dev"

network = {
  region         = "us-east-1"
  vpc_cidr_block = "10.0.0.0/22" # 2^10 = 1024 IPs

  public_subnets = {
    public_subnet_az1 = {
      availability_zone = "us-east-1a"
      subnet_cidr_block = "10.0.0.0/25" # 2^7 = 128 IPs
    }
    public_subnet_az2 = {
      availability_zone = "us-east-1b"
      subnet_cidr_block = "10.0.0.128/25" # 2^7 = 128 IPs
    }
    public_subnet_az3 = {
      availability_zone = "us-east-1c"
      subnet_cidr_block = "10.0.1.0/25" # 2^7 = 128 IPs
    }
  }
}

api = {
  domain = "cruxproject.io"

  task = {
    cpu    = 256
    memory = 512
  }

  container = {
    cpu    = 256
    memory = 512
    image  = "650503560686.dkr.ecr.us-east-1.amazonaws.com/crux-api:latest"
    port   = 3000
  }
}

database = {
  postgres_version = "16.3"
  instance_class   = "db.t4g.micro"

  name = "cruxdb"

  user = {
    username = "cruxdbadmin"
    password = "cruxdbpassword"
  }

  storage = {
    allocated_storage     = 20
    max_allocated_storage = 50
  }
}

bastion = {
  enabled       = false
  instance_type = "t3.micro"

  # Allow SSH from anywhere
  allowed_ssh_cidr_blocks = ["0.0.0.0/0"]
}
