#! /bin/bash

if ! command -v ytt &> /dev/null
then
    echo "ytt could not be found; for installation instruction visit https://carvel.dev/ytt/docs/latest/install"
    exit
fi


ytt \
    -f gen/template.yaml \
    -f gen/data.yaml \
    -f gen/schema.yaml \
    --file-mark 'template.yaml:path=keep-clients.yaml' \
    --output-files .

echo '# File generated with gen.sh - DO NOT EDIT' | cat - keep-clients.yaml > keep-clients.yaml.tmp && mv keep-clients.yaml.tmp keep-clients.yaml
