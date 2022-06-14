variable "storeids" {
  description = "JSON string for store ids"
  type        = list(any)
}

variable "static_client_token" {
  description = "Static client secret to use management REST API"
  type        = string
}

variable "volumeMounts" {
  description = "An array of volume mounts for the service, used for component versions"
  type        = map(string)
}

variable "volumes" {
  description = "An array of volumes for the volume mounts"
  type        = map(string)
}

variable "image" {
  description = "The image to be used for the service when deployed. (ex. localhost:32000/credential-agent, etc.)"
  type        = string
}

variable "host" {
  description = "The frontend host IP or FQDN"
  type        = string
}

variable "scheme" {
  description = "The protocol scheme to use for the frontend (http or https)"
  type        = string
}

variable "postgres_admin_configmap" {
  description = "The configmap name with the credentials of the admin of the Postgres instance being used so a database can be created for the service if needed"
  type        = string
}

variable "cluster_namespace" {
  description = "The namespace in the cluster where the service and associated resources should be deployed"
  type        = string
}

variable "component_postgres_configmap_name" {
  description = "The configmap name with the credentials to use for the service's database"
  type        = string
}

variable "log_level" {
  description = "The log level for the component"
  type        = string
}