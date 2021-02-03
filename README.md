# m-aws-kubernetes-service

Epiphany Module: AWS Kubernetes Service

AwsKS module is reponsible for providing managed Kubernetes service [(Amazon EKS)](https://aws.amazon.com/eks/).

# Basic usage

## Prepare AWS access key

Have a
look [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys).

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

The AwsKS cluster and new private subnets will be created in the vpc from
the [AwsBI Module](https://github.com/epiphany-platform/m-aws-basic-infrastructure) or you can create the AwsKS cluster
in an already existing subnets. Amazon EKS requires subnets in at least two Availability Zones. The existing VPC must
meet specific [requirements](https://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html) for use with Amazon
EKS.

At this stage you should have already /tmp/shared directory with your ssh-keys and AwsBI module data.

* Initialize the AwsKS module on top of [AwsBI Module](https://github.com/epiphany-platform/m-aws-basic-infrastructure):

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest init
  ```

  This commad will create configuration file of AwsKS module in /tmp/shared/awsks/awsks-config.yml. You can investigate
  what is stored in that file. Available parameters are listed in the [inputs](docs/INPUTS.adoc) document.

  Note: M_REGION and M_NAME have to be the same as in AwsBI module

* Initialize the AwsKS module in already existing subnets:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest init M_REGION="region of existing VPC" M_VPC_ID="existiing vpc id" M_SUBNET_IDS="[existing_subnet1_id,existing_subnet2_id,...]"
  ```

  This commad will create configuration file of AwsKS module in /tmp/shared/awsks/awsks-config.yml. You can investigate
  what is stored in that file. Available parameters are listed in the [inputs](docs/INPUTS.adoc) document.

* Plan and apply AwsKS module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest plan M_AWS_ACCESS_KEY="access key id" M_AWS_SECRET_KEY="access key secret"
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest apply M_AWS_ACCESS_KEY="access key id" M_AWS_SECRET_KEY="access key secret"
  ```

  Running those commands should create EKS service. You can verify it in AWS Management Console.

* Share kubeconfig with `epicli` tool:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/awsks:latest kubeconfig
  ```

  This command will create file `/tmp/shared/kubeconfig`. You will need to move this file manually
  to `/tmp/shared/build/your-cluster-name/kubeconfig`.

## Run module with provided example

* Prepare your own variables in vars.mk file to use in the building process. Sample file (
  examples/basic_flow/vars.mk.sample):

  ```shell
  AWS_ACCESS_KEY_ID = "access key id"
  AWS_ACCESS_KEY_SECRET = "access key secret"
  ```

* Create environment

  ```shell
  cd examples/basic_flow
  make all
  ```

  This command will create AWS Basic Infrastructure and AWS EKS on top of it.

## Run module with provided example in existing subnets

* Prepare your own variables in vars.mk file to use in the building process. Sample file (
  examples/create_in_existing_subnets/vars.mk.sample):

  ```shell
  AWS_ACCESS_KEY = "access key id"
  AWS_SECRET_KEY = "access key secret"
  M_REGION="region of existing VPC"
  M_VPC_ID="existing vpc id"
  M_SUBNET_IDS="[existing_subnet1_id,existing_subnet2_id,...]"
  ```

* Create environment

  ```shell
  cd examples/create_in_existing_subnets
  make all
  ```

  This command will create AwsKS cluster in already existing subnets.

## Destroy EKS cluster

  ```shell
  cd examples/basic_flow
  make destroy
  ```

## Release module

  ```shell
  make release
  ```

or if you want to set different version number:

  ```shell
  make release VERSION=number_of_your_choice
  ```

## Integration tests execution

Prior to run integration tests specify environment variables:

### Required environment variables

- AWS_ACCESS_KEY_ID - your access key
- AWS_SECRET_ACCESS_KEY - your secret
- AWSBI_IMAGE_TAG - full tag of docker image that you want to test e.g. "epiphanyplatform/awsbi:0.0.1"
- AWSKS_IMAGE_TAG - full tag of docker image that you want to test e.g. "epiphanyplatform/awsks:0.0.1"

### Optional environment variables

- M_NAME - prefix to create AWS resources with
- M_REGION - AWS region to create resources in

and after that run shell command:

```shell
  make test
```

## Notes

* The cluster autoscaler major and minor versions must match your cluster. For example if you are running a 1.16 EKS
  cluster set version to v1.16.5. For more details
  check [documentation](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/docs/autoscaling.md#notes)

## Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the
included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the
examples the same way but then from then ```examples/basic_flow_devcontainer``` folder
or ```examples/create_in_existing_subnets_devcontainer```.

## Module dependencies

| Component                       | Version | Repo/Website                                                                                                | License                                                           |
| ------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                       | 0.13.2  | https://www.terraform.io/                                                                                   | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider          | 3.21.0   | https://github.com/terraform-providers/terraform-provider-aws                                               | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Terraform Kubernetes provider   | 1.13.3  | https://github.com/hashicorp/terraform-provider-kubernetes                                                  | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-kubernetes/blob/master/LICENSE) |
| Terraform Helm Provider         | 1.3.1   | https://github.com/hashicorp/terraform-provider-helm                                                        | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-helm/blob/master/LICENSE) |
| Terraform TLS provider          | 3.0.0   | https://github.com/hashicorp/terraform-provider-tls                                                         | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-tls/blob/master/LICENSE) |
| Terraform Template Provider     | 2.2.0   | https://github.com/hashicorp/terraform-provider-template                                                    | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-template/blob/master/LICENSE) |
| Terraform Metrics Server Module | 0.9.0   | https://github.com/cookielab/terraform-kubernetes-metrics-server                                            | [MIT License](https://github.com/cookielab/terraform-kubernetes-metrics-server/blob/master/LICENSE.md) |
| Cluster Autoscaler Helm Chart   | 7.3.4   | https://github.com/helm/charts/tree/master/stable/cluster-autoscaler (deprecated)                           | [Apache License 2.0](https://github.com/kubernetes/autoscaler/blob/master/LICENSE) |
| Make                            | 4.3     | https://www.gnu.org/software/make/                                                                          | [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                              | 3.3.4   | https://github.com/mikefarah/yq/                                                                            | [MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
