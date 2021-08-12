terraform {
  required_version = ">=0.13"
}

# Configure IAM Policy EC2 role
module "policy_ec2_all_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-policy"
  version = "~> 3.0"

  name        = "policy-ec2-all-role"
  path        = "/"
  description = "Allow user to get all access to EC2 resource"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:*"
            ],
            "Resource": "arn:aws:ec2:*"
        }
    ]
}
EOF
}

