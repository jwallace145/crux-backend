output "user_arn" {
  value       = aws_iam_user.cicd_user.arn
  description = "The ARN of the CICD user."
}

output "user_name" {
  value       = aws_iam_user.cicd_user.name
  description = "The name of the CICD user."
}
