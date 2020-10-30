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

# Necessary for egress internet access from private networks
variable "public_subnet_id" {
  description = "Subnet id to attach NAT gateway to"
  type        = string
}

variable "private_route_table_id" {
  description = "Private route table id for table associations"
  type        = string
}
