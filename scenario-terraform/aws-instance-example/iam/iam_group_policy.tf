terraform {
  required_version = ">=0.13"
}

module "group_policy_ec2_all_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-group-with-policies"
  version = "~> 3.0"

  name = "group-policy-ec2-all-role"

  group_users = [
    "bonek-mania"
  ]

  attach_iam_self_management_policy = false

  custom_group_policy_arns = [
    "arn:aws:iam::783095911817:policy/policy-ec2-all-role"
  ]

  depends_on = [
    module.policy_ec2_all_role
  ]
}