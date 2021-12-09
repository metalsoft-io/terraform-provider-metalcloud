output "credentials" {
    value = module.tenancy_cluster[*].credentials
    sensitive = true
}