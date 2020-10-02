variable "name" {
  description = "Prefix for resource names"
  type        = string
  default     = "default"
}

variable "vpc_id" {
  description = "VPC id to join to"
  type        = string
}

variable "worker_groups" {
  description = "Worker groups definition list"
  type        = list(object({
    name                 = string
    instance_type        = string
    asg_desired_capacity = number
    asg_min_size         = number
    asg_max_size         = number
  }))
  default     = [
    {
      name                 = "default_wg"
      instance_type        = "t2.small"
      asg_desired_capacity = 1
      asg_min_size         = 1
      asg_max_size         = 1
    }
  ]
}

variable "region" {
  description = "Region for AWS resources"
  type        = string
}

variable "autoscaler_name" {
  description = "EKS Autoscaler name"
  type        = string
  default     = "eks-autoscaler"
}

# https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-does-scale-down-work
variable "autoscaler_scale_down_utilization_threshold" {
  description = "Node utilization level, defined as sum of requested resources divided by capacity"
  type        = string
  default     = "0.65"
}

# Necessary for egress internet access from private networks
variable "public_subnet_id" {
  description = "Subnet id to attach NAT gateway to"
  type        = string
}
