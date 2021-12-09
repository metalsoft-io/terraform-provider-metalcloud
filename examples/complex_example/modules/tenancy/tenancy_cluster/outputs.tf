output "credentials" {
   description = "the credentials of instances"
   sensitive = true
   value = metalcloud_instance_array.cluster[*].instances
   #value = {
   #  for k, ia in  metalcloud_instance_array.cluster[*].instances :  ia.instance_array_label => 
   #  {  for ilabel,i  in jsondecode("${ia.instances}"): ilabel => i.instance_credentials }
   #}
   
}
