locals {
  bucket_name = "${var.project_name}-${var.environment}"

  # Build read access policy statements
  read_statements = flatten([
    for key, access in var.read_access : [
      {
        sid    = "ReadAccess${replace(title(key), "/[^a-zA-Z0-9]/", "")}List"
        effect = "Allow"
        principals = {
          type        = "AWS"
          identifiers = [access.principal]
        }
        actions   = ["s3:ListBucket"]
        resources = [aws_s3_bucket.this.arn]
        conditions = length(access.prefixes) > 0 ? [
          {
            test     = "StringLike"
            variable = "s3:prefix"
            values   = access.prefixes
          }
        ] : []
      },

      {
        sid    = "ReadAccess${replace(title(key), "/[^a-zA-Z0-9]/", "")}Objects"
        effect = "Allow"
        principals = {
          type        = "AWS"
          identifiers = [access.principal]
        }
        actions = [
          "s3:GetObject",
          "s3:GetObjectVersion"
        ]
        resources = flatten([
          for prefix in access.prefixes : "${aws_s3_bucket.this.arn}/${trimprefix(prefix, "/")}*"
        ])
        conditions = []
      }
    ]
  ])

  # Build write access policy statements
  write_statements = [
    for key, access in var.write_access : {
      sid    = "WriteAccess${replace(title(key), "/[^a-zA-Z0-9]/", "")}"
      effect = "Allow"
      principals = {
        type        = "AWS"
        identifiers = [access.principal]
      }
      actions = [
        "s3:PutObject",
        "s3:PutObjectAcl"
      ]
      resources = flatten([
        for prefix in access.prefixes : "${aws_s3_bucket.this.arn}/${trimprefix(prefix, "/")}*"
      ])
      conditions = []
    }
  ]

  # Build delete access policy statements
  delete_statements = [
    for key, access in var.delete_access : {
      sid    = "DeleteAccess${replace(title(key), "/[^a-zA-Z0-9]/", "")}"
      effect = "Allow"
      principals = {
        type        = "AWS"
        identifiers = [access.principal]
      }
      actions = [
        "s3:DeleteObject",
        "s3:DeleteObjectVersion"
      ]
      resources = flatten([
        for prefix in access.prefixes : "${aws_s3_bucket.this.arn}/${trimprefix(prefix, "/")}*"
      ])
      conditions = []
    }
  ]

  # Build full access policy statements
  full_access_statements = length(var.full_access_principals) > 0 ? [{
    sid    = "FullAccess"
    effect = "Allow"
    principals = {
      type        = "AWS"
      identifiers = var.full_access_principals
    }
    actions = [
      "s3:*"
    ]
    resources = [
      aws_s3_bucket.this.arn,
      "${aws_s3_bucket.this.arn}/*"
    ]
    conditions = []
  }] : []

  # Combine all policy statements
  all_statements = concat(
    local.read_statements,
    local.write_statements,
    local.delete_statements,
    local.full_access_statements
  )
}

# =========
# S3 Bucket
# =========

resource "aws_s3_bucket" "this" {
  bucket        = local.bucket_name
  force_destroy = var.force_destroy
}

# ===================
# Block Public Access
# ===================

resource "aws_s3_bucket_public_access_block" "this" {
  count = var.block_public_access ? 1 : 0

  bucket = aws_s3_bucket.this.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# ==========
# Versioning
# ==========

resource "aws_s3_bucket_versioning" "this" {
  count = var.enable_versioning ? 1 : 0

  bucket = aws_s3_bucket.this.id

  versioning_configuration {
    status = "Enabled"
  }
}

# ======================
# Server-Side Encryption
# ======================

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# ===============
# Lifecycle Rules
# ===============

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  count = length(var.lifecycle_rules) > 0 ? 1 : 0

  bucket = aws_s3_bucket.this.id

  dynamic "rule" {
    for_each = var.lifecycle_rules
    content {
      id     = rule.value.id
      status = rule.value.enabled ? "Enabled" : "Disabled"

      dynamic "filter" {
        for_each = rule.value.prefix != null ? [1] : []
        content {
          prefix = rule.value.prefix
        }
      }

      dynamic "expiration" {
        for_each = rule.value.expiration_days != null ? [1] : []
        content {
          days = rule.value.expiration_days
        }
      }

      dynamic "noncurrent_version_expiration" {
        for_each = rule.value.noncurrent_version_expiration != null ? [1] : []
        content {
          noncurrent_days = rule.value.noncurrent_version_expiration
        }
      }

      dynamic "transition" {
        for_each = rule.value.transition_days != null && rule.value.transition_storage_class != null ? [1] : []
        content {
          days          = rule.value.transition_days
          storage_class = rule.value.transition_storage_class
        }
      }
    }
  }
}

# =============
# Bucket Policy
# =============

resource "aws_s3_bucket_policy" "this" {
  count = length(local.all_statements) > 0 ? 1 : 0

  bucket = aws_s3_bucket.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      for stmt in local.all_statements : merge({
        Sid    = stmt.sid
        Effect = stmt.effect
        Principal = {
          (stmt.principals.type) = stmt.principals.identifiers
        }
        Action   = stmt.actions
        Resource = stmt.resources
        },
        # Conditionally add the Condition block only if it exists
        length(stmt.conditions) > 0 ? {
          Condition = {
            for cond in stmt.conditions :
            cond.test => {
              (cond.variable) = cond.values
            }
          }
      } : {})
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.this]
}
