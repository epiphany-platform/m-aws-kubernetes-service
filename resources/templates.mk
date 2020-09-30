define M_METADATA_CONTENT
labels:
  version: $(M_VERSION)
  name: AWS Kubernetes Service
  short: $(M_MODULE_SHORT)
  kind: infrastructure
  provider: aws
endef

define M_CONFIG_CONTENT
kind: $(M_MODULE_SHORT)-config
$(M_MODULE_SHORT):
  name: $(M_NAME)
  k8s_version: $(M_K8S_VERSION)
  vpc_id: $(M_VPC_ID)
  region: $(M_REGION)
  autoscaler_version: $(M_K8S_AUTOSCALER_VERSION)
  public_subnet_id: $(M_PUBLIC_SUBNET_ID)
endef

define M_STATE_INITIAL
kind: state
$(M_MODULE_SHORT):
  status: initialized
endef
