  #!/usr/bin/env bash
set -e
echo 'Deploying SPI OAuth2 config'

OAUTH_URL='oauth.'$(minikube ip)'.nip.io'
tmpfile=/tmp/config.yaml

spiConfig=$(cat <<EOF

sharedSecret: $(openssl rand -hex 20)
serviceProviders:
  - type: GitHub
    clientId: $SPI_GITHUB_CLIENT_ID
    clientSecret: $SPI_GITHUB_CLIENT_SECRET
baseUrl: https://$OAUTH_URL

EOF
)

echo "Please go to https://github.com/settings/developers."
echo "And register new Github OAuth application for callback https://"$OAUTH_URL"/github/callback"

CONFIG_SECRET=$(kubectl get secrets  -l app.kubernetes.io/part-of=service-provider-integration-operator  -n spi-system -o json | jq '.items[0].metadata.name' -r)
echo $CONFIG_SECRET
#kubectl delete secret/$CONFIG_SECRET -n spi-system
echo "$spiConfig" > "$tmpfile"
cat $tmpfile


kubectl create secret generic  $CONFIG_SECRET \
--save-config --dry-run=client \
--from-file="$tmpfile"  \
-o yaml |
kubectl apply -n spi-system  -f -



#oc create namespace spi-system --dry-run=client -o yaml | oc apply -f -
#oc create secret generic $CONFIG_SECRET \
#    --save-config --dry-run=client \
#    --from-file="$tmpfile" \
#    -n spi-system \
#    -o yaml |
#oc apply -f -

rm "$tmpfile"


kubectl rollout restart  deployment/spi-controller-manager  -n spi-system
kubectl rollout restart  deployment/spi-oauth-service  -n spi-system
