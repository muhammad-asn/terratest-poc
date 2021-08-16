output "ami_instance_id" {
  value = data.aws_ami.ubuntu-bionic.id
}

output "instance_id" {
  value = module.ec2_instance_t2_micro.id[0]
}

output "subnet_id" {
  value = data.aws_subnet_ids.all.ids
}

output "public_address_instance" {
  value = module.ec2_instance_t2_micro.public_ip
}