# m-aws-kubernetes-service
Epiphany Module: AWS Kubernetes Service

## Prepare AWS access key

Have a look [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys).

## Build image

In main directory run:

```bash
make build
```

## Run module

```bash
cd examples/basic_flow
AWS_ACCESS_KEY="access key id" AWS_SECRET_KEY="access key secret" make all
```

Or use config file with credentials:

```bash
cd examples/basic_flow
cat >awsks.mk <<'EOF'
AWS_ACCESS_KEY ?= "access key id"
AWS_SECRET_KEY ?= "access key secret"
EOF
make all
```

## Destroy EKS cluster

```
cd examples/basic_flow
make -k destroy
```

## Release module

```bash
make release
```

or if you want to set different version number:

```bash
make release VERSION=number_of_your_choice
```

## Notes

- The cluster autoscaler major and minor versions must match your cluster.
For example if you are running a 1.16 EKS cluster set version to v1.16.5.
For more details check [documentation](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/docs/autoscaling.md#notes)

## Module dependencies

| Component                     | Version | Repo/Website                                                                                                | License                                                           |
| ----------------------------- | ------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                     | 0.13.2  | https://www.terraform.io/                                                                                   | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AWS provider        | 3.7.0   | https://github.com/terraform-providers/terraform-provider-aws                                               | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-aws/blob/master/LICENSE) |
| Terraform Kubernetes provider | 1.13.2  | https://github.com/hashicorp/terraform-provider-kubernetes                                                  | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-kubernetes/blob/master/LICENSE) |
| Terraform Helm Provider       | 1.3.1   | https://github.com/hashicorp/terraform-provider-helm                                                        | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform-provider-helm/blob/master/LICENSE) |
| Terraform AWS EKS module      | 12.2.0  | https://github.com/terraform-aws-modules/terraform-aws-eks                                                  | [Apache License 2.0](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/LICENSE) |
| Terraform AWS IAM module      | 2.21.0  | https://github.com/terraform-aws-modules/terraform-aws-iam/tree/master/modules/iam-assumable-role-with-oidc | [Apache License 2.0](https://github.com/terraform-aws-modules/terraform-aws-iam/blob/master/LICENSE) |
| Make                          | 4.3     | https://www.gnu.org/software/make/                                                                          | [ GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                            | 3.3.4   | https://github.com/mikefarah/yq/                                                                            | [ MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
