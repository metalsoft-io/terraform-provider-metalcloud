provider "metalcloud" {
   user = var.user
   api_key = var.api_key 
   endpoint = var.endpoint
}


resource "metalcloud_infrastructure" "my-infra35" {
  
  infrastructure_label = "my-terraform-infra57"
  datacenter_name = "us-santaclara"

  instance_arrays = [
    "test1" : {
      instance_array_label = "test1"
      instance_array_instance_count = 2
      drive_array {
        drive_array_storage_type = "iscsi_hdd"
        volume_template_label = "centos7-6"
        drive_size_mbytes_default = 40960
      }
  
  ]

  instance_array {
      instance_array_label = "test2"
      instance_array_instance_count = 2
      drive_array {
        drive_array_storage_type = "iscsi_hdd"
        volume_template_label = "centos7-6"
        drive_size_mbytes_default = 40961
      }
  }

 
}
