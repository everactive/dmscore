data "kubectl_file_documents" "devicetwin_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-devicetwin.yaml",
    {
      IMAGE                             = local.docker_image
      DEVICETWIN_VERSION                = local.component_version
      STORAGE_CLASS                     = local.storage_class
      POSTGRES_ADMIN_CONFIGMAP          = local.postgres_admin_configmap
      COMPONENT_POSTGRES_CONFIGMAP_NAME = local.configmap_name
      LOG_LEVEL                         = local.log_level
    }
  )
}

resource "kubectl_manifest" "devicetwin" {
  depends_on         = [kubectl_manifest.iot_mosquitto]
  override_namespace = local.namespace
  count              = 4
  yaml_body          = element(data.kubectl_file_documents.devicetwin_manifests.documents, count.index)
  wait               = true
}
