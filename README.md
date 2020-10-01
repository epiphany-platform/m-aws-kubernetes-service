# m-aws-kubernetes-service
Epiphany Module: AWS Kubernetes Service

# Prepare AWS access key

Have a look [here](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys).

# Build image

In main directory run:

```bash
make build
```

# Run module

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

# Destroy EKS cluster

```
cd examples/basic_flow
make -k destroy
```

# Release module

```bash
make release
```

or if you want to set different version number:

```bash
make release VERSION=number_of_your_choice
```

# Notes

- The cluster autoscaler major and minor versions must match your cluster.
For example if you are running a 1.16 EKS cluster set version to v1.16.5.
For more details check [documentation](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/docs/autoscaling.md#notes)
