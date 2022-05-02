#!/usr/bin/env bash
set -e
echo 'Updating host variables'


SCM_HOST_VALUE='file-retriever-server-service-spi-system.'$(oc get ingresses.config/cluster -o jsonpath={.spec.domain})
OAUTH_URL='spi-oauth-route-spi-system.'$( oc get ingresses.config/cluster -o jsonpath={.spec.domain})
echo "scm="$SCM_HOST_VALUE
echo "oauth="$OAUTH_URL

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
#
#yq -i e '.spec.rules[0].host = "'$SCM_HOST_VALUE'"' $SCRIPT_DIR'/../config/k8s/ingress.yaml'
#jq 'map(select(.op == "replace").value |= "'$OAUTH_URL'")' $SCRIPT_DIR'/../config/k8s/ingress-patch.json' > tmp.$$.json && mv tmp.$$.json $SCRIPT_DIR'/../config/k8s/ingress-patch.json'
#
#
yq -i e '.sharedSecret = "'$(openssl rand -hex 20)'"' $SCRIPT_DIR'/../config/os/config.yaml'
yq -i e '.baseUrl = "https://'$OAUTH_URL'"' $SCRIPT_DIR'/../config/os/config.yaml'
#
#
echo "Please go to https://github.com/settings/developers."
echo "And register new Github OAuth application for callback https://"$OAUTH_URL"/github/callback"
echo "After that update Github's clientId and clientSecret in "$SCRIPT_DIR'/../config/os/config.yaml'
