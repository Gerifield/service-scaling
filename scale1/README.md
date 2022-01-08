# Scale1

This is the first increment which introduces some caching.


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

The project contains a Dockerfile and a docker-compose file. The build is a standard 2 stage process where the first container is responsible to create the binary and the second will be the output image. (It copies the binary from the build conainer.)

The final image is in a `gcr.io/distroless/static` container which contains the certs, timezone data and some other useful stuff and the binary will run with nonroot user too.

The docker-compose environment contains also a MySQL image which will act as a database.

### App Structure

The `main` for the project is under the `cmd/srv` folder. This way you can add multiple commands for the same project if it's needed.

The app itself is separated to different, independent layers (and packages):

- `app` - contains the business logic (which is quite empty now) it connects the database layer to the transport part and could have additional validations and check in it
- `db` - is the database layer which is implemented now only to use SQL, but it could have different implementations using different backends (file, Redis, NoSQL etc.). You just need to implement the same interface.
- `model` - act as the "common" data representation which moves between the layers. This way the communication between the packages could be well (and pre-) defined and will prevent the circular depency calls.
- `server` - currently this package implements the HTTP communication layer. Its responsibilty is to convert the input (json) data into the input model representation for a given logic call and then convert back the response model to the wire format.


This way you can very simply change the communication layer (add for example gRPC or websocket protocol support) or chang the storage layer without touching the application logic at all.

The interfaces are defined at the location where they are used, this way we can easily add tests for each layer because all the 3rd party package dependecies are handled through an interface.
