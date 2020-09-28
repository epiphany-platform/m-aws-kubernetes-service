output "kubeconfig" {
  description = "Kubeconfig as generated by the module"
  value       = module.eks.kubeconfig
  sensitive   = true
}
