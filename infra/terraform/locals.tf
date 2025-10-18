locals {
  valid_regions = ["us-east-1"]

  valid_availability_zones = {
    "us-east-1" : [
      "us-east-1a",
      "us-east-1b",
      "us-east-1c",
      "us-east-1d",
      "us-east-1e",
    ]
  }

  valid_environments = [
    "dev",
    "stg",
    "prod"
  ]
}
