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
