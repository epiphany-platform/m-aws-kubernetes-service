terraform {
  required_version = "0.13.2"

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

    tls = {
      version = "3.0.0"
    }

    template = {
      version = "2.2.0"
    }
  }
}
