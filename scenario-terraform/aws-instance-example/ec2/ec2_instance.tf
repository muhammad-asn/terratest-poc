terraform {
  required_version = ">=0.13"
}

# Configure the AWS Provider
provider "aws" {
  region = "ap-southeast-1"
}

# AWS_VPC Default
data "aws_vpc" "default" {
  default = true
}

# AWS Subnet ID
data "aws_subnet_ids" "all" {
  vpc_id = data.aws_vpc.default.id
}

# Configure AMI (Ubuntu 18.04)
data "aws_ami" "ubuntu-bionic" {

  most_recent = true

  owners = ["099720109477"]

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "image-type"
    values = ["machine"]
  }

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }
}


# Configure the Instance EC2
module "ec2_instance_t2_micro" {
  source                 = "terraform-aws-modules/ec2-instance/aws"
  version                = "~> 2.0"
  name                   = var.ec2_instance_name
  ami                    = data.aws_ami.ubuntu-bionic.id
  instance_count         = 1
  instance_type          = var.ec2_instance_type
  key_name               = "aws-terraform-test-mac"
  subnet_id              = tolist(data.aws_subnet_ids.all.ids)[0]
  vpc_security_group_ids = ["sg-001fb6c4f008a67db"] # launch-wizard-10

  tags = {
    "Name" = var.ec2_instance_name
  }

}



