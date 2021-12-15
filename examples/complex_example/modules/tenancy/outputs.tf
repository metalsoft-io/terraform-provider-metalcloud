output "credentials" {
    value = module.tenancy_cluster[*].credentials
    sensitive = true
}

output "shared_drive_targets" {
    value = module.tenancy_cluster[*].shared_drive_targets
}