### Note
All YAML files for deployments are in the k8s folder.

### Docker Image
The docker image repository is located at:
* https://hub.docker.com/repository/docker/tonylixu/go_http_server

### To Create Deployment
* Go into the k8s folder and run
```bash
$ kubectl create -f http_server_deployment.yaml
```

### To Create Service
* Go into the k8s folder and run
```bash
$ kubectl create -f http_server_service.yaml
```

### To Create Nginx ingress
```bash
$ kubectl create -f nginx-ingress-deployment.yaml
$ kubectl create -f ingress.yaml
```

### To Create secret
```bash
$ kubectl create -f secret.yaml
```