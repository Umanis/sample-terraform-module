module "group" {
  source       = "../../"
  tags         = var.tags
  location     = var.location
  custom_name = var.custom_name
  description  = var.description
  caf_prefixes = var.caf_prefixes
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  features {}
}


variable "tags" {
  type = map(string)
}

variable "location" {
  type = string
}

variable "custom_name" {
  type = string
  default=""
}

variable "description" {
  type    = string
  default = ""
}

variable "caf_prefixes" {
  type = list(string)
  default = []
}