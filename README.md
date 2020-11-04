# m-aws-kubernetes-service

Epiphany Module: AWS Kubernetes Service

AwsKS module is reponsible for providing managed Kubernetes service [(Amazon EKS)](https://aws.amazon.com/eks/).

# Basic usage

## Prepare AWS access key

Have a look [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys).

## Build image

In main directory run:

  ```shell
  make build
  ```

or directly using Docker:

  ```shell
  cd m-aws-basic-infrastructure/
  docker build --tag epiphanyplatform/awsks:latest .
  ```

## Run module

The AwsKS cluster and new private subnets will be created in the vpc from the [AwsBI Module](https://github.com/epiphany-platform/m-aws-basic-infrastructure) or you can create the AwsKS cluster in an already existing subnets. Amazon EKS requires subnets in at least two Availability Zones. The existing VPC must meet specific [requirements](https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html) for use with Amazon EKS.

At this stage you should have already /tmp/shared directory with your ssh-keys and AwsBI module data.

* Initialize the AwsKS module on top of [AwsBI Module](https://github.com/epiphany-platform/m-aws-basic-infrastructure):

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest init
  ```

  This commad will create configuration file of AwsKS module in /tmp/shared/awsks/awsks-config.yml. You can investigate what is stored in that file.
  Available variables:
  * M_NAME = < module name > (default: epiphany)
  * M_REGION = < aws region > (default: eu-central-1)

  Note: M_REGION and M_NAME have to be the same as in AwsBI module

* Initialize the AwsKS module in already existing subnets:

  ```shell
  #TO-DO
  ```

* Plan and apply AwsKS module on top of [AwsBI Module](https://github.com/epiphany-platform/m-aws-basic-infrastructure):

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest plan M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx M_NAME=xxx
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest apply M_AWS_ACCESS_KEY=xxx M_AWS_SECRET_KEY=xxx M_NAME=xxx
  ```

  Running those commands should create EKS service. You can verify it in AWS Management Console.

* Share kubeconfig with `epicli` tool:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest kubeconfig
  ```

  This command will create file `/tmp/shared/kubeconfig`. You will need to move this file manually to `/tmp/shared/build/your-cluster-name/kubeconfig`.

## Run module with provided example

* Prepare your own variables in vars.mk file to use in the building process. Sample file (examples/basic_flow/vars.mk.sample):

  ```shell
  AWS_ACCESS_KEY_ID = "xxx"
  AWS_ACCESS_KEY_SECRET = "xxx"
  ```

* Create environment

  ```shell
  cd examples/basic_flow
  make all
  ```

  This command will create AWS Basic Infrastructure and AWS EKS on top of it.

## Run module with provided example in existing subnets

  ```shell
  #TO-DO
  ```

## Destroy EKS cluster

  ```shell
  cd examples/basic_flow
  make -k destroy
  ```

## Release module

  ```shell
  make release
  ```

or if you want to set different version number:

  ```shell
  make release VERSION=number_of_your_choice
  ```

## Notes

* The cluster autoscaler major and minor versions must match your cluster.
For example if you are running a 1.16 EKS cluster set version to v1.16.5.
For more details check [documentation](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/docs/autoscaling.md#notes)

## Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then ```examples/basic_flow_devcontainer``` folder.

## Module dependencies

| Component                     | Version | Repo/Website                                                                                                | License                                                           |
| ----------------------------- | ------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                     | 0.13.2  | https://www.terraform.io/                                                                                   | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider        | 3.7.0   | https://github.com/terraform-providers/terraform-provider-aws                                               | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Terraform Kubernetes provider | 1.13.2  | https://github.com/hashicorp/terraform-provider-kubernetes                                                  | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-kubernetes/blob/master/LICENSE) |
| Terraform Helm Provider       | 1.3.1   | https://github.com/hashicorp/terraform-provider-helm                                                        | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-helm/blob/master/LICENSE) |
| Terraform AWS EKS module      | 12.2.0  | https://github.com/terraform-aws-modules/terraform-aws-eks                                                  | [Apache License 2.0](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/LICENSE) |
| Terraform AWS IAM module      | 2.21.0  | https://github.com/terraform-aws-modules/terraform-aws-iam/tree/master/modules/iam-assumable-role-with-oidc | [Apache License 2.0](https://github.com/terraform-aws-modules/terraform-aws-iam/blob/master/LICENSE) |
| Make                          | 4.3     | https://www.gnu.org/software/make/                                                                          | [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                            | 3.3.4   | https://github.com/mikefarah/yq/                                                                            | [MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
