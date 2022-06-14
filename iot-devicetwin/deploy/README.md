# Deployment

Deploys the IoT Device Twin Service to a Kubernetes cluster. This deployment uses Terraform to target a Kubernetes cluster. It is intended to be used as a module in
an overall Device Management Service deployment with specifics to the targeted cluster (like EKS).

## Terraform Docs


<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_kubectl"></a> [kubectl](#provider\_kubectl) | 1.10.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [kubectl_manifest.devicetwin](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.iot_mosquitto](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_file_documents.devicetwin_manifests](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |
| [kubectl_file_documents.mosquitto_manifests](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_namespace"></a> [cluster\_namespace](#input\_cluster\_namespace) | The namespace in the cluster where the service and associated resources should be deployed | `string` | n/a | yes |
| <a name="input_component_postgres_configmap_name"></a> [component\_postgres\_configmap\_name](#input\_component\_postgres\_configmap\_name) | The configmap name with the credentials to use for the service's database | `string` | n/a | yes |
| <a name="input_component_version"></a> [component\_version](#input\_component\_version) | The git summary version of the component to be used to create a configmap for consumption in the cluster | `string` | `"UNSET"` | no |
| <a name="input_docker_image"></a> [docker\_image](#input\_docker\_image) | The docker image, full path, to use for the device twin service | `string` | n/a | yes |
| <a name="input_log_level"></a> [log\_level](#input\_log\_level) | The log level to use for the service | `string` | n/a | yes |
| <a name="input_postgres_admin_configmap"></a> [postgres\_admin\_configmap](#input\_postgres\_admin\_configmap) | The configmap name with the credentials of the admin of the Postgres instance being used so a database can be created for the service if needed | `string` | n/a | yes |
| <a name="input_storage_class"></a> [storage\_class](#input\_storage\_class) | The storage class to use for the device twin's persistent volume | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
