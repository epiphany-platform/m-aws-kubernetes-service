terraform {
  required_version = "0.13.2"

  # https://github.com/terraform-aws-modules/terraform-aws-eks#requirements
  required_providers {
    aws = {
      version = "3.7.0"
    }

    kubernetes = {
      version = "1.13.2"
    }

    helm = {
      version = "1.3.1"
    }

    local = {
      version = "1.4"
    }

    null = {
      version = "2.1"
    }

    random = {
      version = "2.1"
    }

    template = {
      version = "2.1"
    }
  }
}
