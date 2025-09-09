---
page_title: "MetalCloud: Core Concepts & Terminology"
description: |-
  Essential concepts and terminology for understanding MetalCloud's infrastructure management approach.
---

# Core Concepts & Terminology

MetalCloud uses a unique approach to infrastructure management that differs from traditional cloud providers. Understanding these core concepts is essential for effective use of the platform.

## Key Terminology

### 1. Infrastructure

A **Infrastructure** is the top-level organizational unit in MetalCloud. It serves as:

- A logical grouping for related resources (servers, drives, networks)
- A security boundary for access control
- An isolated environment for workloads

Users can create multiple infrastructures to separate different projects, environments, or teams.

### 2. ServerInstanceGroup

A **ServerInstanceGroup** is a collection of identical server instances that are managed as a single unit. Key characteristics:

- All instances in the group share identical configurations (OS template, drives, networks)
- Operations are performed at the group level, not individual instances
- Can be dynamically scaled up or down using the `instance_count` property
- Ideal for managing application clusters or load-balanced workloads

### 3. Drive

**Drives** are iSCSI LUNs (Logical Unit Numbers) that provide persistent storage:

- Can be attached to ServerInstanceGroups
- Shared across multiple instances within a group
- Persistent across instance lifecycle changes
- Essential for stateful applications like VMware, Kubernetes, databases

### 4. OSTemplate

An **OSTemplate** defines the operating system and initial configuration for server instances:

- Contains the base OS image and configuration scripts
- Applied to all instances in a ServerInstanceGroup
- Can include custom software, configurations, and automation scripts

### 5. LogicalNetwork

A **LogicalNetwork** is an abstraction layer for network connectivity:

- Implementation varies based on the underlying network fabric
- Spans multiple physical switches for redundancy
- Can be shared across multiple ServerInstanceGroups
- Supports various network profiles (VLAN, VXLAN, etc.)

### 6. Site

A **Site** represents a physical location or resource pool:

- Contains physical servers, storage systems, and network equipment
- Provides geographical distribution for workloads
- Enables resource allocation and capacity planning

## Operational Modes

### Deploy vs. Design Mode

MetalCloud uses a two-stage provisioning approach that separates planning from execution:

#### Design Mode (`prevent_deploy=true`)

- Make configuration changes without affecting running infrastructure
- Plan and validate changes before implementation
- Review resource requirements and dependencies
- Safe environment for experimentation

#### Deploy Mode (`prevent_deploy=false`)

- Apply all pending changes to the physical infrastructure
- Provision new resources and modify existing ones
- Changes are applied atomically across the entire infrastructure

This approach enables:

- **Change validation** before committing resources
- **Rollback capability** to revert to the last deployed state
- **Bulk operations** for efficient resource management

## Physical vs. Logical Separation

### Servers vs. Instances

MetalCloud maintains a clear separation between physical and logical resources:

#### Physical Server

- The actual hardware (CPU, RAM, storage, network interfaces)
- Managed by MetalCloud's orchestration layer
- Can be reassigned between workloads as needed

#### Instance

- The logical representation of compute resources
- Maintains identity across hardware changes
- Preserves configurations (firewall rules, DNS records, network settings)
- Can be "stopped" to release physical resources while retaining configuration

### Benefits of This Separation

1. **Hardware Flexibility**: Replace or upgrade physical servers without affecting workload configuration
2. **Resource Optimization**: Efficiently allocate physical resources based on demand
3. **Maintenance Windows**: Move workloads during hardware maintenance without reconfiguration
4. **Consistent Identity**: Maintain network and security configurations across hardware changes

## Resource Relationships

### ServerInstanceGroups and Drives

- **One-to-Many**: A ServerInstanceGroup can have multiple drives attached
- **Shared Access**: All instances in the group can access attached drives
- **Persistent Storage**: Drives maintain data across instance lifecycle changes
- **Performance**: Multiple instances can simultaneously access shared drives

### Instance Lifecycle Management

- **Provisioning**: Physical servers are allocated and configured
- **Running**: Instances are active and consuming resources
- **Stopped**: Physical resources released, but configuration preserved
- **Terminated**: Both physical resources and configuration removed

### Network Connectivity

- **LogicalNetworks** span multiple physical switches for redundancy
- **ServerInstanceGroups** can be connected to multiple networks
- **Network isolation** is maintained at the infrastructure level
- **Cross-infrastructure** communication requires explicit configuration

## Best Practices

1. **Use ServerInstanceGroups** for applications that benefit from horizontal scaling
2. **Separate concerns** using multiple infrastructures for different environments
3. **Plan changes** in design mode before deploying to production
4. **Monitor resource usage** across sites for optimal placement
5. **Implement proper drive management** for stateful applications

## Important Notes

> **Data Persistence**: When releasing servers, all local drive content is permanently wiped. Use attached drives for persistent data storage.
> **Network Security**: LogicalNetworks provide isolation, but proper firewall configuration is essential for security.
> **Resource Planning**: Consider site capacity and geographical requirements when designing infrastructures.
