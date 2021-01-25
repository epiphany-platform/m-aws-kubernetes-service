resource "aws_eks_node_group" "eks_nodes" {
  count           = length(var.worker_groups)
  cluster_name    = var.name
  node_group_name = var.worker_groups[count.index].name == null ? "${var.name}-node-group${count.index}" : var.worker_groups[count.index].name
  node_role_arn   = aws_iam_role.eks_nodes_iam_role.arn
  subnet_ids      = var.subnet_ids
  instance_types  = [var.worker_groups[count.index].instance_type]
  disk_size       = var.disk_size
  ami_type        = var.ami_type

  remote_access {
    ec2_ssh_key  = var.ec2_ssh_key
  }

  scaling_config {
    desired_size = var.worker_groups[count.index].asg_desired_capacity
    max_size     = var.worker_groups[count.index].asg_max_size
    min_size     = var.worker_groups[count.index].asg_min_size
  }

  # Ensure that IAM Role permissions are created before and deleted after EKS Node Group handling.
  # Otherwise, EKS will not be able to properly delete EC2 Instances and Elastic Network Interfaces.
  depends_on = [
    aws_iam_role_policy_attachment.AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.AmazonEC2ContainerRegistryReadOnly,
  ]

  # Allow external changes without Terraform plan difference
  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }

  # Add necessary tags for cluster autoscaler
  # https://docs.aws.amazon.com/eks/latest/userguide/cluster-autoscaler.html#ca-ng-considerations
  tags = local.eks_node_tags
}
