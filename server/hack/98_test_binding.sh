#!/usr/bin/env bash
set -e
echo 'Testing SPIAccessTokenBinding'


cat <<EOF | kubectl apply -n spi-system -f -
apiVersion: appstudio.redhat.com/v1beta1
kind: SPIAccessTokenBinding
metadata:
  name: read-private-repo-read
spec:
  permissions:
    required:
      - type: r
        area: admin:repo_hook
  repoUrl: https://github.com/skabashnyuk/some-private-repo
  secret:
    type: kubernetes.io/basic-auth
EOF



cat <<EOF | kubectl apply -n spi-system -f -
apiVersion: appstudio.redhat.com/v1beta1
kind: SPIAccessTokenBinding
metadata:
  name: read-private-repo-read-write
spec:
  permissions:
    required:
      - type: rw
        area: admin:repo_hook
  repoUrl: https://github.com/skabashnyuk/some-private-repo
  secret:
    type: kubernetes.io/basic-auth
EOF
