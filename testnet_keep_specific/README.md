## KEEP network specific deployment steps

#### Install openVPN

```
helm install stable/openvpn --name keep-vpn --namespace=internal
```

#### Install Nginx Ingress controller

```
helm install stable/nginx-ingress --name keep-ingress-controller --namespace=internal --set controller.service.loadBalancerSourceRanges={10.0.0.0/8} --set rbac.create=true

```

#### Create TLS certificate

Use wildcard *.keep.network when asked for *common name*

```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls-key.key -out keep-tls-testnet.crt

kubectl create secret tls keep-tls-testnet --key tls-key.key --cert tls-testnet.crt
```


#### Deploy KEEP multihost ingress resource

```
kubectl apply -f keep/keep-testnet-ingress.yaml

```

#### Domain names

Get ip address (EXTERNAL-IP) of Nginx ingress controller

```
kubectl get svc keep-ingress-controller-nginx-ingress-controller
```

Add DNS records for the domains below or modify your local `etc/hosts` file:

```
<ip_address> testnet.keep.network
<ip_address> testnet-monitor.keep.network
<ip_address> testnet-token-dashboard.keep.network
<ip_address> testnet-multisig-wallet.keep.network
```

You should be able to access the services while connected to the VPN:

* [https://testnet.keep.network](https://testnet.keep.network)
* [https://testnet-monitor.keep.network](https://testnet.keep.network)
* [https://testnet-token-dashboard.keep.network](https://testnet-token-dashboard.keep.network)
* [https://testnet-multisig-wallet-ui.keep.network](https://testnet-multisig-wallet.keep.network)
