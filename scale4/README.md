# Scale4

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


Now the app not just use the MySQL database but also has a cache layer using Redis.
We have a 10 second TTL on the cached value by default and has a cache invalidation call in case of a Save. If it's not a problem to have inconsistent data on the client side (which probably okay if we have a short cache timeout and a single page app/client side model which could "hide" it) then we can skip the invalidation. This will lead to a bit inconsistent data (we don't return the actual values), but gives us more performance on the read. (We'll return the data from the cache.)


#### Caching with lists
In this app we could store each message as an individual item in a Redis List instead of a single key-value par too and this way at save we could just add an item/update the cache, which could also benefit us, but in this case we should be extra careful about the data consistency.

The cache data could be lost any time and we should threat them like that.
Let's imagine we switched to a different cache server, but we have 20 messages in the db. The user calls the Save which will append 1 value to the DB (21 items) and 1 value to the Redis list (which is empty so 1 item).
Then in the GetAll we just check if the cache key has any values (which has, that 1) and return with that.
Now the user will see 1 item until the cache hits it's TTL.

So we should mitigate this issue somehow.
For example we can introduce a new key-value pair which can show us if the cache is ready. For example `cacheKey:something:ready` and it could have a true/false or any value.
Now each time we call the `GetAll` in the app, we check this first and if it's `true` then we can check and return from the cache.
When the TTL is timed out or the cache ready value is `false` then we can go to the DB, fetch the data, add them to the list one by one (be carefurl with the ordering!) and then set the `cacheKey:something:ready`.

The `Save` in this case could always just push to the list it does not need to check the ready key at all, because the `GetAll` will set it anyway.

Btw don't forget to re-set the cache TTL at the `Save`.
