data "aws_vpc" "vpc" {
  id = var.vpc_id
}

# Need to create private subnets as it's not done by awsbi module.
# When "aws_subnet_ids" datasource is used in a such way
#
# data "aws_subnet_ids" "private" {
#  vpc_id = var.vpc_id
#  tags = {
#    Tier = "Private"
#  }
#}
#
# and there is no private subnets in VPC, result is not empty, but an error:
# https://github.com/hashicorp/terraform/issues/16380

# Subnets in at least 2 availability zones are required for EKS

# Following part could be created in aws-basic-infrastructure module
# ----------------------------------------------------------------------------------------------------------------------
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

data "aws_route_table" "awsbi_route_table_private" {
  tags = {
    Name         = "rt-private-${var.name}"
    cluster_name = var.name
  }
}

resource "aws_route_table_association" "private1" {
  subnet_id      = aws_subnet.eks-subnet1.id
  route_table_id = data.aws_route_table.awsbi_route_table_private.id
}

resource "aws_route_table_association" "private2" {
  subnet_id      = aws_subnet.eks-subnet2.id
  route_table_id = data.aws_route_table.awsbi_route_table_private.id
}

# https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html#vpc-tagging
resource "aws_ec2_tag" "eks-vpc" {
  resource_id = data.aws_vpc.vpc.id
  key         = "kubernetes.io/cluster/${var.name}-eks"
  value       = "shared"
}
# ----------------------------------------------------------------------------------------------------------------------

module "awsks" {
  source                                      = "./modules/awsks"
  name                                        = var.name
  k8s_version                                 = "1.17"
  vpc_id                                      = data.aws_vpc.vpc.id
  subnets                                     = [aws_subnet.eks-subnet1.id,aws_subnet.eks-subnet2.id]
  worker_groups                               = var.worker_groups
  region                                      = var.region
  autoscaler_name                             = var.autoscaler_name
  autoscaler_version                          = "v1.17.3"
  autoscaler_chart_version                    = "7.3.4"
  autoscaler_scale_down_utilization_threshold = var.autoscaler_scale_down_utilization_threshold
}
