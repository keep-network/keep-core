# Keep Client

## Configuration

### Generation

Keep Client Nodes manifests are generated with [`ytt`](https://carvel.dev/ytt/).

To generate the YAML configuration for nodes run `./gen.sh`.

ytt configuration consists of 3 files:

- [`template.yaml`](.gen/template.yaml) - template for kubernetes manifest
- [`schema.yaml`](./gen/schema.yaml) - properties with default values
- [`data.yaml`](./gen/data.yaml) - values for generation

### Resources

Manifests for `StatefulSet` and `Service` for all the nodes are generated into the [`keep-clients.yaml`](./keep-clients.yaml) file.

A node manifest reads values from following resources:

Config Maps:

- [`keep-client-config`](./keep-client-config.yaml)
- [`eth-account-info`](../eth-account-info-configmap.yaml)

Secrets:

- `eth-network-goerli`
- `eth-account-passphrases`
- `eth-account-privatekeys`
