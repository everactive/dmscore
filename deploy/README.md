# Deployment

This deployment uses Terraform to target a Kubernetes cluster. It is intended to be used as a module in
an overall Device Management Service deployment with specifics to the targeted cluster (like EKS).


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
| [kubectl_manifest.iot_management](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.iot_management_config](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_file_documents.management_config_manifests](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |
| [kubectl_file_documents.management_manifests](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_namespace"></a> [cluster\_namespace](#input\_cluster\_namespace) | The namespace in the cluster where the service and associated resources should be deployed | `string` | n/a | yes |
| <a name="input_component_postgres_configmap_name"></a> [component\_postgres\_configmap\_name](#input\_component\_postgres\_configmap\_name) | The configmap name with the credentials to use for the service's database | `string` | n/a | yes |
| <a name="input_host"></a> [host](#input\_host) | The frontend host IP or FQDN | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | The image to be used for the service when deployed. (ex. localhost:32000/credential-agent, etc.) | `string` | n/a | yes |
| <a name="input_log_level"></a> [log\_level](#input\_log\_level) | The log level for the component | `string` | n/a | yes |
| <a name="input_postgres_admin_configmap"></a> [postgres\_admin\_configmap](#input\_postgres\_admin\_configmap) | The configmap name with the credentials of the admin of the Postgres instance being used so a database can be created for the service if needed | `string` | n/a | yes |
| <a name="input_scheme"></a> [scheme](#input\_scheme) | The protocol scheme to use for the frontend (http or https) | `string` | n/a | yes |
| <a name="input_static_client_token"></a> [static\_client\_token](#input\_static\_client\_token) | Static client secret to use management REST API | `string` | n/a | yes |
| <a name="input_storeids"></a> [storeids](#input\_storeids) | JSON string for store ids | `list(any)` | n/a | yes |
| <a name="input_volumeMounts"></a> [volumeMounts](#input\_volumeMounts) | An array of volume mounts for the service, used for component versions | `map(string)` | n/a | yes |
| <a name="input_volumes"></a> [volumes](#input\_volumes) | An array of volumes for the volume mounts | `map(string)` | n/a | yes |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
