module "umanis_tagging" {
  source = <<path_to_module>>

  location          = "France Central"
  client            = "XY2"
  project           = "1234"
  budget            = "FE4567"
  rgpd_personal     = true
  rgpd_confidential = false
}

module "umanis_resource_group" {
  source = <<path_to_module>>

  tags         = module.umanis_tagging.tags
  location     = var.location
  description  = "Test resource group"
  custom_name  = "Custom name"
  caf_prefixes = module.umanis_naming.resource_group_prefixes
}
