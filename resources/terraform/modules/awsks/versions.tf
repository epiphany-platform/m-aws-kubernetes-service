terraform {
  required_version = ">= 0.13.2"

  # https://github.com/terraform-aws-modules/terraform-aws-eks#requirements
  required_providers {
    aws = {
      version = ">= 2.28.1"
    }

    kubernetes = {
      version = ">= 1.11.1"
    }

    local = {
      version = "~> 1.2"
    }

    null = {
      version = "~> 2.1"
    }

    random = {
      version = "~= 2.1"
    }

    template = {
      version = "~> 2.1"
    }
  }
}
