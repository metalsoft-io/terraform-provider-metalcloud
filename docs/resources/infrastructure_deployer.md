---
page_title: "metalcloud_infrastructure_deployer Resource - terraform-provider-metalcloud"
description: |-
  Controls the deployment of MetalCloud infrastructure and all its elements such as instance arrays and networks.
---

# metalcloud_infrastructure_deployer

The `metalcloud_infrastructure_deployer` is the central resource for managing MetalCloud infrastructure deployment. It orchestrates the provisioning and lifecycle management of all infrastructure components including:

* Deployment control flags and configuration options
* One or more [instance_array](./instance_array.html.md) blocks (ServerInstanceGroups)
* One or more [network](./network.html.md) blocks (LogicalNetworks)
* Custom variables and workflow tasks

## Key Concepts

### Design vs. Deploy Mode

MetalCloud uses a two-stage provisioning approach:

- **Design Mode** (`prevent_deploy = true`): Plan and validate changes without affecting physical infrastructure
- **Deploy Mode** (`prevent_deploy = false`): Apply all pending changes to provision actual resources

This separation allows you to safely prepare infrastructure configurations and deploy them atomically.

## Example Usage

### Basic Infrastructure Deployment

```terraform
# Retrieve existing infrastructure
data "metalcloud_infrastructure" "infra" {
    infrastructure_label = "my-app-prod"
    datacenter_name = "dc-west-1"
}

# Deploy the infrastructure
resource "metalcloud_infrastructure_deployer" "infrastructure_deployer" {
    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

    # Deployment control
    prevent_deploy = false
    await_deploy_finished = true
    
    # Safety settings
    allow_data_loss = false
    
    # Shutdown behavior
    attempt_soft_shutdown = true
    soft_shutdown_timeout_seconds = 300
    hard_shutdown_after_timeout = true

    # Ensure all infrastructure components are created before deployment
    depends_on = [
        metalcloud_instance_array.web_servers,
        metalcloud_instance_array.database_servers,
        metalcloud_network.private_network
    ]
}
```

### Development/Testing Configuration

```terraform
resource "metalcloud_infrastructure_deployer" "dev_infrastructure" {
    infrastructure_id = data.metalcloud_infrastructure.dev.infrastructure_id

    # Prevent actual deployment for testing
    prevent_deploy = true
    
    # Allow data loss for testing scenarios
    allow_data_loss = true
    
    # Custom variables for development
    infrastructure_custom_variables = {
        environment = "development"
        debug_mode = "enabled"
        log_level = "debug"
    }

    depends_on = [
        metalcloud_instance_array.test_cluster
    ]
}
```

### Production with Workflow Tasks

```terraform
resource "metalcloud_infrastructure_deployer" "prod_infrastructure" {
    infrastructure_id = data.metalcloud_infrastructure.prod.infrastructure_id

    prevent_deploy = false
    await_deploy_finished = true
    allow_data_loss = false

    # Custom variables for production environment
    infrastructure_custom_variables = {
        environment = "production"
        backup_enabled = "true"
        monitoring_endpoint = "https://monitoring.company.com"
    }

    # Post-deployment workflow tasks
    workflow_task {
        stage_definition_id = data.metalcloud_workflow_task.install_monitoring.id
        run_level = 1
        stage_run_group = "post_deploy"
    }

    workflow_task {
        stage_definition_id = data.metalcloud_workflow_task.configure_backup.id
        run_level = 2
        stage_run_group = "post_deploy"
    }

    depends_on = [
        metalcloud_instance_array.app_servers,
        metalcloud_instance_array.db_cluster,
        metalcloud_network.app_network,
        metalcloud_network.db_network
    ]
}
```

## Argument Reference

### Required Arguments

* `infrastructure_id` - (Required) The ID of the infrastructure to deploy. Use the `metalcloud_infrastructure` data source to retrieve this ID.

### Deployment Control

