variable create_version {
  description = "Determines whether the module should just create a version configmap or the service"
  type = bool
}

variable image {
  description = "The image to be used for the service when deployed. (ex. localhost:32000/credential-agent, etc.)"
  type = string
}

variable postgres_admin_configmap {
  description = "The configmap name with the credentials of the admin of the Postgres instance being used so a database can be created for the service if needed"
  type = string
}

variable cluster_namespace {
  description = "The namespace in the cluster where the service and associated resources should be deployed"
  type = string
}

variable "component_version" {
  description = "The git summary version of the component"
  type        = string
}

variable mqtt_host_port {
  description = "The port of the MQTT server to supply to devices when they enroll"
  type = number
  default = 8883
}

variable mqtt_host_address {
  description = "The MQTT service host address FQDN without port to supply to devices when they enroll (ex. mqtt.example.com)"
  type = string
  default = "localhost"
}

variable "log_level" {
  type = string
}
