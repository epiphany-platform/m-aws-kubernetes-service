data "aws_vpc" "vpc" {
  id = var.vpc_id
}

# Need to create private subnet as it's not done by awsbi module.
# When "aws_subnet_ids" datasource is used in a such way
#
# data "aws_subnet_ids" "private" {
#  vpc_id = var.vpc_id
#  tags = {
#    Tier = "Private"
#  }
#}
#
# there is not an empty result, but an error:
# https://github.com/hashicorp/terraform/issues/16380

# Subnets in at least 2 availability zones are required
resource "aws_subnet" "eks-subnet1" {
  vpc_id     = data.aws_vpc.vpc.id
  cidr_block = cidrsubnet(data.aws_vpc.vpc.cidr_block, 4, 14)
  availability_zone = "${var.region}a"
  tags       = {
    Name                                    = "${var.name}-eks-subnet1"
    cluster_name                            = var.name
    # https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-subnet-tagging
    "kubernetes.io/cluster/${var.name}-eks" = "shared"
    "kubernetes.io/role/internal-elb"       = 1
  }
}

resource "aws_subnet" "eks-subnet2" {
  vpc_id     = data.aws_vpc.vpc.id
  cidr_block = cidrsubnet(data.aws_vpc.vpc.cidr_block, 4, 15)
  availability_zone = "${var.region}b"
  tags       = {
    Name                                    = "${var.name}-eks-subnet2"
    cluster_name                            = var.name
    # https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-subnet-tagging
    "kubernetes.io/cluster/${var.name}-eks" = "shared"
    "kubernetes.io/role/internal-elb"       = 1
  }
}

# https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-tagging
resource "aws_ec2_tag" "eks-vpc" {
  resource_id = data.aws_vpc.vpc.id
  key         = "kubernetes.io/cluster/${var.name}-eks"
  value       = "shared"
}

module "awsks" {
  source        = "./modules/awsks"
  name          = var.name
  k8s_version   = var.k8s_version
  vpc_id        = data.aws_vpc.vpc.id
  subnets       = [aws_subnet.eks-subnet1.id,aws_subnet.eks-subnet2.id]
  worker_groups = var.worker_groups
}
