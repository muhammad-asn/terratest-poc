terraform {
  required_version = ">=0.13"
}

/**
Configure IAM Create User name defined

Rules: 
- name: <name>-mania
**/
module "bonek-mania" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-user"
  version = "~> 3.0"

  name                    = "bonek-mania"
  force_destroy           = "true"
  pgp_key                 = "keybase:whitespace01"
  password_reset_required = true
  create_iam_access_key   = false
  tags = {
    Organization    = "bonek"
    Project         = "bonek"
    Terraform       = "true"
    Email           = "bonek-mania@mania.com"
    BastionUserName = "bonek-mania"
  }
}


# module "jamrud" {
#   source  = "terraform-aws-modules/iam/aws//modules/iam-user"
#   version = "~> 3.0"

#   name                    = "jamrud"
#   force_destroy           = "true"
#   pgp_key                 = "keybase:whitespace01"
#   password_reset_required = true
#   create_iam_access_key   = false
#   tags = {
#     Organization    = "jamrud"
#     Project         = "jamrud"
#     Terraform       = "true"
#     Email           = "jamrud@jamrud.com"
#     BastionUserName = "jamrud-jamrud"
#   }
# }


