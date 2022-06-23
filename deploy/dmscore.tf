resource "kubernetes_secret" "dmscore-secrets" {
  metadata {
    name = "dmscore-secrets"
    namespace = local.namespace
  }
  data = {
    static-token=local.static_client_token
  }
}

resource "kubernetes_config_map" "dmscore-config" {
  metadata {
    name = "dmscore-config"
    namespace = local.namespace
  }

  data = {
    HOST=local.host
    SCHEME=local.scheme
  }
}

resource "kubernetes_deployment" "dmscore" {
  depends_on = [kubernetes_config_map.dmscore-config, kubernetes_secret.dmscore-secrets]
  metadata {
    name      = "dmscore"
    namespace = local.namespace
  }
  spec {
    selector {
      match_labels = {
        app   = "dmscore"
        tier = "frontend"
        track = "stable"
      }
    }
    template {
      metadata {
        labels = {
          app   = "dmscore"
          tier = "frontend"
          track = "stable"
        }
      }
      spec {
        volume {
          name = "sql"
          config_map {
            name = local.component_postgres_configmap_name
            items {
              key = "DB_CREATE_SQL"
              path = "create-db.sql"
            }
          }
        }
        volume {
          name = "sql-identity"
          config_map {
            name = local.identity_component_postgres_configmap_name
            items {
              key = "DB_CREATE_SQL"
              path = "create-db.sql"
            }
          }
        }
        volume {
          name = "sql-devicetwin"
          config_map {
            name = local.devicetwin_component_postgres_configmap_name
            items {
              key = "DB_CREATE_SQL"
              path = "create-db.sql"
            }
          }
        }
        volume {
          name = "identity-certs"
          secret {
            secret_name = "identity-certs"
          }
        }
        volume {
          name = "devicetwin-certs"
          secret {
            secret_name = "devicetwin-certs"
          }
        }
        init_container {
          name = "create-db"
          image = "jbergknoff/postgresql-client"
          command = ["psql", "-d", "$(DATASOURCE)", "-f", "/sql/create-db.sql"]
          env {
            name = "DATASOURCE"
            value_from {
              config_map_key_ref {
                name = local.postgres_admin_configmap
                key = "DATASOURCE"
              }
            }
          }
          volume_mount {
            mount_path = "/sql"
            name       = "sql"
          }
        }
        init_container {
          name = "create-db-identity"
          image = "jbergknoff/postgresql-client"
          command = ["psql", "-d", "$(DATASOURCE)", "-f", "/sql/create-db.sql"]
          env {
            name = "DATASOURCE"
            value_from {
              config_map_key_ref {
                name = local.postgres_admin_configmap
                key = "DATASOURCE"
              }
            }
          }
          volume_mount {
            mount_path = "/sql"
            name       = "sql-identity"
          }
        }
        init_container {
          name = "create-db-devicetwin"
          image = "jbergknoff/postgresql-client"
          command = ["psql", "-d", "$(DATASOURCE)", "-f", "/sql/create-db.sql"]
          env {
            name = "DATASOURCE"
            value_from {
              config_map_key_ref {
                name = local.postgres_admin_configmap
                key = "DATASOURCE"
              }
            }
          }
          volume_mount {
            mount_path = "/sql"
            name       = "sql-devicetwin"
          }
        }
        container {
          image             = local.image
          name              = "dmscore"
          image_pull_policy = "Always"
          volume_mount {
            mount_path = "/srv/identity-certs"
            name       = "identity-certs"
          }
          volume_mount {
            mount_path = "/srv/devicetwin-certs"
            name       = "devicetwin-certs"
          }
          env {
            name  = "LOG_LEVEL"
            value = local.log_level
          }
          env {
            name  = "DMS_DATABASE_DRIVER"
            value = "postgres"
          }
          env {
            name  = "DMS_DATABASE_CONNECTION_STRING"
            value_from {
              config_map_key_ref {
                name = local.component_postgres_configmap_name
                key = "DATASOURCE"
              }
            }
          }
          env {
            name  = "DMS_SERVICE_HOST"
            value_from {
              config_map_key_ref {
                name = "dmscore-config"
                key = "HOST"
              }
            }
          }
          env {
            name  = "DMS_SERVICE_SCHEME"
            value_from {
              config_map_key_ref {
                name = "dmscore-config"
                key = "SCHEME"
              }
            }
          }
          env {
            name = "DMS_STORE_URL"
            value = "https://api.snapcraft.io/api/v1/"
          }
          env {
            name = "LOG_FORMAT"
            value = "json"
          }
          env {
            name = "DMS_STORE_IDS"
            value = local.storeids
          }
          env {
            name = "DMS_MQTT_HOST_ADDRESS"
            value = local.mqtt_host_address
          }
          env {
            name = "DMS_MQTT_HOST_PORT"
            value = local.mqtt_host_port
          }
          env {
            name = "DMS_SERVICE_CLIENT_TOKEN_PROVIDER"
            value = local.client_token_provider
          }
          env {
            name = "DMS_SERVICE_AUTH_PROVIDER"
            value = local.auth_provider
          }
          env {
            name = "DMS_SERVICE_AUTH_DISABLED"
            value = local.auth_disabled
          }
          env {
            name  = "DMS_STATIC_CLIENT_TOKEN"
            value_from {
              secret_key_ref {
                name = "dmscore-secrets"
                key = "static-token"
              }
            }
          }
          env {
            name = "DMS_SERVICE_JWTSECRET"
            value = "this-is-a-really-long-secret-for-jwt-do-not-use-in-production"
          }
          env {
            name = "DMS_SERVICE_MIGRATIONS_SOURCE"
            value = "/migrations/dmscore"
          }
          env {
            name = "DMS_DEVICETWIN_DATABASE_CONNECTION_STRING"
            value_from {
              config_map_key_ref {
                name = local.devicetwin_component_postgres_configmap_name
                key = "DATASOURCE"
              }
            }
          }
          env {
            name = "DMS_IDENTITY_SERVICE_PORT_ENROLL"
            value = local.identity_service_enroll_port
          }
          env {
            name = "DMS_IDENTITY_DATABASE_CONNECTION_STRING"
            value_from {
              config_map_key_ref {
                name = local.identity_component_postgres_configmap_name
                key = "DATASOURCE"
              }
            }
          }
          port {
            name = "mgmt-svc"
            container_port = 8010
          }
          port {
            name = "id-svc"
            container_port = 8040
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "dmscore-internal" {
  metadata {
    name      = "dmscore-internal"
    namespace = local.namespace
  }
  spec {
    selector = {
      app   = "dmscore"
      tier = "frontend"
      track = "stable"
    }
    port {
      port        = "8010"
      protocol = "TCP"
    }
  }
}

resource "kubernetes_service" "dmscore-enroll" {
  metadata {
    name      = "dmscore-enroll"
    namespace = local.namespace
  }
  spec {
    selector = {
      app   = "dmscore"
      tier = "frontend"
      track = "stable"
    }
    port {
      port        = local.identity_service_enroll_port
      protocol = "TCP"
    }
  }
}