data "kubectl_file_documents" "mosquitto_manifests" {
  content = templatefile(
    "${path.module}/pods/k8s-mosquitto.yaml",
    {
    }
  )
}

resource "kubectl_manifest" "iot_mosquitto" {
  override_namespace = local.namespace
  count              = 3
  yaml_body          = element(data.kubectl_file_documents.mosquitto_manifests.documents, count.index)
  wait               = true
}