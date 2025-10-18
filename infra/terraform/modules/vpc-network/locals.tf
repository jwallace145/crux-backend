locals {
  valid_regions = [
    "us-east-1"
  ]

  valid_availability_zones = {
    "us-east-1" : [
      "us-east-1a",
      "us-east-1b",
      "us-east-1c",
      "us-east-1d",
      "us-east-1e",
    ]
  }

  availability_zone_abbreviations = {
    "us-east-1a" = "use1a"
    "us-east-1b" = "use1b"
    "us-east-1c" = "use1c"
    "us-east-1d" = "use1d"
  }
}
