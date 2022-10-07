## Monitoring

This documentation was made upon this [tutorial](https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/)

The Kubernetes Prometheus monitoring stack has the following components:

1. Prometheus Server
2. Alert Manager
3. Grafana

![Schema](https://devopscube.com/wp-content/uploads/2022/01/kubernetes-1024x498.png)

### Prometheus

Version of Prometheus used in this deployment is [v2.37.0](https://hub.docker.com/layers/prometheus/prom/prometheus/v2.37.0/images/sha256-8ab20bc5a8bee3b8107bb2f533deea35da5641a608f9b0c16e683d6c60d3ee84?context=explore)

#### I. Create Namespace & ClusterRole

1. Execute the following command to create a new namespace named *monitoring*.
  
  ```bash
  kubectl create namespace monitoring
  ```

2. Additional Cluster Role Binding is required for prometheus  to view k8s metrics. In GCP a `clusterrolebinding` is required to be created.

```bash
ACCOUNT=$(gcloud info --format='value(config.account)')
kubectl create clusterrolebinding owner-cluster-admin-binding \
    --clusterrole cluster-admin \
    --user $ACCOUNT
```

*! This role can be created by Owner in GCP IAM*

3. Create `clusterRole`

  ```bash
  kubectl create -f clusterrole.yaml
  ```

#### II. Create `ConfigMap` to put configuration externally

1. Execute `create` command to create configmap
    
  ```bash
  kubectl create -f configmap.yaml
  ```

Some words to that `ConfigMap`:

All configurations are hold in `prometheus.yml`

Alert rules for Alertmanager are configured in `promethus.rules`

By externalizing Prometheus configs to a Kubernetes config map, there is no need to build Prometheus image whenever it needs additional configuration or some to be removed. Updating the config map and restarting the Prometheus pods to apply the new configuration is everything that needs to be done.

The config map with all the Prometheus scrape config and alerting rules gets mounted to the Prometheus container in `/etc/prometheus` location as prometheus.yaml and prometheus.rules files.

Beside that there is another ConfigMap located, isolatd from the one that belongs to Prometheus. It's created for Grafana.

All datasources are stored in `prometheus.yaml` file. It holds information about prometheus server location so that grafana would have that extension already configured

Dashboards config is stored in `dashboards.yaml`

For keep-client monitoring purpose `client-dashboards.json` holds template config for that service

One ConfigMap YAML file contains two section of configmaps configuration for two different services. 

#### III. Create Prometheus PersistentVolumeClaim

Create a PVC object

  ```bash
  kubectl create -f prometheus-pvc.yaml
  ```

#### IV. Create Prometheus Deployment

In this configuration, config map is mounted as a file inside /etc/prometheus as explained in the previous section.

  ```bash
  kubectl create  -f prometheus-deployment.yaml
  ```

#### IV. Create Prometheus Service

This object will use `NodePort` type for accessing the Prometheus service by using additional VPN service that is routed to kubernetes network. 

  ```bash
  kubectl create -f prometheus-service.yaml
  ```

Service should be available under

http://prometheus.monitoring.svc.cluster.local:9090

`prometheus` and `monitoring` in FQDN came from service object:

```
metadata:
  name: prometheus
  namespace: monitoring
```


### Grafana

#### I. Create Grafana ConfigMap

   ```bash
   kubectl create -f grafana-configmap.yaml
   ```

Notice the **url** key in the ConfigMap file that refers to prometheus FQDN that comes from combination of `metadata` in the prometheus-service YAML file.

#### II. Create Grafana PersistentVolumeClaim

  ```bash
  kubectl create -f grafana-pvc.yaml
  ```
#### IV. Create Grafana Deployment

  ```bash
  kubectl create -f grafana-deployment.yaml
  ```

#### V. Create Grafana Service

  ```bash
  kubectl create -f grafana-service.yaml
  ```

Service should be available under

http://grafana.monitoring.svc.cluster.local:3000/
