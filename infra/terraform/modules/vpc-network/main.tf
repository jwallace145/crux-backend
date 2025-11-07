resource "aws_vpc" "network" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support   = var.enable_dns_support

  tags = {
    Name = "${var.project_name}-network-${var.environment}"
  }
}

resource "aws_internet_gateway" "network_igw" {
  vpc_id = aws_vpc.network.id

  tags = {
    Name = "${var.project_name}-igw-${var.environment}"
  }
}

resource "aws_subnet" "public_subnets" {
  for_each = var.public_subnets

  vpc_id                  = aws_vpc.network.id
  availability_zone       = each.value.availability_zone
  cidr_block              = each.value.subnet_cidr_block
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-public-subnet-${lookup(local.availability_zone_abbreviations, each.value.availability_zone)}"
  }
}

resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.network.id

  tags = {
    Name = "${var.project_name}-public-rt"
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

# Private Subnets
resource "aws_subnet" "private_subnets" {
  for_each = var.private_subnets

  vpc_id                  = aws_vpc.network.id
  availability_zone       = each.value.availability_zone
  cidr_block              = each.value.subnet_cidr_block
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project_name}-private-subnet-${lookup(local.availability_zone_abbreviations, each.value.availability_zone)}"
  }
}

# Elastic IPs for NAT Gateways
resource "aws_eip" "nat_gateway_eips" {
  for_each = {
    for key, nat in var.nat_gateways : key => nat
    if nat.enabled
  }

  domain = "vpc"

  tags = {
    Name = "${var.project_name}-nat-eip-${lookup(local.availability_zone_abbreviations, each.value.availability_zone)}"
  }

  depends_on = [aws_internet_gateway.network_igw]
}

# NAT Gateways
resource "aws_nat_gateway" "nat_gateways" {
  for_each = {
    for key, nat in var.nat_gateways : key => nat
    if nat.enabled
  }

  allocation_id = aws_eip.nat_gateway_eips[each.key].id
  subnet_id = [
    for subnet in aws_subnet.public_subnets :
    subnet.id
    if subnet.availability_zone == each.value.availability_zone
  ][0]

  tags = {
    Name = "${var.project_name}-nat-gateway-${lookup(local.availability_zone_abbreviations, each.value.availability_zone)}"
  }

  depends_on = [aws_internet_gateway.network_igw]
}

# Private Route Tables (one per AZ with private subnets)
resource "aws_route_table" "private_route_tables" {
  for_each = toset([
    for subnet in var.private_subnets :
    subnet.availability_zone
  ])

  vpc_id = aws_vpc.network.id

  tags = {
    Name = "${var.project_name}-private-rt-${lookup(local.availability_zone_abbreviations, each.value)}"
  }
}

# Routes to NAT Gateways
resource "aws_route" "private_nat_gateway" {
  for_each = aws_route_table.private_route_tables

  route_table_id         = each.value.id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id = lookup(
    {
      for key, nat in aws_nat_gateway.nat_gateways :
      nat.tags["Name"] => nat.id
    },
    "${var.project_name}-nat-gateway-${lookup(local.availability_zone_abbreviations, each.key)}",
    # Fallback to first available NAT gateway if no NAT gateway in this AZ
    length(aws_nat_gateway.nat_gateways) > 0 ? values(aws_nat_gateway.nat_gateways)[0].id : null
  )
}

# Associate private subnets with their route tables
resource "aws_route_table_association" "private_route_table_association" {
  for_each = aws_subnet.private_subnets

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private_route_tables[each.value.availability_zone].id
}
