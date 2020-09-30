variable "name" {
  description = "Prefix for resource names"
  type        = string
}

variable "k8s_version" {
  description = "Kubernetes version to install"
  type        = string
}

variable "vpc_id" {
  description = "VPC id to join to"
  type        = string
}

variable "subnets" {
  description = "List of subnets to use in EKS"
  type        = list(string)
}

variable "worker_groups" {
  type = list(object({
    name                 = string
    instance_type        = string
    asg_desired_capacity = number
    asg_min_size         = number
    asg_max_size         = number
  }))
}

variable "region" {
  description = "Region for AWS resources"
  type        = string
}

# The cluster autoscaler major and minor versions must match your cluster.
# For example if you are running a 1.16 EKS cluster set version to v1.16.5
# See https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/docs/autoscaling.md#notes
variable "autoscaler_version" {
  description = "EKS autoscaler image tag"
  type        = string
}

variable "autoscaler_name" {
  description = "EKS Autoscaler name"
  type        = string
}

variable "autoscaler_chart_version" {
  description = "EKS chart version"
  type        = string
}

# https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-scale-down-work
variable "autoscaler_scale_down_utilization_threshold" {
  description = "Node utilization level, defined as sum of requested resources divided by capacity"
  type        = string
}
