terraform {
  required_providers {
    kubectl = {
      source = "gavinbunney/kubectl"
    }
  }
}

locals {
  storage_class            = var.storage_class
  component_version        = var.component_version
  docker_image             = var.docker_image
  postgres_admin_configmap = var.postgres_admin_configmap
  namespace                = var.cluster_namespace
  configmap_name           = var.component_postgres_configmap_name
  log_level                = var.log_level
}