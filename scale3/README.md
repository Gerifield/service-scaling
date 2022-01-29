# Scale3

This is the third increment with message queue.


Star the services:

```
$ docker-compose -d up
```

Save some data:
```
$ curl -X POST 127.0.0.1:8080/save -d '{"content":"something"}'
```

Get back data:
```
$ curl -X GET 127.0.0.1:8080/
```

## Details

### Build

Everything is the same there's only an additional Redis instance.

### App Structure

The `main` for the project is under the `cmd/srv` folder. This way you can add multiple commands for the same project if it's needed.

The app itself is separated to different, independent layers (and packages):

- `app` - contains the business logic (which is quite empty now) it connects the database layer to the transport part and could have additional validations and check in it
- `db` - is the database layer which is implemented now only to use SQL, but it could have different implementations using different backends (file, Redis, NoSQL etc.). You just need to implement the same interface.
- `model` - act as the "common" data representation which moves between the layers. This way the communication between the packages could be well (and pre-) defined and will prevent the circular depency calls.
- `server` - currently this package implements the HTTP communication layer. Its responsibilty is to convert the input (json) data into the input model representation for a given logic call and then convert back the response model to the wire format.
- `cache` - contains a small Redis based cache layer for a given functionality
- `queue` - contains a small Redis based message queue implementation

Now the app separates the write and read sides, every write to the database goes through a message queue and a separate worker app does the actualy MySQL writing.

For this we've introduced a message format which we should put in the queue and should contain every information which is needed for the database action.

This way we can handle much higher traffic and spikes, because the queue put is much faster than a db write and we don't have to way for the answer.
As soon as the message is in the queue (and we assume the queue is durable and should not loose messages) we can return and answer to the user.
The increased usage/traffic spikes will just mean a bit more message waiting in the queue. This will increase the processing time, but won't lead to the potential failure of the database server.

On the other side we could have a fixed number of workers to do the queue reading and DB writing whihch leads to a much more predictable usage pattern and also possible to scale the `api` and the `worker` parts individually.

