## Monitoring

This documentation was made upon this [tutorial](https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/)

The Kubernetes Prometheus monitoring stack has the following components.

1. Prometheus Server
2. Alert Manager
3. Grafana

![Schema](https://devopscube.com/wp-content/uploads/2022/01/kubernetes-1024x498.png)


Version of Prometheus used in this deployment is [v2.37.0](https://hub.docker.com/layers/prometheus/prom/prometheus/v2.37.0/images/sha256-8ab20bc5a8bee3b8107bb2f533deea35da5641a608f9b0c16e683d6c60d3ee84?context=explore)

#### I. Create Namespace & ClusterRole

1. Execute the following command to create a new namespace named *monitoring*.
  
  ```bash
  kubectl create namespace monitoring
  ```

2. Additionaly Cluster Role Binding is required for prometheus  to view k8s metrics. In GCP a `clusterrolebinding` is required to be created.

```bash
ACCOUNT=$(gcloud info --format='value(config.account)')
kubectl create clusterrolebinding owner-cluster-admin-binding \
    --clusterrole cluster-admin \
    --user $ACCOUNT
```

*! This role can be created by Owner in GCP IAM*

3. Create `clusterRole`

  ```bash
  kubectl create -f clusterRole.yaml
  ```

#### II. Create `ConfigMap` to put configuration externally

1. Execute `create` command to create configmap
    
  ```bash
  kubectl create -f config-map.yaml
  ```

>Some words to that `ConfigMap`
All configurations are hold in `prometheus.yml`
Alert rules for Alertmanager are configured in `promethus.rules`

>By externalizing Prometheus configs to a Kubernetes config map, there is no need to build Prometheus image whenever it needs additional configuration or some to be removed. Updating the config map and restarting the Prometheus pods to apply the new configuration is everything that needs to be done.

>The config map with all the Prometheus scrape config and alerting rules gets mounted to the Prometheus container in /etc/prometheus location as prometheus.yaml and prometheus.rules files.


#### III. Create Prometheus Deployment

1.  In this configuration, config map is mounted as a file inside /etc/prometheus as explained in the previous section.

  ```bash
  kubectl create  -f prometheus-deployment.yaml
  ```

## TO BE CONTINUED :)
