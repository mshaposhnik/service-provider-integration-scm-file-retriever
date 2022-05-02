#!/usr/bin/env bash
set -e
echo 'Deploying SPI OAuth2 config'

OAUTH_URL='spi-oauth-route-spi-system.'$( oc get ingresses.config/cluster -o jsonpath={.spec.domain})
tmpfile=/tmp/config.yaml

spiConfig=$(cat <<EOF

sharedSecret: $(openssl rand -hex 20)
serviceProviders:
  - type: GitHub
    clientId: $SPI_GITHUB_CLIENT_ID
    clientSecret: $SPI_GITHUB_CLIENT_SECRET
baseUrl: https://spi-oauth-route-spi-system.$( oc get ingresses.config/cluster -o jsonpath={.spec.domain})

EOF
)

echo "Please go to https://github.com/settings/developers."
echo "And register new Github OAuth application for callback https://"$OAUTH_URL"/github/callback"

echo "$spiConfig" > "$tmpfile"
oc create namespace spi-system --dry-run=client -o yaml | oc apply -f -
oc create secret generic oauth-config \
    --save-config --dry-run=client \
    --from-file="$tmpfile" \
    -n spi-system \
    -o yaml |
oc apply -f -

rm "$tmpfile"


oc rollout restart  deployment/spi-controller-manager  -n spi-system
oc rollout restart  deployment/spi-oauth-service  -n spi-system
