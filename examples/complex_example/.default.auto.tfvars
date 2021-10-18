tenancy_config = {
    "customer_name" = "test-114",
    "datacenter" = "us-chi-qts01-dc",

     "esxi_vlan_id"= 101
     "vmotion_vlan_id"= 201
     "tep_vlan_id" = 301
     "pri_transit_a_vlan_id" = 401
     "pri_transit_b_vlan_id" = 501
     "pub_transit_a_vlan_id" = 601
     "pub_transit_b_vlan_id" = 701

    "clusters" = [
      {
        "clustername" = "cluster01",
        "customer_prefix" = "10.120.70.0/24",
        "compute_node_size" = "small",

        "compute_nodes" = [
          {
            "compute_node_name" : "ubuntu-db-01-2"
          },
          {
            "compute_node_name" : "ubuntu-db-02"
          }
        ],
        "datastores" = [
          {
            "datastore_name" : "datastore-db-01",
            "shared_drive_size" = "2TB"
          },
          {
            "datastore_name" : "datastore-db-02",
            "shared_drive_size" = "4TB"
          },
          {
            "datastore_name" : "datastore-db-03",
            "shared_drive_size" = "6TB"
          }
        ],
        "instance_array_ram_gbytes" = "2",
        "instance_array_processor_count" = "1",
        "instance_array_processor_core_count" = "2",
        "customer_compute_version" = "ubuntu2004"
      },

      {
        "clustername" = "cluster02",
        "customer_prefix" = "10.120.71.0/24",
        "compute_node_size" = "medium",

        "compute_nodes" = [
          {
            "compute_node_name" : "esxi-web-01"
          },
          {
            "compute_node_name" : "esxi-web-02"
          },
          {
            "compute_node_name" : "esxi-web-03"
          }
        ],
        "datastores" = [
          {
            "datastore_name" : "datastore-web-01",
            "shared_drive_size" = "1TB"
          },
          {
            "datastore_name" : "datastore-web-02",
            "shared_drive_size" = "3TB"
          }
        ],
        "instance_array_ram_gbytes" = "4",
        "instance_array_processor_count" = "2",
        "instance_array_processor_core_count" = "4",
        "customer_compute_version" = "esxi7"
      },

      {
        "clustername" = "cluster03",
        "customer_prefix" = "10.120.72.0/24",
        "compute_node_size" = "large",

        "compute_nodes" = [
          {
            "compute_node_name" : "win-app-01"
          },
          {
            "compute_node_name" : "win-app-02"
          }
        ],
        "datastores" = [
          {
            "datastore_name" : "datastore-app-01",
            "shared_drive_size" = "500gb"
          },
          {
            "datastore_name" : "datastore-app-02",
            "shared_drive_size" = "750g"
          },
          {
            "datastore_name" : "datastore-app-03",
            "shared_drive_size" = "450Gb"
          }
        ],
        "instance_array_ram_gbytes" = "2",
        "instance_array_processor_count" = "1",
        "instance_array_processor_core_count" = "2",
        "customer_compute_version" = "win2012r2"
       }
    ]
}