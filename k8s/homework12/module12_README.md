### Module 12 homework
Requirements:
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：
* 如何实现安全保证；
* 七层路由规则；
* 考虑 open tracing 的接入。

### Docker Image
The docker image repository is located at:
* https://hub.docker.com/repository/docker/tonylixu/go_http_server

### Install istio
```bash
$ curl -L https://istio.io/downloadIstio | sh -
$ cd istio-1.12.1/
$ sudo cp bin/istioctl /usr/local/bin
$ istioctl install --set profile=demo -y
✔ Istio core installed
✔ Istiod installed
✔ Ingress gateways installed
✔ Egress gateways installed
✔ Installation complete                                                                                                                                   Making this installation the default for injection and validation.

Thank you for installing Istio 1.12.  Please take a few minutes to tell us about your install/upgrade experience!  https://forms.gle/FegQbc9UvePd4Z9z7
```

### Deploy httpserver
* Create namespace `httpserver`
```bash
$ kubectl create ns httpserver
```
* Enable istio for new namespace:
```bash
$ kubectl label ns httpserver istio-injection=enabled
```
* Create httpserver service
```
$ kubectl create -f httpserver.yaml
```

### Deploy new code to k8s
```bash
$ k create -f k8s/module10.yaml

$ k get po -n http-server
NAME                                      READY   STATUS    RESTARTS   AGE
http-server-deployment-66d4fbc946-2r2b5   1/1     Running   0          7m32s
http-server-deployment-66d4fbc946-fpp2c   1/1     Running   0          7m32s
http-server-deployment-66d4fbc946-hqznz   1/1     Running   0          7m32s

$ k get svc -n http-server
NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
http-server-service   ClusterIP   10.111.177.185   <none>        80/TCP    7m56s

$ curl http://10.111.177.185/metrics
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 4.0942e-05
go_gc_duration_seconds{quantile="0.25"} 4.0942e-05
go_gc_duration_seconds{quantile="0.5"} 8.7557e-05
go_gc_duration_seconds{quantile="0.75"} 8.7557e-05
go_gc_duration_seconds{quantile="1"} 8.7557e-05
go_gc_duration_seconds_sum 0.000128499
go_gc_duration_seconds_count 2
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 10
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.17.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
....
```

### Expose prometheus server
```bash
$ k edit svc loki-prometheus-server
service/loki-prometheus-server edited
ubuntu@ip-172-31-82-231:~/go_http_server/k8s$ k get svc
NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
kubernetes                      ClusterIP   10.96.0.1        <none>        443/TCP        50m
loki                            ClusterIP   10.107.65.27     <none>        3100/TCP       45m
loki-grafana                    NodePort    10.108.243.191   <none>        80:30339/TCP   45m
loki-headless                   ClusterIP   None             <none>        3100/TCP       45m
loki-kube-state-metrics         ClusterIP   10.98.58.192     <none>        8080/TCP       45m
loki-prometheus-alertmanager    ClusterIP   10.108.52.140    <none>        80/TCP         45m
loki-prometheus-node-exporter   ClusterIP   None             <none>        9100/TCP       45m
loki-prometheus-pushgateway     ClusterIP   10.104.17.17     <none>        9091/TCP       45m
loki-prometheus-server          NodePort    10.110.163.119   <none>        80:32044/TCP   45m
```
You can also check the prometheus web UI now.
* loki-prometheus screenshot
![Prometheus Web UI](./images/loki-prometheus.png)

### Check httpserver delayed metrics
* Delayed metrics
![Service httpserver delay](./images/loki-prometheus-httpserver.png)


### Grafana Dashboard
* Httpserver grafana dashboard screenshot
![Httpserver dashboard](./images/grafana-dashboard.png)
