# Exmaple k8s files

I used docker-desktop for the example, here to have everyting working we'll need to install an Ingress controller first (after you enabled the kubernetes in the settings, this is an nginx one):

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/cloud/deploy.yaml
```

Then you can apply the configs:

```
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

If everything started and configured correctly you should be able to access the service through the Ingress -> Service -> Pod path on this url: [http://127.0.0.1/](http://127.0.0.1/)
