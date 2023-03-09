# This is the Stackr base module for Azure
# It will build a resource group in a given subscription
# It will create a VNet if required and n subnets

resource "azurerm_resource_group" "stackr_system_rg" {
  name = var.stackrName
  location = var.stackrLocation

  tags = {
    "lzVersion" = "3.2",
    "billingOwningApp" = "SomeApp"
  }
}

