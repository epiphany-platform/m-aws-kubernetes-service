data "aws_vpc" "vpc" {
  id = var.vpc_id
}

module "eks" {
  source        = "modules/awsks"
  name          = var.name
  k8s_version   = var.k8s_version
  vpc_id        = data.aws_vpc.vpc.id
  subnets       = data.aws_vpc.vpc.private_subnets
  worker_groups = var.worker_groups
}
