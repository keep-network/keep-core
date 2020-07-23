# Bare script to download each secret from a cluster.
# This assumes you're using the intended Kube context.
# This assumes you're on the correct VPN for that context.

# Downloaded secrets will have key values base64 encoded.  
# The last applied actuals should be in metadata.

# If you want decoded values in one swoop, third party
# tooling is required.  e.g. https://github.com/ashleyschuett/kubernetes-secret-decode

CURRENT_CONTEXT=$(kubectl config current-context)

printf "current kube context: [${CURRENT_CONTEXT}]\n\n"
printf "SECRETS TO BE DOWNLOADED:\n"

kubectl get secret --no-headers

kubectl get secret --no-headers | awk '{print $1}' | \
  xargs -I{} sh -c 'kubectl get secret -o yaml "$1" > "$1.yaml"' - {}