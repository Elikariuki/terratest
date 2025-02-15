# Define output values for later reference
output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}
 
output "vm_name" {
  value = azurerm_linux_virtual_machine.webserver.name
}
 
output "nic_name" {
  value = azurerm_network_interface.webserver.name
}
 
output "public_ip" {
  value = azurerm_linux_virtual_machine.webserver.public_ip_address
}
 
output "vm_nic" {
  value = azurerm_linux_virtual_machine.webserver.network_interface_ids
}
 
output "nic_id" {
  value = azurerm_network_interface.webserver.id
}
 
output "vm_version" {
  value = azurerm_linux_virtual_machine.webserver.source_image_reference[0].version
}
 
output "vm_sku" {
  value = azurerm_linux_virtual_machine.webserver.source_image_reference[0].sku
}