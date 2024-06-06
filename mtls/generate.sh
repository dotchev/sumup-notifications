#!/bin/bash

SCRIPT_DIR=$(realpath $(dirname "$0"))
docker run -it --rm -v $SCRIPT_DIR:/app -w /app --entrypoint /bin/sh hashicorp/terraform:1.8 -c 'terraform init && terraform apply -auto-approve'