variable "infrastructure_id" {
  
   description = "infrastructure to use"
   type = number
   default = 0

    validation {
      condition     = var.infrastructure_id!=0
      error_message = "The infrastructure_id must be provided. Use a metalsoft_infrastructure data block to retrieve it."
    }
}


/*
variable "infrastructure_label" {
  
   description = "infrastructure to use"
   type = string
   default = ""

    validation {
      condition     = var.infrastructure_label!=""
      error_message = "The infrastructure_label must be provided. If the infrastructure does not exist it will be created."
    }
}

variable "datacenter_name" {
  
   description = "Datacenter to use."
   type = string
   default = ""

    validation {
      condition     = var.datacenter_name!=""
      error_message = "The datacenter_name must be provided."
    }
}
*/
variable "clustername" {
   description = "Cluster's name"
   type = string
   default = "cluster"
}


variable "customer_prefix" {
   description = "IP subnet to use in CIDR notation"
   type = string
   default = ""
}

variable "compute_nodes" {
  description = "array of names of compute nodes"
  type = list(object({
    compute_node_name = string
  }))
  default=[]
}

variable "datastores" {
  description = "datastore name and size"
   type = list(object({
    datastore_name = string
    shared_drive_size = string
  }))
  default=[]
}


variable "instance_array_ram_gbytes" {
  description = "Minimum amount of RAM required for instances"
  type = number
  default = 2
}

variable "instance_array_processor_count" {
  description = "Minimum amount of CPU count"
  type = number
  default = 1
}

variable "instance_array_processor_core_count" {
  description = "Minimum amount of HT cores per cpu"
  type = number
  default = 1
}





