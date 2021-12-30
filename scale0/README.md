# Scale0

This is the most basic app, without anything.


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
