locals {
  k8s_service_account_namespace               = "kube-system"
  k8s_service_account_name                    = "cluster-autoscaler-aws-cluster-autoscaler"

  tags = map(
    "resource_group", var.name
  )
}
