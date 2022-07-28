# dmscore
The core (backend) of the Device Management Service (DMS)


## Mocks

Unit testing mocks are being created with mockery and previous manually created
mocks are being phased out as time allows. The mocks are committed for convenience
to guarantee a match for each commit.

When creating mocks (or updating existing ones) output the mock to a `mocks` directory
that is peer to the interface file.

Ex.

```shell
mockery --dir iot-management/service/manage --name Manage  --output iot-management/service/manage/mocks
```

If you are creating or updating mocks for external packages output them to the top-level directory `mocks`
in the subdirectory `external`. Also output them to a directory matching their package name.

Ex. 

```shell
mockery --srcpkg="github.com/juju/usso/openid" --name=NonceStore --output mocks/external/openid
```