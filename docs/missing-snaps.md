# Overview 

This feature is intended to provide a way to add snaps to a device that did not exist on the device when it was flashed.
If remodeling is not an option or does not fit the use case then DMS will add snaps to devices when they check 
in (heartbeat) if those snaps are not present.

# Design

3 additional tables will be added to the dmscore main database to have a table of models and ids and a table of snap
names and the model id they are associated with. The additional table is described in [#Heartbeat additions](#heartbeat-additions).

3 new REST API endpoints will be created to create and delete a snap requirement for a given model as well
as an endpoint to get all the required snaps for a given model.

* POST /v1/models/:model-name/required
* DELETE /v1/models/:model-name/required
* GET /v1/models/:model-name/required

A JSON payload similar to the following will be used for POST and DELETE:

```
{
  "name": "snap-name",
  "model": "model-name"
}
```

For GET it will just be a model:

```
{
  "model": "name"
}
```

The response will be an array:

```
{
  "snaps": [ "snap1", "snap2" ] 
}
```


## Heartbeat additions

Two additional fields will be added to the heartbeat:

```
{
   ... existing fields ...
  "snapListHash": "<hash of the JSON document that results from marshaling the result of the function List on the SnapdClient after sorting them by name alphabetically>",
  "installedSnapsHash": "<hash of the JSON document that results from adding the names of the snaps to an array, sorted alphabetically>"
}
```

These two additional pieces of information should give DMS all the information it needs to determine if it needs to request a new list and depending on
whether it has the latest or not and determine if this device needs to be told to install some additional snaps.

A table will be added to the dmscore main database to track this information containing a deviceID, snapListHash and installedSnapsHash. 

## PublishSnaps additions

For DMS to know what snaps a device has and whether it has the most current information the device will include the hashes outlined above when it 
publishes its list of snaps. This makes it so that DMS never has to worry about calculating the values itself. If the hashes every don't match in the database
it just requests the lastest snap list from the device which includes updated hashes.

# iot-agent

A new action will need to be added to iot-agent to install required snaps. This is advantageous over using the existing install action because
it will allow iot-agent to prioritize and persist the need, if necessary, to make sure it is done.

