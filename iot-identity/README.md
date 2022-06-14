# IoT Identity Service

Managing the identity, ownership, credentials and authorization of an IoT device plays a crucial role in the security story. Those details need to be managed as the device goes through its lifecycle - from the manufacturer, distributor, system integrator, to end customer; from commissioning, repurposing to decommissioning the device.

The Identity service plays the role of managing these assets and enabling the connected systems to communicate with secure credentials.

The Identity Service is primarily in focus when the new device comes online. The device will be preconfigured to connect to the Identity Service, providing its Model and Serial assertions. The Identity Service registry will contain the primary ownership details for the device (customer name, store ID) and generates certificates and credentials for the device.

## Build

```shell
go build bin/identity/identity.go
```

## Run
```shell
go run bin/identity/identity.go
```

## Configuration

The service currently runs two different servers for endpoints; one for internal REST APIs used by the 
management service and one for external clients to enroll. By default those ports are 8031 and 8030 respectively.

Configuration is set using either a YAML file or environment variables. Environment variables are the
YAML keys in uppercase with periods replaced with underscores and prefaced with IOTIDENTITY_. Ex. `database.driver` 
becomes `IOTIDENTITY_DATABASE_DRIVER` as an environment variable.

### Keys

`database.driver` - The driver to use for the database, supported values are: `memory` and `postgres`.
See [factory.go](./service/factory/factory.go)

`database.connection.string` - A connection string formatted for the datasource to consume; passed during factory.
Currently just `postgres` utilizes this. See [config.go](./config/config.go) for an example.

`service.port.internal` - The port to use for internal REST API endpoints, consumed by the management service

`service.port.enroll` - The port to expose for clients to enroll with the Device Management Service

`mqtt.url` - The URL of the MQTT broker without the port, i.e. `localhost` or `mqtt.somedomain.com`

`mqtt.port` - The MQTT broker port used to communicate with the clients. 

`mqtt.certificate.path` - The path to the certificates to use to connect to the MQTT broker