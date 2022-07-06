resource "kubernetes_config_map" "mqtt-config" {
  metadata {
    name      = "mqtt-config"
    namespace = local.namespace
  }

  data = {
    "iot.conf" = file("${path.module}/iot.conf")
  }
}

resource "kubernetes_deployment" "mqtt" {
  metadata {
    name      = "mqtt"
    namespace = local.namespace
  }
  spec {
    selector {
      match_labels = {
        app   = "mqtt"
        tier  = "backend"
        track = "stable"
      }
    }
    template {
      metadata {
        labels = {
          app   = "mqtt"
          tier  = "backend"
          track = "stable"
        }
      }
      spec {
        container {
          image             = "eclipse-mosquitto"
          name              = "mqtt"
          image_pull_policy = "Always"
          volume_mount {
            mount_path = "/mosquitto/config"
            name       = "config-volume"
          }
          volume_mount {
            mount_path = "/mosquitto/certs"
            name       = "certs"
          }
          port {
            container_port = 8883
          }
          port {
            container_port = 31883
          }
          liveness_probe {
            failure_threshold     = 3
            initial_delay_seconds = 10
            period_seconds        = 10
            success_threshold     = 1
            tcp_socket {
              port = "31883"
            }
            timeout_seconds = 2
          }
          readiness_probe {
            failure_threshold     = 3
            initial_delay_seconds = 10
            period_seconds        = 10
            success_threshold     = 1
            tcp_socket {
              port = "31883"
            }
            timeout_seconds = 2
          }
        }
        volume {
          name = "config-volume"
          config_map {
            name = "mqtt-config"
            items {
              key  = "iot.conf"
              path = "mosquitto.conf"
            }
          }
        }
        volume {
          name = "certs"
          secret {
            secret_name = "mqtt-certs"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "mqtt" {
  metadata {
    name      = "mqtt"
    namespace = local.namespace
  }
  spec {
    selector = {
      app  = "mqtt"
      tier = "backend"
    }
    port {
      name     = "mqtts"
      port     = "8883"
      protocol = "TCP"
    }
  }
}