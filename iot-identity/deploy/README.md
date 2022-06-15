# Deployment

Deploys the IoT Identity Service to a Kubernetes cluster. This deployment uses Terraform to target a Kubernetes cluster. It is intended to be used as a module in
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
| [kubectl_manifest.iot_identity](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.iot_identity_manifests_version](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_file_documents.identity_manifests](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_namespace"></a> [cluster\_namespace](#input\_cluster\_namespace) | The namespace in the cluster where the service and associated resources should be deployed | `string` | n/a | yes |
| <a name="input_component_version"></a> [component\_version](#input\_component\_version) | The git summary version of the component | `string` | n/a | yes |
| <a name="input_create_version"></a> [create\_version](#input\_create\_version) | Determines whether the module should just create a version configmap or the service | `bool` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | The image to be used for the service when deployed. (ex. localhost:32000/credential-agent, etc.) | `string` | n/a | yes |
| <a name="input_mqtt_host_address"></a> [mqtt\_host\_address](#input\_mqtt\_host\_address) | The MQTT service host address FQDN without port to supply to devices when they enroll (ex. mqtt.example.com) | `string` | `"localhost"` | no |
| <a name="input_mqtt_host_port"></a> [mqtt\_host\_port](#input\_mqtt\_host\_port) | The port of the MQTT server to supply to devices when they enroll | `number` | `8883` | no |
| <a name="input_postgres_admin_configmap"></a> [postgres\_admin\_configmap](#input\_postgres\_admin\_configmap) | The configmap name with the credentials of the admin of the Postgres instance being used so a database can be created for the service if needed | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_identity_service"></a> [identity\_service](#output\_identity\_service) | n/a |
<!-- END_TF_DOCS -->