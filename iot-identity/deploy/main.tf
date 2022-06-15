terraform {
  required_providers {
    kubectl = {
      source = "gavinbunney/kubectl"
    }
  }
}
locals {
  mqtt_host_address = var.mqtt_host_address
  mqtt_host_port = var.mqtt_host_port
  image = var.image
  component_version = var.component_version
  postgres_admin_configmap = var.postgres_admin_configmap
  namespace                = var.cluster_namespace
  log_level                = var.log_level
}
