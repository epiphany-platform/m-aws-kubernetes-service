data "aws_caller_identity" "current" {}

resource "helm_release" "cluster-autoscaler" {
  name            = "cluster-autoscaler"
  chart           = "stable/cluster-autoscaler"
  version         = var.autoscaler_chart_version
  cleanup_on_fail = "true"
  namespace       = "kube-system"
  timeout         = 300

  set {
    name  = "cloudProvider"
    type  = "string"
    value = "aws"
  }
  set {
    name  = "awsRegion"
    type  = "string"
    value = var.region
  }
  set {
    name  = "autoDiscovery.clusterName"
    type  = "string"
    value = "${var.name}-eks"
  }
  set {
    name  = "autoDiscovery.enabled"
    type  = "string"
    value = "true"
  }
  set {
    name  = "image.repository"
    type  = "string"
    value = "k8s.gcr.io/autoscaling/cluster-autoscaler"
  }
  set {
    name  = "image.tag"
    type  = "string"
    value = var.autoscaler_version
  }
  set {
    name  = "extraArgs.scale-down-utilization-threshold"
    type  = "auto"
    value = var.autoscaler_scale_down_utilization_threshold
  }
  set {
    name  = "rbac.serviceAccountAnnotations.eks\\.amazonaws\\.com/role-arn"
    type  = "string"
    value = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.autoscaler_name}"
  }
}
