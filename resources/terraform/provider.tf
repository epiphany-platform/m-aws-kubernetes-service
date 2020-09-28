provider "aws" {
  # No need to pass it as a variable as region is defined in m-aws-basic-infrastructure
  # and this default must be set to some value
  region = "eu-central-1"
}
