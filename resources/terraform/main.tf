data "aws_vpc" "vpc" {
  id = var.vpc_id
}

data "aws_route_table" "private_route_table" {
  route_table_id = var.private_route_table_id
}

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

resource "aws_route_table_association" "private1" {
  subnet_id      = aws_subnet.eks-subnet1.id
  route_table_id = data.aws_route_table.private_route_table.route_table_id
}

resource "aws_route_table_association" "private2" {
  subnet_id      = aws_subnet.eks-subnet2.id
  route_table_id = data.aws_route_table.private_route_table.route_table_id
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
  autoscaler_version                          = "v1.17.3"
  autoscaler_chart_version                    = "7.3.4"

}
