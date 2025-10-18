variable "repository_name" {
  description = "Name of the ECR repository"
  type        = string
}

variable "image_tag_mutability" {
  description = "The tag mutability setting for the repository (MUTABLE or IMMUTABLE)"
  type        = string
  default     = "MUTABLE"

  validation {
    condition     = contains(["MUTABLE", "IMMUTABLE"], var.image_tag_mutability)
    error_message = "image_tag_mutability must be either MUTABLE or IMMUTABLE"
  }
}

variable "scan_on_push" {
  description = "Enable image scanning on push"
  type        = bool
  default     = true
}

variable "encryption_type" {
  description = "The encryption type to use for the repository (AES256 or KMS)"
  type        = string
  default     = "AES256"

  validation {
    condition     = contains(["AES256", "KMS"], var.encryption_type)
    error_message = "encryption_type must be either AES256 or KMS"
  }
}

variable "kms_key_arn" {
  description = "The ARN of the KMS key to use for encryption (only used if encryption_type is KMS)"
  type        = string
  default     = null
}

variable "lifecycle_policy_enabled" {
  description = "Enable lifecycle policy for image retention"
  type        = bool
  default     = true
}

variable "image_count_to_keep" {
  description = "Number of tagged images to keep"
  type        = number
  default     = 10
}

variable "tag_prefixes_to_keep" {
  description = "List of tag prefixes to apply retention policy to"
  type        = list(string)
  default     = ["dev", "prod", "staging"]
}

variable "untagged_days_to_keep" {
  description = "Number of days to keep untagged images before deletion"
  type        = number
  default     = 7
}

variable "repository_policy" {
  description = "JSON policy document for repository access control"
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags to apply to the ECR repository"
  type        = map(string)
  default     = {}
}
