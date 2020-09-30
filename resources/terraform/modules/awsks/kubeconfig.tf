data "template_file" "kubeconfig" {
  template = file("${path.module}/templates/kubeconfig.tpl")
  vars     = {
    kubeconfig = module.eks.kubeconfig
  }
}

resource "local_file" "kubeconfig" {
  sensitive_content = data.template_file.kubeconfig.rendered
  filename          = "/workdir/.kube/config"
}
