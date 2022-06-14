terraform {
  required_providers {
    kubectl = {
      source = "gavinbunney/kubectl"
    }
  }
}

locals {
  storeids                 = jsonencode(var.storeids)
  static_client_token      = var.static_client_token
  volumeMounts             = var.volumeMounts
  volumes                  = var.volumes
  host                     = var.host
  scheme                   = var.scheme
  image                    = var.image
  postgres_admin_configmap = var.postgres_admin_configmap
  namespace                = var.cluster_namespace
  log_level                = var.log_level
  component_postgres_configmap_name = var.core_component_postgres_configmap_name
  identity_component_postgres_configmap_name = var.identity_component_postgres_configmap_name
  devicetwin_component_postgres_configmap_name = var.devicetwin_component_postgres_configmap_name
}