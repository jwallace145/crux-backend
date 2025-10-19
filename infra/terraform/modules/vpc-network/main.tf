resource "aws_vpc" "network" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support   = var.enable_dns_support

  tags = {
    Name = "${var.service_name}-network-${var.environment}"
  }
}

resource "aws_internet_gateway" "network_igw" {
  vpc_id = aws_vpc.network.id

  tags = {
    Name = "${var.service_name}-igw-${var.environment}"
  }
}

resource "aws_subnet" "public_subnets" {
  for_each = var.public_subnets

  vpc_id                  = aws_vpc.network.id
  availability_zone       = each.value.availability_zone
  cidr_block              = each.value.subnet_cidr_block
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.service_name}-public-subnet-${lookup(local.availability_zone_abbreviations, each.value.availability_zone)}"
  }
}

resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.network.id

  tags = {
    Name = "${var.service_name}-public-rt"
  }
}

resource "aws_route" "public_internet_gateway" {
  route_table_id         = aws_route_table.public_route_table.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.network_igw.id
}

resource "aws_route_table_association" "public_route_table_association" {
  for_each       = aws_subnet.public_subnets
  subnet_id      = each.value.id
  route_table_id = aws_route_table.public_route_table.id
}
