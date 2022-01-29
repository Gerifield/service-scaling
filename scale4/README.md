# Scale4

This is the third increment with environment variable read.

We've switched the `flag` package to the `github.com/namsral/flag` to be able to pass env variable to the services.

This step will help us to use the container in a Kubernetes environment, becuase mounting a `configmap` to be able to use it's content as environment variables is much easier than realying on flags.
