variable "storage_class" {
  description = "The storage class to use for the device twin's persistent volume"
  type        = string
}

variable "docker_image" {
  description = "The docker image, full path, to use for the device twin service"
  type        = string
}

variable "component_version" {
  description = "The git summary version of the component to be used to create a configmap for consumption in the cluster"
  default     = "UNSET"
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
  description = "The log level to use for the service"
  type        = string
}