output "credentials" {
    value = module.tenancy.credentials
    sensitive = true
}

output "shared_drive_targets" {
    value = module.tenancy.shared_drive_targets
}

