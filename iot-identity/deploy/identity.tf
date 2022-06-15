data "kubectl_file_documents" "identity_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-identity.yaml",
    {
      IMAGE        = local.image
      MQTT_HOST_ADDRESS      = local.mqtt_host_address
      MQTT_HOST_PORT = local.mqtt_host_port
      IDENTITY_VERSION = local.component_version
      POSTGRES_ADMIN_CONFIGMAP  = local.postgres_admin_configmap
      LOG_LEVEL                = local.log_level
    }
  )
}

resource "kubectl_manifest" "iot_identity_manifests_version" {
  override_namespace = local.namespace
  count              = 1
  yaml_body          = data.kubectl_file_documents.identity_manifests.documents[2]
  wait               = true
}

resource "kubectl_manifest" "iot_identity" {
  override_namespace = local.namespace
  count              = var.create_version ? 0 : 3
  yaml_body          = element(data.kubectl_file_documents.identity_manifests.documents, count.index)
  wait               = true
}