* `prevent_deploy` - (Optional, default: `true`) Controls whether to actually provision physical resources:
  - `true`: Design mode - validate configuration without provisioning
  - `false`: Deploy mode - provision actual infrastructure
  
* `await_deploy_finished` - (Optional, default: `true`) Whether to wait for deployment completion:
  - `true`: Terraform waits until deployment finishes before continuing
  - `false`: Terraform continues while deployment runs in background

### Safety and Data Protection

* `allow_data_loss` - (Optional, default: `false`) Controls operations that may cause data loss:
  - `true`: Allow destructive operations (stopping/deleting drives)
  - `false`: Block operations that could cause data loss
  - **⚠️ Use with caution in production environments**

* `keep_detaching_drives` - (Optional, default: `true`) Behavior when reducing instance count:
  - `true`: Preserve drives from detached instances
  - `false`: Delete drives from detached instances

### Shutdown Behavior

* `attempt_soft_shutdown` - (Optional, default: `true`) Whether to attempt graceful shutdown:
  - `true`: Send ACPI shutdown signal first
  - `false`: Perform immediate hard shutdown

* `soft_shutdown_timeout_seconds` - (Optional, default: `180`) Time to wait for graceful shutdown before forcing

* `hard_shutdown_after_timeout` - (Optional, default: `true`) Action when soft shutdown times out:
  - `true`: Force hard shutdown after timeout
  - `false`: Wait indefinitely (requires manual intervention)

### Advanced Configuration

* `infrastructure_custom_variables` - (Optional) Key-value pairs passed to OS templates and workflows:
  ```terraform
  infrastructure_custom_variables = {
      environment = "production"
      app_version = "v2.1.0"
      backup_schedule = "daily"
  }
  ```

* `skip_ansible` - (Optional, default: `false`) Skip automatic provisioning steps (advanced use only)

### Workflow Integration

* `workflow_task` - (Optional) Define post-deployment automation tasks:
  ```terraform
  workflow_task {
      stage_definition_id = data.metalcloud_workflow_task.task_name.id
      run_level = 1  # Execution order (lower numbers run first)
      stage_run_group = "post_deploy"  # or "pre_deploy"
  }
  ```

### Deprecated Arguments

* `server_allocation_policy` - (DEPRECATED) Use instance array configurations instead

## Attributes Reference

* `infrastructure_id` - The infrastructure ID, also used as the resource ID

## Important Considerations

### Dependencies

Always use `depends_on` to ensure all infrastructure components are created before deployment:

```terraform
depends_on = [
    metalcloud_instance_array.component1,
    metalcloud_instance_array.component2,
    metalcloud_network.network1
]
```

### Data Safety

- Set `allow_data_loss = false` in production to prevent accidental data loss
- Use `prevent_deploy = true` when testing configurations
- Always backup critical data before major infrastructure changes

### Deployment Timing

- Use `await_deploy_finished = true` for sequential deployments
- Set `await_deploy_finished = false` for parallel deployments (advanced)
- Configure appropriate shutdown timeouts for your workloads

### Best Practices

1. **Start with Design Mode**: Always test with `prevent_deploy = true` first
2. **Use Custom Variables**: Pass environment-specific configuration through `infrastructure_custom_variables`
3. **Implement Proper Dependencies**: Use `depends_on` to ensure correct provisioning order
4. **Plan for Data Safety**: Configure appropriate `allow_data_loss` and shutdown settings
5. **Monitor Deployments**: Use `await_deploy_finished = true` to track deployment progress

## Related Resources

- [metalcloud_infrastructure](../data-sources/infrastructure.html.md) - Data source for retrieving infrastructure information
- [metalcloud_instance_array](./instance_array.html.md) - Define server instance groups
- [metalcloud_network](./network.html.md) - Configure logical networks
- [Core Concepts Guide](../guides/concepts.html.md) - Understanding MetalCloud terminology and concepts
