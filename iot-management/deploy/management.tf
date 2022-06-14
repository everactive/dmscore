data "kubectl_file_documents" "management_config_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-management-config.yaml",
    {
      HOST                = local.host
      STATIC_CLIENT_TOKEN = local.static_client_token
      SCHEME              = local.scheme
    }
  )
}

resource "kubectl_manifest" "iot_management_config" {
  override_namespace = local.namespace
  count              = 2
  yaml_body          = element(data.kubectl_file_documents.management_config_manifests.documents, count.index)
  wait               = true
}

data "kubectl_file_documents" "management_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-management.yaml",
    {
      IMAGE                             = local.image
      STOREIDS                          = local.storeids
      VOLUME_MOUNTS                     = local.volumeMounts
      VOLUMES                           = local.volumes
      POSTGRES_ADMIN_CONFIGMAP          = local.postgres_admin_configmap
      COMPONENT_POSTGRES_CONFIGMAP_NAME = var.component_postgres_configmap_name
      LOG_LEVEL                         = local.log_level
    }
  )
}

resource "kubectl_manifest" "iot_management" {
  depends_on = [
    kubectl_manifest.iot_management_config
  ]
  override_namespace = local.namespace
  count              = 2
  yaml_body          = element(data.kubectl_file_documents.management_manifests.documents, count.index)
  wait               = true
}