locals {
  worker_groups = [
    for wg in var.worker_groups:
    {
      name                    = wg.name
      instance_type           = wg.instance_type
      asg_desired_capacity    = wg.asg_desired_capacity
      asg_min_size            = wg.asg_min_size
      asg_max_size            = wg.asg_max_size
      tags = [
        {
          key                 = "k8s.io/cluster-autoscaler/enabled"
          propagate_at_launch = "false"
          value               = "true"
        },
        {
          key                 = "k8s.io/cluster-autoscaler/${module.eks.cluster_name}"
          propagate_at_launch = "false"
          value               = "true"
        }
      ]
    }
  ]
}

module "eks" {
  # there is no information about EKS addition to ARG
  # neither in module: https://github.com/terraform-aws-modules/terraform-aws-eks
  # nor in AWS documentation: https://docs.aws.amazon.com/ARG/latest/userguide/supported-resources.html
  source          = "terraform-aws-modules/eks/aws"
  version         = "12.2.0"
  cluster_name    = "${var.name}-eks"
  cluster_version = var.k8s_version
  subnets         = var.subnets
  vpc_id          = var.vpc_id
  worker_groups   = local.worker_groups

  tags = {
    Environment = var.name
  }
}

data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}
